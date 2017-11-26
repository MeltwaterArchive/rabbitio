// Copyright © 2017 Meltwater
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rmq

import (
	"log"

	"github.com/streadway/amqp"
)

// NewPublisher creates and sets up a RabbitMQ Publisher
func NewPublisher(amqpURI, exchange, queue, tag string, prefetch int) *RabbitMQ {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Fatalf("writer failed to connect to Rabbit: %s", err)
		return nil
	}

	go func() {
		log.Printf("writer closing: %s", <-conn.NotifyClose(make(chan *amqp.Error)))
		log.Printf("writer blocked by rabbit: %v", <-conn.NotifyBlocked(make(chan amqp.Blocking)))
	}()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("writer failed to get a channel from Rabbit: %s", err)
		return nil
	}

	if err = channel.ExchangeDeclarePassive(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // noWait
		nil,      // arguments
	); err != nil {
		log.Fatalf("Exchange Declare: %s", err)
	}

	r := &RabbitMQ{
		conn:            conn,
		channel:         channel,
		exchange:        exchange,
		contentType:     "application/json",
		contentEncoding: "UTF-8",
	}
	log.Print("RabbitMQ connected: ", amqpURI)

	return r
}

// Publish Takes stream of messages and publish them to rabbit
func (r *RabbitMQ) Publish(messages chan Message, o Override) {
	for m := range messages {

		// override routingKey stored in Message with the executed options
		var routingKey string
		if o.RoutingKey != "#" {
			routingKey = o.RoutingKey
		} else {
			routingKey = m.RoutingKey
		}

		if err := r.channel.Publish(
			r.exchange,
			routingKey,
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				Headers:         m.Headers,
				ContentType:     r.contentType,
				ContentEncoding: r.contentEncoding,
				Body:            m.Body,
				DeliveryMode:    amqp.Persistent,
			},
		); err != nil {
			log.Fatalf("writer failed to write document to rabbit: %s", err)
		}
		r.Wg.Done()
	}
}
