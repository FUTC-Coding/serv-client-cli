/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"serv-client-cli/helper"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		listTables(helper.ListOfTables())
	},
}

func listTables(tables []string) {
	var data [][]string

	for i := 0 ; i < len(tables); i++ {
		helper.StatsFromHostname(tables[i])
		memPerc := (float64(helper.MemUsed)/float64(helper.MemTotal))*100
		storagePerc := (float64(helper.DiskUsed)/(float64(helper.DiskFree)+float64(helper.DiskUsed)))*100 //calculate percentage
		s := []string{tables[i], FloatToString(helper.CpuPerc), FloatToString(memPerc), FloatToString(storagePerc), helper.Uptime, helper.Time}
		data = append(data, s)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{"Server", "Cpu %", "Mem %", "Storage %", "Uptime", "Time Updated"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
