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
	"fmt"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
)

// Message contains the most basic about the message
type Message struct {
	Body        []byte
	RoutingKey  string
	Headers     amqp.Table
	DeliveryTag uint64
}

// Verify will be used to Ack Message from the queue
type Verify struct {
	Tag      uint64
	MultiAck bool
}

// ToPAXRecords takes amqp headers and convert them to PAXRecords compatible
func (m *Message) ToPAXRecords() map[string]string {

	pax := make(map[string]string)
	var headerType string

	for k, v := range m.Headers {
		switch v.(type) {
		case int, int32, int64:
			headerType = "int"
		case float32, float64:
			headerType = "float"
		case bool:
			headerType = "bool"
		case string:
			headerType = "string"
		}
		pax[fmt.Sprintf("RABBITIO.amqp.Headers.%s.%s", headerType, k)] = fmt.Sprintf("%v", v)
	}
	return pax
}

// NewMessage will create a new message from a byte slice and attributes
func NewMessage(bytes []byte, xattr map[string]string) *Message {

	// add amqp header information to the Message
	var headers = make(amqp.Table)
	var routingKey string

	// need to support more than just string here for v
	for k, v := range xattr {

		switch {
		case k == "RABBITIO.amqp.routingKey":
			routingKey = v
		case strings.HasPrefix(k, "RABBITIO.amqp.Headers."):
			// th is now [type, header]
			th := strings.SplitN(strings.TrimPrefix(k, "RABBITIO.amqp.Headers."), ".", 2)
			headerType := th[0]
			header := strings.Join(th[1:], ".")

			switch headerType {
			case "bool":
				if b, err := strconv.ParseBool(v); err == nil {
					headers[header] = b
				}
			case "int":
				if i, err := strconv.ParseInt(v, 10, 64); err == nil {
					headers[header] = i
				}
			case "float":
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					headers[header] = f
				}
			case "string":
				headers[header] = v
			}
		}
	}

	// create a message
	m := &Message{
		Body:       bytes,
		RoutingKey: routingKey,
		Headers:    headers,
	}

	return m
}
