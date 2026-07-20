package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flogd",
	Short: "flogd monitors log files or command output and runs a command when a pattern is matched",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("type", "t", "process", "type of the log stream to monitor")
	rootCmd.PersistentFlags().StringP("regex", "r", "*", "regex to match")
	rootCmd.PersistentFlags().IntP("count", "n", 10, "number of times to match before triggering the command")
	rootCmd.PersistentFlags().IntP("interval", "i", 5, "interval to check the log stream for matches in seconds")
	rootCmd.PersistentFlags().StringP("do", "d", "", "command to execute when a match is found <count> times in <interval> time interval")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
}
