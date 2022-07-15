/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flogd",
	Short: "iFLogDo is a tool to monitor log files or commands output and execute commands when a specific pattern is found",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.flogd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	// flogd monitor
	// -t process 'docker logs container-name -f'
	// -r 'Login failed for (?P<ip>d\.\d\.d\.\d)'
	// -n 10
	// -i 5s
	// -u ip
	// -d 'echo ip found {{ip}}'
	rootCmd.PersistentFlags().StringP("type", "t", "process", "type of the log stream to monitor")
	rootCmd.PersistentFlags().StringP("regex", "r", "*", "regex to match")
	rootCmd.PersistentFlags().IntP("count", "n", 10, "number of times to match before triggering the command")
	rootCmd.PersistentFlags().IntP("interval", "i", 5, "interval to check the log stream for matches in seconds")
	rootCmd.PersistentFlags().StringP("do", "d", "", "command to execute when a match is found <count> times in <interval> time interval")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
}
