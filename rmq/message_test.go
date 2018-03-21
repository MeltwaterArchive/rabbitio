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

var (
	myStringHeader  = "amqp.Headers.string.myStringHeader"
	myInt32Header   = "amqp.Headers.int.myInt32Header"
	myInt64Header   = "amqp.Headers.int.myInt64Header"
	myFloat32Header = "amqp.Headers.float.myFloat32Header"
	myFloat64Header = "amqp.Headers.float.myFloat64Header"
	myBoolHeader    = "amqp.Headers.bool.myBoolHeader"
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
	attrHeaders[myStringHeader] = "myString"
	attrHeaders[myInt32Header] = "32"
	attrHeaders[myInt64Header] = "64"
	attrHeaders[myFloat32Header] = "32.32"
	attrHeaders[myFloat64Header] = "64.64"
	attrHeaders[myBoolHeader] = "true"

	attrs := message.ToXAttrs()

	assert.NoError(t, messageHeaders.Validate(), "should be valid Headers")

	assert.Equal(t, attrHeaders[myStringHeader], attrs[myStringHeader])
	assert.Equal(t, attrHeaders[myInt32Header], attrs[myInt32Header])
	assert.Equal(t, attrHeaders[myInt64Header], attrs[myInt64Header])
	assert.Equal(t, attrHeaders[myFloat32Header], attrs[myFloat32Header])
	assert.Equal(t, attrHeaders[myFloat64Header], attrs[myFloat64Header])
	assert.Equal(t, attrHeaders[myBoolHeader], attrs[myBoolHeader])
}

func TestNewMessage(t *testing.T) {
	var headers = make(map[string]string)
	headers["amqp.routingKey"] = "routingKey from tarball XAttrs"
	headers[myStringHeader] = "myString"
	headers[myInt32Header] = "3232"
	headers[myInt64Header] = "6464"
	headers[myFloat32Header] = "32.123"
	headers[myFloat64Header] = "64.123"
	headers[myBoolHeader] = "true"

	m := NewMessage([]byte("Message"), headers)

	assert.Equal(t, "routingKey from tarball XAttrs", m.RoutingKey)
	assert.Equal(t, []byte("Message"), m.Body)
	assert.NoError(t, m.Headers.Validate())
}
