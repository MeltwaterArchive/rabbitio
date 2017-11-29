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

// Message contains the most basic about the message
type Message struct {
	Body        []byte
	RoutingKey  string
	Headers     map[string]interface{}
	DeliveryTag uint64
}

// Verify will be used to Ack Message from the queue
type Verify struct {
	Tag      uint64
	MultiAck bool
}

// NewMessageFromAttrs will create a new message from a byte slice and attributes
func NewMessageFromAttrs(bytes []byte, attrs map[string]string) *Message {

	// add amqp header information to the Message
	var headers = make(map[string]interface{})
	var key string

	// need to support more than just string here for v
	for k, v := range attrs {
		switch k {
		// use the routing key from tarball header configuration
		case "amqp.routingKey":
			key = v
		default:
			headers[k] = v
		}
	}

	// create a message
	m := &Message{
		Body:       bytes,
		RoutingKey: key,
		Headers:    headers,
	}

	return m
}
