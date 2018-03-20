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
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestToXAttrs(t *testing.T) {
	messageHeaders := make(amqp.Table)
	messageHeaders["myStringHeader"] = "myString"
	messageHeaders["myInt32Header"] = int32(32)
	messageHeaders["myInt64Header"] = int64(64)
	messageHeaders["myFloat32Header"] = float32(32.32)
	messageHeaders["myFloat64Header"] = float64(64.64)
	messageHeaders["myBoolHeader"] = true
	message := &Message{Headers: messageHeaders}

	var attrHeaders = make(map[string]string)
	attrHeaders["amqp.Headers.string.myStringHeader"] = "myString"
	attrHeaders["amqp.Headers.int.myInt32Header"] = "32"
	attrHeaders["amqp.Headers.int.myInt64Header"] = "64"
	attrHeaders["amqp.Headers.float.myFloat32Header"] = "32.32"
	attrHeaders["amqp.Headers.float.myFloat64Header"] = "64.64"
	attrHeaders["amqp.Headers.bool.myBoolHeader"] = "true"

	attrs := message.ToXAttrs()

	assert.Equal(t, attrHeaders, attrs)
	assert.NoError(t, messageHeaders.Validate())
}

func TestNewMessage(t *testing.T) {
	var headers = make(map[string]string)
	headers["amqp.routingKey"] = "routingKey from tarball XAttrs"
	headers["amqp.Headers.string.myHeader"] = "myString"
	headers["amqp.Headers.int.myIntHeader"] = "456"
	headers["amqp.Headers.float.myFloatHeader"] = "123.123"
	headers["amqp.Headers.bool.myBoolHeader"] = "true"

	m := NewMessage([]byte("Message"), headers)

	assert.Equal(t, "routingKey from tarball XAttrs", m.RoutingKey)
	assert.Equal(t, []byte("Message"), m.Body)
	assert.NoError(t, m.Headers.Validate())
}
