package main

import (
	"fmt"
	"log"
	"os/exec"
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

func track(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s, execution time %s\n", name, time.Since(start))
	}
}

func main() {
	conn, err := amqp.Dial("amqp://" + rbUser + ":" + rbPass + "@" + rbHost + ":" + rbPort + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {

			var parameter string = string(d.Body)
			log.Printf("Asterisk: %s", parameter)

			out, err := AsteriskCMD(parameter)
			failOnError(err, "asterisk not fount")

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          []byte(out),
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Asterisk Server Remote CLI")
	<-forever
}

func AsteriskCMD(command string) (string, error) {
	defer track("AsteriskCMD")()
	prg := "asterisk"

	arg1 := "-rx"
	arg2 := command

	cmd := exec.Command(prg, arg1, arg2)
	stdout, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return string(stdout), nil
}
