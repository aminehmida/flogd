/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aminehmida/flogd/exec"
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
		log.Debug().Msg("monitor called")
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

				var wg sync.WaitGroup

				go func() {
					for err := range tailerErrPipe {
						if err != nil {
							log.Error().Msg(err.Error())
						} else {
							log.Info().Msg("Tailer finished")
						}
						wg.Wait()
						close(matcherOutPipe)
						log.Debug().Msg("Closed matcherOutPipe")
					}
				}()

				go tailer.ProcessTailer(args[0], tailerLineOutPipe, tailerErrPipe)
				go matcher.Monitor(regex, count, interval, wg, tailerLineOutPipe, matcherOutPipe)

				if strings.Contains(command, "%s") {
					for match := range matcherOutPipe {
						wg.Add(1)
						go func(m string) {
							defer wg.Done()
							c := fmt.Sprintf(command, m)
							log.Info().Msg(fmt.Sprintf(" ==> Executing: %s", c))
							stdout, stderr, retCode, err := exec.Execute(c)
							if err != nil {
								log.Error().Msg(fmt.Sprintf("Error executing command: %v", err))
							} else {
								multulineInfoPrefixPrint(c+"; stdout", stdout)
								multulineInfoPrefixPrint(c+"; stderr", stderr)
								log.Info().Msg(fmt.Sprintf("Command returned: %d", retCode))
							}
							log.Debug().Msg("Execution finish for: " + c)

						}(match)
					}
				} else if command != "" {
					for range matcherOutPipe {
						wg.Add(1)
						go func() {
							defer wg.Done()
							log.Info().Msg(fmt.Sprintf(" ==> Executing: %s. ", command))
							stdout, stderr, retCode, err := exec.Execute(command)
							if err != nil {
								log.Error().Msg(fmt.Sprintf("Error executing command: %v", err))
							} else {
								multulineInfoPrefixPrint(command+"; stdout", stdout)
								multulineInfoPrefixPrint(command+"; stderr", stderr)
								log.Info().Msg(fmt.Sprintf("Command returned: %d", retCode))
							}
							log.Debug().Msg("Execution finish for: " + command)
						}()
					}
				} else {
					for match := range matcherOutPipe {
						wg.Add(1)
						// log.Info().Msg(fmt.Sprintf(" ==> Match: %s", match))
						go func(m string) {
							defer wg.Done()
							log.Info().Msg("Match found: " + m)
						}(match)
					}
				}
			} else {
				log.Error().Msg("Stream type not supported")
			}
		} else {
			fmt.Println("Using config file:", config)
		}
	},
}

func monitorCommand(command, do, regexp string, count, interval int) {
	
})

func multulineInfoPrefixPrint(prefix, s string) {
	for _, line := range strings.Split(s, "\n") {
		log.Info().Msg(fmt.Sprintf("%s: %s", prefix, line))
	}
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
