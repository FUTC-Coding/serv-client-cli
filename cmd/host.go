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
	"github.com/spf13/cast"
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

var cmdCpu = &cobra.Command{
	Use:   "cpu",
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
		termuiRunCpu(args[0])
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
			updateNetData()
			drawFunction()
		}
	}
}

func termuiRunCpu(args string) {
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
			updateCpuData() //todo add another stats from hostname which only gets cpu data to save bandwidth
			drawCpu()
		}
	}
}

var cpuUsageArray = make([][]float64, 1)
//appends new data and deletes oldest datapoint
func updateCpuData() {
	if len(cpuUsageArray) == 0 {
		cpuUsageArray[0] = []float64{}
	}
	//delete newest datapoint if there are more than 40 datapoints already
	if len(cpuUsageArray) > 40 {
		//todo deleting of old datapoint doesnt work
		cpuUsageArray[0] = append(cpuUsageArray[0][:0], cpuUsageArray[0][1:]...)
	}
	//append new datapoint
	cpuUsageArray[0] = append(cpuUsageArray[0], helper.CpuPerc)
	//fmt.Println(cpuUsageArray)
}

var rxData = make([]float64, 1)
var txData = make([]float64, 1)
func updateNetData() {
	if len(rxData) > 40 {
		rxData = append(rxData[:0], rxData[1:]...)
		txData = append(txData[:0], rxData[1:]...)
	}
	//append new datapoint
	rxData = append(rxData, cast.ToFloat64(helper.RxBytes))
	txData = append(txData, cast.ToFloat64(helper.TxBytes))
}

//draw the UI elements
func drawFunction() {
	//draw memory barchart
	bc := widgets.NewBarChart()
	bc.Data = []float64{float64(helper.MemFree), float64(helper.MemCached), float64(helper.MemTotal), float64(helper.MemUsed)}
	bc.Labels = []string{"Free", "Cached", "Total", "Used"}
	bc.Title = "Memory (MB)"
	bc.SetRect(0, 0, 50, 10)
	bc.BarWidth = 10
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	//draw CPU usage plot
	plot := widgets.NewPlot()
	plot.Marker = widgets.MarkerDot
	plot.Data = cpuUsageArray
	plot.AxesColor = ui.ColorWhite
	plot.LineColors[0] = ui.ColorYellow
	plot.SetRect(50,0,100,10)
	//plot.DrawDirection = widgets.DrawLeft
	plot.Title = "CPU (%)"

	//draw network usage sparkline
	rxSpark := widgets.NewSparkline()
	rxSpark.Data = rxData
	rxSpark.LineColor = ui.ColorGreen

	txSpark := widgets.NewSparkline()
	txSpark.Data = txData
	txSpark.LineColor = ui.ColorBlue

	sparkGroup := widgets.NewSparklineGroup(rxSpark,txSpark)
	sparkGroup.Title = "Net (Bytes)"
	sparkGroup.SetRect(0, 10,50,20)

	table := widgets.NewTable()
	table.Rows = [][]string {
		{"uptime", "time of data"},
		{helper.Uptime, helper.Time},
	}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(50,10,100,15)

	pie := widgets.NewPieChart()
	pie.Title = "Disk Usage(MB)"
	pie.SetRect(50,15,71,25)
	pie.Data = []float64{float64(helper.DiskUsed), float64(helper.DiskFree)}
	//pie.AngleOffset = -.9 * math.Pi
	pie.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%0.2f", v)
	}

	ui.Render(bc, plot, sparkGroup, table, pie)
}

func drawCpu() {
	//draw CPU usage plot
	plot := widgets.NewPlot()
	plot.Marker = widgets.MarkerDot
	plot.Data = cpuUsageArray
	plot.AxesColor = ui.ColorWhite
	plot.LineColors[0] = ui.ColorYellow
	plot.SetRect(0,0,100,20)
	//plot.DrawDirection = widgets.DrawLeft
	plot.Title = "CPU (%)"
	ui.Render(plot)
}

func displayStats() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"CpuUser","CpuSystem", "CpuIdle", "MemTotal", "MemUsed", "MemCached", "MemFree", "RxBytes", "TxBytes", "Uptime", "Time"})
	table.SetBorder(false)
	//append a row with the data from the last call of StatsFromHostname()
	table.Append([]string{
		FloatToString(helper.CpuPerc),
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
	hostCmd.AddCommand(cmdCpu)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hostCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
