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
	"github.com/marcusolsson/tui-go"
	"github.com/olekukonko/tablewriter"
	"os"

	//"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
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
		//displayStats()
		ui()
	},
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func ui() {

	theme := tui.NewTheme()
	theme.SetStyle("box.focused.border", tui.Style{Fg: tui.ColorYellow, Bg: tui.ColorDefault})

	memUsedUi := tui.NewTextEdit()
	memUsedUi.SetText(strconv.Itoa(helper.MemUsed))
	memUsedBox := tui.NewVBox(memUsedUi)
	memUsedBox.SetTitle("Memory Used")
	memUsedBox.SetBorder(true)
	memUsedBox.SetSizePolicy(tui.Minimum, tui.Minimum)

	memTotalUi := tui.NewTextEdit()
	memTotalUi.SetText(strconv.Itoa(helper.MemTotal))
	memTotalBox := tui.NewVBox(memTotalUi)
	memTotalBox.SetTitle("Memory Total")
	memTotalBox.SetBorder(true)
	memTotalBox.SetSizePolicy(tui.Minimum, tui.Minimum)

	memCachedUi := tui.NewTextEdit()
	memCachedUi.SetText(strconv.Itoa(helper.MemCached))
	memCachedBox := tui.NewVBox(memCachedUi)
	memCachedBox.SetTitle("Memory Cached")
	memCachedBox.SetBorder(true)
	memCachedBox.SetSizePolicy(tui.Minimum, tui.Minimum)

	memFreeUi := tui.NewTextEdit()
	memFreeUi.SetText(strconv.Itoa(helper.MemFree))
	memFreeBox := tui.NewVBox(memFreeUi)
	memFreeBox.SetTitle("Memory Free")
	memFreeBox.SetBorder(true)
	memFreeBox.SetSizePolicy(tui.Minimum, tui.Minimum)

	memoryBox := tui.NewHBox(memUsedBox,memTotalBox,memCachedBox,memFreeBox)
	memoryBox.SetTitle("Memory")
	memoryBox.SetBorder(true)
	//memoryBox.SetSizePolicy(tui.Preferred, tui.Minimum)

	cpuSystemUi := tui.NewTextEdit()
	cpuSystemUi.SetText(FloatToString(helper.CpuSystem))
	cpuSystemBox := tui.NewVBox(cpuSystemUi)
	cpuSystemBox.SetTitle("Cpu System")
	cpuSystemBox.SetBorder(true)

	cpuUserUi := tui.NewTextEdit()
	cpuUserUi.SetText(FloatToString(helper.CpuUser))
	cpuUserBox := tui.NewVBox(cpuUserUi)
	cpuUserBox.SetTitle("Cpu User")
	cpuUserBox.SetBorder(true)

	cpuIdleUi := tui.NewTextEdit()
	cpuIdleUi.SetText(FloatToString(helper.CpuIdle))
	cpuIdleBox := tui.NewVBox(cpuIdleUi)
	cpuIdleBox.SetTitle("Cpu Idle")
	cpuIdleBox.SetBorder(true)

	netRxUi := tui.NewTextEdit()
	netRxUi.SetText(strconv.Itoa(helper.RxBytes))
	netRxBox := tui.NewVBox(netRxUi)
	netRxBox.SetTitle("Rx Bytes")
	netRxBox.SetBorder(true)

	netTxUi := tui.NewTextEdit()
	netTxUi.SetText(strconv.Itoa(helper.TxBytes))
	netTxBox := tui.NewVBox(netTxUi)
	netTxBox.SetTitle("Tx Bytes")
	netTxBox.SetBorder(true)

	uptimeUi := tui.NewTextEdit()
	uptimeUi.SetText(helper.Uptime)
	uptimeBox := tui.NewVBox(uptimeUi)
	uptimeBox.SetTitle("Uptime")
	uptimeBox.SetBorder(true)

	timeUi := tui.NewTextEdit()
	timeUi.SetText(helper.Time)
	timeBox := tui.NewVBox(timeUi)
	timeBox.SetTitle("Last Updated")
	timeBox.SetBorder(true)

	cpuBox := tui.NewHBox(cpuUserBox,cpuSystemBox,cpuIdleBox)
	cpuBox.SetTitle("CPU")
	cpuBox.SetBorder(true)
	//cpuBox.SetSizePolicy(tui.Expanding, tui.Expanding)

	netBox := tui.NewHBox(netRxBox, netTxBox)
	netBox.SetTitle("Network")
	netBox.SetBorder(true)

	root := tui.NewHBox(memoryBox,cpuBox,netBox,uptimeBox,timeBox)

	ui, err := tui.New(root)
	if err != nil {
		panic(err)
	}
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		panic(err)
	}

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
