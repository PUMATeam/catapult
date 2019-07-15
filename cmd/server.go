// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"github.com/PUMATeam/catapult/api"
	"github.com/spf13/cobra"
)

var port int

// restCmd represents the rest command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start catapult server",
	Long:  `Start catapult server`,
	Run: func(cmd *cobra.Command, args []string) {
		handler := api.Bootstrap(port)
		api.Start(handler)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVarP(&port, "port", "p", 8888, "Port for which to listen")
}
