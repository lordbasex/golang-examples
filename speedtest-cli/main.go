package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chelnak/ysmrr"
	"github.com/showwin/speedtest-go/speedtest"
	"github.com/showwin/speedtest-go/speedtest/transport"
)

var (
	commit = "dev"
	date   = "unknown"
)

func main() {
	AppInfo()

	// discard standard log.
	log.SetOutput(io.Discard)

	// 0. speed test setting
	var speedtestClient = speedtest.New()

	// 1. retrieving user information
	taskManager := InitTaskManager(false, false)
	taskManager.AsyncRun("Retrieving User Information", func(task *Task) {
		u, err := speedtestClient.FetchUserInfo()
		task.CheckError(err)
		task.Printf("ISP: %s", u.String())
		task.Complete()
	})

	// 2. retrieving servers
	var err error
	var servers speedtest.Servers
	var targets speedtest.Servers
	taskManager.Run("Retrieving Servers", func(task *Task) {
		servers, err = speedtestClient.FetchServers()
		task.CheckError(err)
		task.Printf("Found %d Public Servers", len(servers))
		targets, err = servers.FindServer([]int{})
		task.CheckError(err)
		task.Complete()
	})
	taskManager.Reset()

	// 3. test each selected server with ping, download and upload.
	for _, server := range targets {
		fmt.Println()
		taskManager.Println("Test Server: " + server.String())
		taskManager.Run("Latency: --", func(task *Task) {
			task.CheckError(server.PingTest(func(latency time.Duration) {
				task.Updatef("Latency: %v", latency)
			}))
			task.Printf("Latency: %v Jitter: %v Min: %v Max: %v", server.Latency, server.Jitter, server.MinLatency, server.MaxLatency)
			task.Complete()
		})

		// 3.0 create a packet loss analyzer, use default options
		analyzer := speedtest.NewPacketLossAnalyzer(nil)

		blocker := sync.WaitGroup{}
		packetLossAnalyzerCtx, packetLossAnalyzerCancel := context.WithTimeout(context.Background(), time.Second*40)
		taskManager.Run("Packet Loss Analyzer", func(task *Task) {
			blocker.Add(1)
			go func() {
				defer blocker.Done()
				err = analyzer.RunWithContext(packetLossAnalyzerCtx, server.Host, func(packetLoss *transport.PLoss) {
					server.PacketLoss = *packetLoss
				})
				if errors.Is(err, transport.ErrUnsupported) {
					packetLossAnalyzerCancel() // cancel early
				}
			}()
			task.Println("Packet Loss Analyzer: Running in background (<= 30 Secs)")
			task.Complete()
		})

		// 3.1 create accompany Echo
		accEcho := newAccompanyEcho(server, time.Millisecond*500)
		taskManager.Run("Download", func(task *Task) {
			accEcho.Run()
			speedtestClient.SetCallbackDownload(func(downRate speedtest.ByteRate) {
				lc := accEcho.CurrentLatency()
				if lc == 0 {
					task.Updatef("Download: %s (Latency: --)", downRate)
				} else {
					task.Updatef("Download: %s (Latency: %dms)", downRate, lc/1000000)
				}
			})
			task.CheckError(server.DownloadTest())
			accEcho.Stop()
			mean, _, std, minL, maxL := speedtest.StandardDeviation(accEcho.Latencies())
			task.Printf("Download: %s (Used: %.2fMB) (Latency: %dms Jitter: %dms Min: %dms Max: %dms)", server.DLSpeed, float64(server.Context.Manager.GetTotalDownload())/1000/1000, mean/1000000, std/1000000, minL/1000000, maxL/1000000)
			task.Complete()
		})

		taskManager.Run("Upload", func(task *Task) {
			accEcho.Run()
			speedtestClient.SetCallbackUpload(func(upRate speedtest.ByteRate) {
				lc := accEcho.CurrentLatency()
				if lc == 0 {
					task.Updatef("Upload: %s (Latency: --)", upRate)
				} else {
					task.Updatef("Upload: %s (Latency: %dms)", upRate, lc/1000000)
				}
			})
			task.CheckError(server.UploadTest())
			accEcho.Stop()
			mean, _, std, minL, maxL := speedtest.StandardDeviation(accEcho.Latencies())
			task.Printf("Upload: %s (Used: %.2fMB) (Latency: %dms Jitter: %dms Min: %dms Max: %dms)", server.ULSpeed, float64(server.Context.Manager.GetTotalUpload())/1000/1000, mean/1000000, std/1000000, minL/1000000, maxL/1000000)
			task.Complete()
		})

		packetLossAnalyzerCancel()
		blocker.Wait()
		taskManager.Println(server.PacketLoss.String())
		taskManager.Reset()
		speedtestClient.Manager.Reset()
	}
	taskManager.Stop()
}

type AccompanyEcho struct {
	stopEcho       chan bool
	server         *speedtest.Server
	currentLatency int64
	interval       time.Duration
	latencies      []int64
}

func newAccompanyEcho(server *speedtest.Server, interval time.Duration) *AccompanyEcho {
	return &AccompanyEcho{
		server:   server,
		interval: interval,
		stopEcho: make(chan bool),
	}
}

func (ae *AccompanyEcho) Run() {
	ae.latencies = make([]int64, 0)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ae.stopEcho:
				cancel()
				return
			default:
				latency, _ := ae.server.HTTPPing(ctx, 1, ae.interval, nil)
				if len(latency) > 0 {
					atomic.StoreInt64(&ae.currentLatency, latency[0])
					ae.latencies = append(ae.latencies, latency[0])
				}
			}
		}
	}()
}

func (ae *AccompanyEcho) Stop() {
	ae.stopEcho <- false
}

func (ae *AccompanyEcho) CurrentLatency() int64 {
	return atomic.LoadInt64(&ae.currentLatency)
}

func (ae *AccompanyEcho) Latencies() []int64 {
	return ae.latencies
}

func AppInfo() {
	fmt.Println()
	fmt.Printf("    speedtest-go v%s (git-%s, built %s) @showwin\n", speedtest.Version(), commit, date)
	fmt.Println()
}

// TaskManager handles the display of task progress with spinners
type TaskManager struct {
	sm         ysmrr.SpinnerManager
	isOut      bool
	noProgress bool
}

// Task represents a single task with a spinner
type Task struct {
	spinner *ysmrr.Spinner
	manager *TaskManager
	title   string
}

func InitTaskManager(jsonOutput, unixOutput bool) *TaskManager {
	isOut := !jsonOutput || unixOutput
	tm := &TaskManager{sm: ysmrr.NewSpinnerManager(), isOut: isOut, noProgress: unixOutput}
	if isOut && !unixOutput {
		tm.sm.Start()
	}
	return tm
}

func (tm *TaskManager) Reset() {
	if tm.isOut && !tm.noProgress {
		tm.sm.Stop()
		tm.sm = ysmrr.NewSpinnerManager()
		tm.sm.Start()
	}
}

func (tm *TaskManager) Stop() {
	if tm.isOut && !tm.noProgress {
		tm.sm.Stop()
	}
}

func (tm *TaskManager) Println(message string) {
	if tm.noProgress {
		fmt.Println(message)
		return
	}
	if tm.isOut {
		context := &Task{manager: tm}
		context.spinner = tm.sm.AddSpinner(message)
		context.Complete()
	}
}

func (tm *TaskManager) Run(title string, callback func(task *Task)) {
	context := &Task{manager: tm, title: title}
	if tm.isOut {
		if tm.noProgress {
			//fmt.Println(title)
		} else {
			context.spinner = tm.sm.AddSpinner(title)
		}
	}
	callback(context)
}

func (tm *TaskManager) AsyncRun(title string, callback func(task *Task)) {
	context := &Task{manager: tm, title: title}
	if tm.isOut {
		if tm.noProgress {
			//fmt.Println(title)
		} else {
			context.spinner = tm.sm.AddSpinner(title)
		}
	}
	go callback(context)
}

func (t *Task) Complete() {
	if t.manager.noProgress {
		return
	}
	if t.spinner == nil {
		return
	}
	t.spinner.Complete()
}

func (t *Task) Updatef(format string, a ...interface{}) {
	if t.spinner == nil || t.manager.noProgress {
		return
	}
	t.spinner.UpdateMessagef(format, a...)
}

func (t *Task) Println(message string) {
	if t.manager.noProgress {
		fmt.Println(message)
		return
	}
	if t.spinner == nil {
		return
	}
	t.spinner.UpdateMessage(message)
}

func (t *Task) Printf(format string, a ...interface{}) {
	if t.manager.noProgress {
		fmt.Printf(format+"\n", a...)
		return
	}
	if t.spinner == nil {
		return
	}
	t.spinner.UpdateMessagef(format, a...)
}

func (t *Task) CheckError(err error) {
	if err != nil {
		if t.spinner != nil {
			t.Printf("Fatal: %s, err: %v", t.title, err)
			t.spinner.Error()
			t.manager.Stop()
		} else {
			fmt.Printf("Fatal: %s, err: %v", t.title, err)
		}
		os.Exit(1)
	}
}
