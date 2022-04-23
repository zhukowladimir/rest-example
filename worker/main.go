package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type QueueMsg struct {
	ID        string `json:"ID"`
	Username  string `json:"Username"`
	Avatar    string `json:"Avatar"`
	Sex       string `json:"Sex"`
	Email     string `json:"Email"`
	WinCount  int    `json:"WinCount"`
	LossCount int    `json:"LossCount"`
	Duration  int    `json:"Duration"`
	Filename  string `json:"Filename"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	wq, err := ch.QueueDeclare(
		"work_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	rq, err := ch.QueueDeclare(
		"reply_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		wq.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [+] Received a message: %s", d.Body)

			qMsg := QueueMsg{}
			json.Unmarshal(d.Body, &qMsg)
			timer := time.Now()
			pdfBytes, err := createPdf(&qMsg)
			if err != nil {
				log.Println("[o] Error:", err.Error())
			}

			err = ch.Publish(
				"",
				rq.Name,
				false,
				false,
				amqp.Publishing{
					DeliveryMode:  amqp.Persistent,
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					MessageId:     d.MessageId,
					Body:          []byte(pdfBytes),
				},
			)
			failOnError(err, "Failed to publish a message")

			log.Println(" [+] Done with", time.Since(timer).Seconds())
		}
	}()

	log.Print(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
