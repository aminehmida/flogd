/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/aminehmida/flogd/matcher"
	"github.com/aminehmida/flogd/tailer"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Start monitoring log stream(s) defined in the command line or in a config file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// logging.Info("monitor called")
		log.Info().Msg("monitor called")
		stype, _ := cmd.Flags().GetString("type")

		regex, _ := cmd.Flags().GetString("regex")

		command, err := cmd.Flags().GetString("do")
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Can not get do argument: %v", err))
			return
		}
		if command == "" {
			log.Warn().Msg("\"do\" argument not defined. Will not execute any command on match")
		}

		count, _ := cmd.Flags().GetInt("count")
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Can not get count argument: %v", err))
			return
		}
		interval, err := cmd.Flags().GetInt("interval")
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Can not get interval argument: %v", err))
			return
		}
		config, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Can not get config argument: %v", err))
			return
		}
		if config == "" {
			if stype == "process" {
				tailerLineOutPipe := make(chan string)
				tailerErrPipe := make(chan error)
				matcherOutPipe := make(chan string)

				go func() {
					for err := range tailerErrPipe {
						if err != nil {
							log.Error().Msg(err.Error())
						}
						close(matcherOutPipe)
					}
				}()

				go tailer.ProcessTailer(args[0], tailerLineOutPipe, tailerErrPipe)
				go matcher.Monitor(regex, count, interval, tailerLineOutPipe, matcherOutPipe)

				if strings.Contains(command, "%s") {
					for match := range matcherOutPipe {
						c := fmt.Sprintf(command, match)
						log.Info().Msg(fmt.Sprintf(" ==> Executing: %s", c))
					}
				} else if command != "" {
					for range matcherOutPipe {
						log.Info().Msg(fmt.Sprintf(" ==> Executing: %s. ", command))
					}
				} else {
					log.Info().Msg(" ==> No command to execute")
				}
			} else {
				log.Error().Msg("Stream type not supported")
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
