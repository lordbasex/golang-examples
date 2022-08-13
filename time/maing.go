package main

import (
	"fmt"
	"time"
)

func inLocTimestamp(inLoc time.Time) int64 {
	_, offset := inLoc.Zone()
	ts := inLoc.Add(time.Duration(offset) * time.Second)
	return ts.Unix()
}

func inLocTimestampMilli(inLoc time.Time) int64 {
	_, offset := inLoc.Zone()
	ts := inLoc.Add(time.Duration(offset) * time.Second)
	return ts.UnixMilli()
}

func main() {
	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
	now := time.Now().Round(0)
	inLoc := now.In(loc)
	fmt.Printf("time       | now: %v | inLoc: %v\n", now, inLoc)
	fmt.Printf("unix       | now: %v | inLoc: %v\n", now.Unix(), inLoc.Unix())
	fmt.Printf("timestamp  | now: %v | inLoc: %v\n", now.Unix(), inLocTimestamp(inLoc))
	fmt.Printf("difference | inLoc: timestamp - unix = %v\n", inLocTimestamp(inLoc)-inLoc.Unix())
}
