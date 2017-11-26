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
	"sync"

	"github.com/streadway/amqp"
)

// RabbitMQ type for talking to RabbitMQ
type RabbitMQ struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	override        Override
	exchange        string
	contentType     string
	contentEncoding string
	queue           string
	tag             string
	prefetch        int
	consume         bool
	publish         bool
	Wg              *sync.WaitGroup
}

// Override will be used to override RabbitMQ settings on publishing messages
type Override struct {
	RoutingKey string
}
