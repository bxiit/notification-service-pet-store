package main

import (
	"encoding/json"
	"github.com/bxiit/notification-service-pet-store/config"
	"github.com/bxiit/notification-service-pet-store/data/dto"
	"github.com/bxiit/notification-service-pet-store/mailer"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	smtpConfig := config.LoadConfig()
	m := mailer.New(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password, smtpConfig.Sender)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"order", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
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

	var forever chan struct{}

	go func() {
		for d := range msgs {
			body := d.Body
			log.Printf("received message! %s", body)
			data := make(map[string]json.RawMessage)
			err := json.Unmarshal(body, &data)
			if err != nil {
				log.Fatalf("Failed to unmarshal JSON: %v", err)
				return
			}
			var orderDTO dto.OrderDTO
			err = json.Unmarshal(data["order_info"], &orderDTO)
			if err != nil {
				log.Printf("failed to unmarshal order info")
				return
			}

			var userDTO dto.UserDTO
			err = json.Unmarshal(data["user_info"], &userDTO)
			if err != nil {
				log.Printf("failed to unmarshal user info")
				return
			}
			messageData := map[string]any{
				"itemName":  orderDTO.Item.Name,
				"username":  userDTO.Username,
				"itemImage": orderDTO.Item.ImageURL,
			}
			err = m.Send(userDTO.Email, "user_welcome.tmpl", messageData)
			if err != nil {
				log.Printf("failed to send mail %v", err)
				return
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
