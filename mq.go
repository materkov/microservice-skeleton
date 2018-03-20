package test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type userCreated struct {
	ID string
}

func onUserCreated(body interface{}) error {
	e := body.(userCreated)
	log.Printf("created user %s", e.ID)
	return nil
}

func userCreatedDecode(body []byte) (interface{}, error) {
	r := userCreated{}
	err := json.Unmarshal(body, &r)
	return r, err
}

// ServeMQ runs queue listener
func ServeMQ() {
	handle("UserCreated", "FetchAvatar", onUserCreated, userCreatedDecode)
}

type mqHandler func(body interface{}) error
type mqBodyDecode func(body []byte) (interface{}, error)

func handle(event, queueName string, handler mqHandler, decoder mqBodyDecode) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("error dialing mq: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("error creating channel: %s", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error declaring queue: %s", err)
	}

	err = ch.ExchangeDeclare(event, "fanout", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error declaring exchange: %s", err)
	}

	err = ch.QueueBind(queueName, "", event, false, nil)
	if err != nil {
		return fmt.Errorf("error binding queue to exchange: %s", err)
	}

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error consuming messages: %s", err)
	}

	for d := range msgs {
		req, err := decoder(d.Body)
		if err != nil {
			log.Printf("error decoding body: %s", err)
			continue
		}

		err = handler(req)
		if err != nil {
			log.Printf("error handling: %s", err)
		}
	}
	return nil
}
