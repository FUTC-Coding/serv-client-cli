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
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"serv-client-cli/helper"
	"strconv"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hostname: " + args[0])
		helper.StatsFromHostname(args[0])
		fmt.Println(helper.CpuUser)
		displayStats()
	},
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func displayStats() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"CpuUser","CpuSystem", "CpuIdle", "MemTotal", "MemUsed", "MemCached", "MemFree", "RxBytes", "TxBytes", "Uptime", "Time"})
	table.SetBorder(false)
	//append a row with the data from the last call of StatsFromHostname()
	table.Append([]string{
		FloatToString(helper.CpuUser),
		FloatToString(helper.CpuSystem),
		FloatToString(helper.CpuIdle),
		strconv.Itoa(helper.MemTotal),
		strconv.Itoa(helper.MemUsed),
		strconv.Itoa(helper.MemCached),
		strconv.Itoa(helper.MemFree),
		strconv.Itoa(helper.RxBytes),
		strconv.Itoa(helper.TxBytes),
		helper.Uptime,
		helper.Time})

	table.Render() //output the table
}

func init() {
	rootCmd.AddCommand(hostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hostCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
