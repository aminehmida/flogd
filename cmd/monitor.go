package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aminehmida/flogd/config"
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
		do, _ := cmd.Flags().GetString("do")
		count, _ := cmd.Flags().GetInt("count")
		interval, _ := cmd.Flags().GetInt("interval")
		configFile, _ := cmd.Flags().GetString("config")

		if configFile == "" {
			if stype != "process" {
				log.Error().Msg("Stream type not supported")
				return
			}
			if len(args) == 0 {
				log.Error().Msg("A command to monitor is required when no config file is provided")
				return
			}
			monitorCommand(args[0], do, regex, count, interval, nil)
		} else {
			log.Debug().Msg("Using config file: " + configFile)
			var configs config.Configs
			// read config file as []byte
			configBytes, err := os.ReadFile(configFile)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("Can not read config file: %v", err))
				return
			}
			err = configs.Decode(configBytes)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("Error loading config file: %v", err))
				return
			}
			var wg sync.WaitGroup
			for _, stream := range configs {
				if stream.SType == "process" {
					log.Info().Msg("Monitoring process: " + stream.Name)
					wg.Add(1)
					go monitorCommand(stream.Command, stream.Do, stream.Regex, stream.Count, stream.Interval, &wg)
				}
			}
			wg.Wait()
		}
	},
}

func monitorCommand(command, do, regex string, count, interval int, mainWg *sync.WaitGroup) {
	if do == "" {
		log.Warn().Msg("\"do\" argument not defined. Will not execute any command on match")
	}
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

	go tailer.ProcessTailer(command, tailerLineOutPipe, tailerErrPipe)
	go matcher.Monitor(regex, count, interval, tailerLineOutPipe, matcherOutPipe, &wg)

	if strings.Contains(do, "%s") {
		for match := range matcherOutPipe {
			wg.Add(1)
			go func(m string) {
				defer wg.Done()
				c := fmt.Sprintf(do, m)
				log.Info().Msg(fmt.Sprintf(" ==> Executing: %s", c))
				stdout, stderr, retCode, err := exec.Execute(c)
				if err != nil {
					log.Error().Msg(fmt.Sprintf("Error executing command: %v", err))
				} else {
					multilineInfoPrefixPrint(c+"; stdout", stdout)
					multilineInfoPrefixPrint(c+"; stderr", stderr)
					log.Info().Msg(fmt.Sprintf("Command returned: %d", retCode))
				}
				log.Debug().Msg("Execution finish for: " + c)

			}(match)
		}
	} else if do != "" {
		for range matcherOutPipe {
			wg.Add(1)
			go func() {
				defer wg.Done()
				log.Info().Msg(fmt.Sprintf(" ==> Executing: %s. ", do))
				stdout, stderr, retCode, err := exec.Execute(do)
				if err != nil {
					log.Error().Msg(fmt.Sprintf("Error executing command: %v", err))
				} else {
					multilineInfoPrefixPrint(do+"; stdout", stdout)
					multilineInfoPrefixPrint(do+"; stderr", stderr)
					log.Info().Msg(fmt.Sprintf("Command returned: %d", retCode))
				}
				log.Debug().Msg("Execution finish for: " + do)
			}()
		}
	} else {
		for match := range matcherOutPipe {
			wg.Add(1)
			go func(m string) {
				defer wg.Done()
				log.Info().Msg("Match found: " + m)
			}(match)
		}
	}
	if mainWg != nil {
		mainWg.Done()
	}
}

func multilineInfoPrefixPrint(prefix, s string) {
	for _, line := range strings.Split(s, "\n") {
		log.Info().Msg(fmt.Sprintf("%s: %s", prefix, line))
	}
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
