package cmd

import (
	"fmt"
	"os"

	"github.com/aminehmida/flogd/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "A brief description of your command",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug().Msg("Save invoked")
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("desc")
		stype, _ := cmd.Flags().GetString("type")
		regex, _ := cmd.Flags().GetString("regex")
		do, _ := cmd.Flags().GetString("do")
		if do == "" {
			log.Warn().Msg("\"do\" argument not defined. Will not execute any command on match")
		}
		count, _ := cmd.Flags().GetInt("count")
		interval, _ := cmd.Flags().GetInt("interval")
		configFileName := cmd.Flag("config").Value.String()
		if configFileName == "" {
			log.Warn().Msg("Config file not defined. Using default config file: ./flogd.yaml")
			configFileName = "./flogd.yaml"
		}
		configData := config.Config{
			Name:        name,
			Description: desc,
			SType:       stype,
			Regex:       regex,
			Do:          do,
			Count:       count,
			Interval:    interval,
			Command:     args[0],
		}
		// Check if config file exists
		if _, err := os.Stat(configFileName); os.IsNotExist(err) {
			configsData := config.Configs{configData}
			configBytes, err := configsData.Encode()
			if err != nil {
				log.Error().Msg(fmt.Sprintf("Can not encode config: %v", err))
				return
			}
			os.WriteFile(configFileName, configBytes, 0644)
			return
		}
		// Read config file
		configBytes, err := os.ReadFile(configFileName)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Can not read config file: %v", err))
			return
		}
		var existingConfigData config.Configs
		err = existingConfigData.Decode(configBytes)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Can not decode config file: %v", err))
			return
		}
		// Append new config
		existingConfigData = append(existingConfigData, configData)
		// Encode config
		configBytes, err = existingConfigData.Encode()
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Can not encode config: %v", err))
			return
		}
		// Write config file
		os.WriteFile(configFileName, configBytes, 0644)
		log.Info().Msg(fmt.Sprintf("Config: %s, saved to file: %s", name, configFileName))

	},
}

func init() {
	rootCmd.AddCommand(saveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	saveCmd.Flags().StringP("name", "m", "", "Config name")
	saveCmd.Flags().StringP("desc", "s", "", "Config description")
}
