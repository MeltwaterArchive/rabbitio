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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version                               string
	uri, exchange, queue, tag, routingKey string
	prefetch                              int
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "rio",
	Short: "Rabbit IO will help backup and restore your messages in RabbitMQ",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ver string) {
	version = ver
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&uri, "uri", "u", "amqp://guest:guest@localhost:5672/", "AMQP URI, uri to for instance RabbitMQ")
	RootCmd.PersistentFlags().StringVarP(&exchange, "exchange", "e", "", "Exchange to connect to")
	RootCmd.PersistentFlags().StringVarP(&queue, "queue", "q", "", "Queue to connect to")
	RootCmd.PersistentFlags().StringVarP(&routingKey, "routingkey", "r", "#", "Routing Key, if specified will override tarball routing key configuration")
	RootCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "Rabbit IO Connector "+version, "AMQP Client Tag")
	RootCmd.PersistentFlags().IntVarP(&prefetch, "prefetch", "p", 100, "Prefetch for batches")
}
