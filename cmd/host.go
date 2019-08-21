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
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"time"

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
		//ui(args[0])
		//displayStats()
		termuiRun(args[0])
	},
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

type MemWidget struct {
	*widgets.BarChart
}

type MemoryInfo struct {
	Total int
	Cached int
	Free int
	Used int
}

func (self *MemWidget) renderMemInfo(info MemoryInfo) {
	self.Data = append(self.Data)
}

func termuiRun(args string) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second * 6).C //ticker of 6 seconds
	for {
		select {
			case e := <-uiEvents:
				switch e.ID {
				case "q", "<C-c>":
					return
				}
		case <- ticker: //run following things every 6 seconds to update stats and widgets
			helper.StatsFromHostname(args)
			updateCpuData()
			drawFunction()
		}
	}
}

var cpuUsageArray = make([][]float64, 2)
//appends new data and deletes oldest datapoint
func updateCpuData() {
	if len(cpuUsageArray) == 0 {
		cpuUsageArray[0] = []float64{}
		cpuUsageArray[1] = []float64{}
	}
	//delete newest datapoint if there are more than 10 datapoints already
	/*
	if len(cpuUsageArray) > 10 {
		cpuUsageArray = append(cpuUsageArray[:0], cpuUsageArray[1:]...)
	} */
	//append new datapoint
	cpuUsageArray[0] = append(cpuUsageArray[0], helper.CpuSystem)
	cpuUsageArray[1] = append(cpuUsageArray[1], helper.CpuUser)
	//fmt.Println(cpuUsageArray)
}

//draw the UI elements
func drawFunction() {
	bc := widgets.NewBarChart()
	bc.Data = []float64{float64(helper.MemFree), float64(helper.MemCached), float64(helper.MemTotal), float64(helper.MemUsed)}
	bc.Labels = []string{"Free", "Cached", "Total", "Used"}
	bc.Title = "Memory (MB)"
	bc.SetRect(0, 0, 40, 8)
	bc.BarWidth = 8
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	line := widgets.NewPlot()
	line.Data = cpuUsageArray
	line.Marker = widgets.MarkerDot
	line.DotMarkerRune = '+'
	line.AxesColor = ui.ColorWhite
	line.LineColors[0] = ui.ColorYellow
	line.DrawDirection = widgets.DrawLeft
	line.Title = "CPU (%)"
	line.SetRect(0,0,50,10)
	ui.Render(bc,line)
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
