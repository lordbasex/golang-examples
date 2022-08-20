package main

import (
        "fmt"
        "log"
        "math/rand"
        "os"
        "time"

        amqp "github.com/rabbitmq/amqp091-go"
)

const (
	      rbHost    = "172.16.20.10"
	      rbPort    = "8888"
	      rbUser    = "rabbitUser"
	      rbPass    = "AASwPslfkjJs2ijsnfiujhaADXKjbsadkjbdasdc222asd11A"
	      queueName = "AsteriskServerRemoteCLI"
)

func failOnError(err error, msg string) {
        if err != nil {
                log.Panicf("%s: %s", msg, err)
        }
}

func randomString(l int) string {
        bytes := make([]byte, l)
        for i := 0; i < l; i++ {
                bytes[i] = byte(randInt(65, 90))
        }
        return string(bytes)
}

func randInt(min int, max int) int {
        return min + rand.Intn(max-min)
}

func fibonacciRPC(n string) (res string, err error) {
        conn, err := amqp.Dial("amqp://" + rbUser + ":" + rbPass + "@" + rbHost + ":" + rbPort + "/")
        failOnError(err, "Failed to connect to RabbitMQ")
        defer conn.Close()

        ch, err := conn.Channel()
        failOnError(err, "Failed to open a channel")
        defer ch.Close()

        q, err := ch.QueueDeclare(
                "",    // name
                false, // durable
                false, // delete when unused
                true,  // exclusive
                false, // noWait
                nil,   // arguments
        )
        failOnError(err, "Failed to declare a queue")

        msgs, err := ch.Consume(
                q.Name, // queue
                "",     // consumer
                true,   // auto-ack
                false,  // exclusive
                false,  // no-local
                false,  // no-wait
                nil,    // args
        )
        failOnError(err, "Failed to register a consumer")

        corrId := randomString(32)

        err = ch.Publish(
                "",        // exchange
                queueName, // routing key
                false,     // mandatory
                false,     // immediate
                amqp.Publishing{
                        ContentType:   "text/plain",
                        CorrelationId: corrId,
                        ReplyTo:       q.Name,
                        Body:          []byte(n),
                })
        failOnError(err, "Failed to publish a message")

        for d := range msgs {
                if corrId == d.CorrelationId {
                        res = string(d.Body)
                        break
                }
        }

        return
}

func track(name string) func() {
        start := time.Now()
        return func() {
                fmt.Print("\n\n")
                log.Printf("%s, execution time %s\n", name, time.Since(start))
        }
}

func main() {
        defer track("main")()

        if os.Args != nil && len(os.Args) > 1 {
                //TRUE
                cmd := os.Args[1]
                rand.Seed(time.Now().UTC().UnixNano())
                log.Printf("CMD: %s", cmd)
                fmt.Printf("Asterisk*CLI>: %s\n", cmd)
                res, err := fibonacciRPC(cmd)
                failOnError(err, "Failed to handle RPC request")
                fmt.Print(res)
        } else {
                //FALSE
                fmt.Print("[*] Required argument. \n Use Example: core show help")
        }
}
