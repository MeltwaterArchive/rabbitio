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
	"errors"
	"sync"

	"github.com/meltwater/rabbitio/file"
	"github.com/meltwater/rabbitio/rmq"
	"github.com/spf13/cobra"
)

var (
	fileInput string
)

// inCmd represents the in command
var inCmd = &cobra.Command{
	Use:   "in",
	Short: "Publishes documents from tarballs into RabbitMQ exchange",
	Long:  `Specify a directory or file and tarballs will be published.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if fileInput == "" {
			return errors.New("please specify a tarball or directory with tarballs using the -f flag")
		}
		channel := make(chan rmq.Message, prefetch)

		override := rmq.Override{RoutingKey: routingKey}
		path, err := file.NewInput(fileInput)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		path.Wg = &wg
		rabbit := rmq.NewPublisher(uri, exchange, queue, tag, prefetch)
		rabbit.Wg = &wg
		defer rabbit.Close()

		go rabbit.Publish(channel, override)

		return path.Send(channel)
	},
}

func init() {
	RootCmd.AddCommand(inCmd)
	inCmd.Flags().StringVarP(&fileInput, "file", "f", "", "File is specified as either file or directory to restore into RabbitMQ")
}
