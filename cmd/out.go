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
	"os/signal"
	"syscall"

	"github.com/meltwater/rabbitio/file"
	"github.com/meltwater/rabbitio/rmq"
	"github.com/spf13/cobra"
)

var (
	outputDirectory string
	batchSize       int
)

// outCmd represents the out command
var outCmd = &cobra.Command{
	Use:   "out",
	Short: "Consumes data out from RabbitMQ and stores to tarballs",
	Long: `Select your output directory and batchsize of the tarballs.
	When there are no more messages in the queue, press CTRL + c, to interrupt
	the consumption and save the last message buffers.`,
	Run: func(cmd *cobra.Command, args []string) {
		channel := make(chan rmq.Message, prefetch*2)
		verify := make(chan uint64)

		rabbit := rmq.NewConsumer(uri, exchange, queue, routingKey, tag, prefetch)
		path := file.NewOutput(outputDirectory, batchSize)

		go rabbit.Consume(channel, verify)

		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			fmt.Println(" Interruption, saving last memory bits..")
			close(channel)
		}()

		path.Receive(channel, verify)
	},
}

func init() {
	RootCmd.AddCommand(outCmd)

	outCmd.Flags().StringVarP(&outputDirectory, "directory", "d", ".", "Output directory for files consumed from RabbitMQ")
	outCmd.Flags().IntVarP(&batchSize, "batch", "b", 1000, "Number of messages stored in each tarball")
}
