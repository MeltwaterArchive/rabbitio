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
	"time"

	"github.com/streadway/amqp"
)

// NewConsumer creates and sets up a RabbitMQ struct best used for consuming messages
func NewConsumer(amqpURI, exchange, queue, routingKey, tag string, prefetch int) *RabbitMQ {
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

	q, err := channel.QueueDeclarePassive(
		queue, // name of the queue
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Queue Declare: %s", err)
	}
	if q.Messages == 0 {
		log.Fatalf("No messages in RabbitMQ Queue: %s", q.Name)
	}
	if err = channel.QueueBind(
		q.Name,     // name of the queue
		routingKey, // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		log.Fatalf("Queue Bind: %s", err)
	}

	r := &RabbitMQ{
		conn:            conn,
		channel:         channel,
		exchange:        exchange,
		contentType:     "application/json",
		contentEncoding: "UTF-8",
		tag:             tag,
	}
	log.Print("RabbitMQ connected: ", amqpURI)
	log.Printf("Bind to Exchange: %q and Queue: %q, Messaging waiting: %d", exchange, queue, q.Messages)

	return r
}

func (r *RabbitMQ) ackMultiple(deliveryTag <-chan Verify) {
	for v := range deliveryTag {
		err := r.channel.Ack(v.Tag, v.MultiAck)
		if err != nil {
			log.Fatalf("rabbit channel ack failed %s", err)
		}
	}
}

func (r *RabbitMQ) SafeStop() {
	err := r.channel.Cancel(
		r.tag, // name
		false, // noWait
	)
	if err != nil {
		log.Fatalf("rabbit channel cancel failed %s", err)
	} else {
		for _ = range time.Tick(2 * time.Second) {
			log.Print("Rabbit Channel canceled")
			break
		}
	}
}

// Consume outputs a stream of Message into a channel from rabbit
func (r *RabbitMQ) Consume(out chan Message, verify <-chan Verify) {
	go r.ackMultiple(verify)

	// set up a channel consumer
	deliveries, err := r.channel.Consume(
		r.queue, // name
		r.tag,   // consumer
		false,   // noAck
		false,   // exclusive
		false,   // noLocal
		false,   // noWait
		nil,     // args
	)
	if err != nil {
		log.Fatalf("rabbit consumer failed %s", err)
	}

	// process deliveries from the queue
	for d := range deliveries {
		// create a new Message for the rabbit message
		msg := Message{
			Body:        d.Body,
			RoutingKey:  d.RoutingKey,
			Headers:     d.Headers,
			DeliveryTag: d.DeliveryTag,
		}
		// write Message to channel
		out <- msg
	}

	log.Print("All messages consumed")
}
