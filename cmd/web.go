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
	"net/http"

	"github.com/gin-gonic/gin"
	rh "github.com/michaelklishin/rabbit-hole"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(webCmd)
}

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Prints the web of Rabbit IO",
	Run: func(cmd *cobra.Command, args []string) {
		rmqc, _ := rh.NewClient("http://127.0.0.1:15672", "guest", "guest")
		// for i, queue := range qs {
		// 	fmt.Println(i, queue.Name, queue.Messages)
		// }

		router := gin.Default()
		router.LoadHTMLGlob("templates/*")
		router.GET("/", func(c *gin.Context) {
			qs, _ := rmqc.ListQueuesIn("/")
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"title":  "Queues",
				"queues": qs,
			})
		})
		router.Run(":9010")
	},
}
