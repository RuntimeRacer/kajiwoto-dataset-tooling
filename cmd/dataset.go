// Package cmd
/*
Copyright © 2021 NAME HERE runtimeracer@gmail.com

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
	"github.com/spf13/cobra"
)

// Flags
var source, target string

// datasetCmd represents the dataset command
var datasetCmd = &cobra.Command{
	Use:   "dataset",
	Short: "Interact with kajiwoto's dataset API or local dataset files",
	Long: `dataset is used to interact with kajiwoto's dataset API or local files containing dataset content. 
Only works with a valid session, otherwise it will fail.
Offers various subcommands for different kinds of dataset interaction.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(datasetCmd)

	// Flags fir dataset commands
	datasetCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "source file or URL")
	datasetCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "target file or URL")
	if err := datasetCmd.MarkPersistentFlagRequired("source"); err != nil {
		fmt.Println(err.Error())
	}
	if err := datasetCmd.MarkPersistentFlagRequired("target"); err != nil {
		fmt.Println(err.Error())
	}

}
