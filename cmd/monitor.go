/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/aminehmida/flogd/tailer"
	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Start monitoring log stream(s) defined in the command line or in a config file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("monitor called")
		stype, _ := cmd.Flags().GetString("type")

		regex, _ := cmd.Flags().GetString("regex")

		command, err := cmd.Flags().GetString("do")
		if err != nil {
			fmt.Println("Can not get do argument:", err)
			return
		}
		if command == "" {
			fmt.Println("do argument not defined. Will not execute any command on match")
		}

		count, _ := cmd.Flags().GetInt("count")
		if err != nil {
			fmt.Println("Can not get count argument:", err)
			return
		}
		interval, err := cmd.Flags().GetInt("interval")
		if err != nil {
			fmt.Println("Can not get interval argument:", err)
			return
		}
		config, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Println("Can not get config argument:", err)
			return
		}
		if config == "" {
			if stype == "process" {
				fmt.Println("monitor process")
				tailerOutPipe := make(chan string)
				tailerErrPipe := make(chan error)

				matcherInPipe := make(chan string)
				matcherOutPipe := make(chan string)

				tailer.ProcessTailer(command, tailerOutPipe, tailerErrPipe)

			} else {
				fmt.Println("Stream type not supported")
			}
		} else {
			fmt.Println("Using config file:", config)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
