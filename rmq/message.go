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
	Body       []byte
	RoutingKey string
	Headers    map[string]interface{}
}

// NewMessageFromAttrs will create a new message from a byte slice and attributes
func NewMessageFromAttrs(bytes []byte, attrs map[string]string) *Message {

	// add header information to the Message
	var headers = make(map[string]interface{})
	var key string
	for k, v := range attrs {
		switch k {
		// use the provided routing key to override tarball configuration
		case "amqp.routingKey":
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
