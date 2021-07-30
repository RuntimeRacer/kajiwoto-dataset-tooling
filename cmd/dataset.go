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
	"github.com/google/go-cmp/cmp"
	"github.com/runtimeracer/kajitool/constants"
	"github.com/runtimeracer/kajitool/query"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
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

// DatasetEntry Helper Type for Dataset entries
type DatasetEntry struct {
	ID          string
	UserMessage string
	Message     string
	/*	ASM Cheat sheet

		... No idea what "ASM" stands for in this context. But its holding the emotional values for the Kaji dialogues.

		HAPPY => Happy or Excited
		SAD   => Sad
		HUNGRY => Hungry
		FULL => Full
		EXCITED => Excited
		ANGRY => Angry
		SCARED => Scared
		BULLIED => Bullied
		ATTACKED => Attacked
		POWERED => Powered
		DRUNK => Drunk
		SICK => Sick
		SLEEPY => Sleepy
	*/
	ASM string
	/*  Condition cheat sheet

	Seems to be oriented from linux permissions. five digits; last two seem to be never used.

	//// Attachment Keys
	XX 1 XX Disliked
	XX 2 XX Any-Emotion
	XX 3 XX Liked
	XX 4 XX -- NOT USED --
	XX 5 XX Disliked/Neutral

	//// Daytime Keys
	// Default (single) conditions
	1 XXXX Early Morning AM
	2 XXXX Morning
	3 XXXX Afternoon
	4 XXXX Evening
	5 XXXX Middle of Sleep AM
	6 XXXX -- NOT USED --
	// Combined Conditions
	7 XXXX Early Morning AM - Morning
	8 XXXX Evening AM - Middle of Sleep AM
	9 XXXX Morning - Afternoon

	//// Last seen keys
	X 1 XXX Seen 2 hrs ago
	X 2 XXX Seen 12 hrs ago
	X 3 XXX Seen 5 days ago
	X 4 XXX Seen 5 days+ ago
	*/
	Condition string
	Deleted   bool
	// History contains possible preceding user dialogues
	History      []string
	DuplicateIDs []string
}

/*
	ToCSV converts a Dataset entry into a String array used for writing it to a CSV file.

	Remarks:
	- I have no idea if more than 1 Element is even supported. However, it is an array. Allowing 5 for now.

	Mapping:
	- 0:  ID
	- 1:  UserMessage
	- 2:  Message
	- 3:  ASM
	- 4:  Attachment Key => Condition
	- 5:  Daytime Key => Condition
	- 6:  Last Seen Key => Condition
	- 7:  Deleted
	- 8:  History
	- 9:  DuplicateIDs

*/
func (e *DatasetEntry) ToCSV() []string {
	result := make([]string, 10)

	result[0] = e.ID
	result[1] = e.UserMessage
	result[2] = e.Message
	result[3] = e.ASM

	// Condition Split
	attachmentVars := strings.Split(e.Condition, "")
	result[4] = attachmentVars[2]
	result[5] = attachmentVars[0]
	result[6] = attachmentVars[1]

	result[7] = strconv.FormatBool(e.Deleted)
	result[8] = strings.Join(e.History, constants.CSVListSeparator)
	result[9] = strings.Join(e.DuplicateIDs, constants.CSVListSeparator)

	return result
}

/*
	FromCSV converts a String array read from a CSV file into a Dataset entry.

	Remarks:
	- I have no idea if more than 1 Element is even supported. However, it is an array. Allowing 5 for now.

	Mapping:
	- 0:  ID
	- 1:  UserMessage
	- 2:  Message
	- 3:  ASM
	- 4:  Attachment Key => Condition
	- 5:  Daytime Key => Condition
	- 6:  Last Seen Key => Condition
	- 7:  Deleted
	- 8:  History
	- 9:  DuplicateIDs

*/
func (e *DatasetEntry) FromCSV(src []string) DatasetEntry {
	// Convert from Array elements
	condition := strings.Join([]string{src[5], src[6], src[4], "00"}, "")
	deleted, err := strconv.ParseBool(src[7])
	if err != nil {
		fmt.Println("WARNING: Error Parsing Bool from String!")
	}

	return DatasetEntry{
		ID:           src[0],
		UserMessage:  src[1],
		Message:      src[2],
		ASM:          src[3],
		Condition:    condition,
		Deleted:      deleted,
		History:      strings.Split(src[8], constants.CSVListSeparator),
		DuplicateIDs: strings.Split(src[9], constants.CSVListSeparator),
	}
}

func (e *DatasetEntry) FromAITrained(src query.AITrained) DatasetEntry {
	// Convert from History Element
	history := make([]string, len(src.History))
	for i, itm := range src.History {
		history[i] = string(itm)
	}

	return DatasetEntry{
		ID:           string(src.ID),
		UserMessage:  string(src.UserMessage),
		Message:      string(src.Message),
		ASM:          string(src.ASM),
		Condition:    string(src.Condition),
		Deleted:      bool(src.Deleted),
		History:      history,
		DuplicateIDs: make([]string, 0),
	}
}

func (e *DatasetEntry) isDuplicate(c *DatasetEntry) bool {
	if e.Message == c.Message &&
		e.UserMessage == c.UserMessage &&
		e.ASM == c.ASM &&
		e.Condition == c.Condition &&
		cmp.Equal(e.History, c.History) {
		return true
	}
	return false
}
