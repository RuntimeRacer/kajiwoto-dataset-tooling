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
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
	"github.com/runtimeracer/go-graphql-client"
	"github.com/runtimeracer/kajitool/constants"
	"github.com/runtimeracer/kajitool/query"
	"github.com/spf13/cobra"
)

const (
	csvSize     = 10
	emptyColumn = "EMPTY"
)

var (
	asmMap = map[string]string{
		"none":     "emotion_any",
		"HAPPY":    "emotion_happy_or_excited",
		"SAD":      "emotion_sad",
		"HUNGRY":   "emotion_hungry",
		"FULL":     "emotion_full",
		"EXCITED":  "emotion_ecited",
		"ANGRY":    "emotion_angry",
		"SCARED":   "emotion_scared",
		"BULLIED":  "emotion_bullied",
		"ATTACKED": "emotion_attacked",
		"POWERED":  "emotion_powered",
		"DRUNK":    "emotion_drunk",
		"SICK":     "emotion_sick",
		"SLEEPY":   "emotion_sleepy",
	}

	attachmentMap = map[string]string{
		"0": "attachment_none_NOT_USED",
		"1": "attachment_disliked",
		"2": "attachment_any",
		"3": "attachment_liked",
		"4": "attachment_UNKNOWN",
		"5": "attachment_disliked_neutral",
	}

	daytimeMap = map[string]string{
		"0": "daytime_any",
		"1": "daytime_early_morning",
		"2": "daytime_morning",
		"3": "daytime_afternoon",
		"4": "daytime_evening",
		"5": "daytime_middle_of_sleep",
		"6": "daytime_UNKNOWN",
		"7": "daytime_early_morning_till_morning",
		"8": "daytime_evening_till_middle_of_sleep",
		"9": "daytime_morning_till_afternoon",
	}

	lastSeenMap = map[string]string{
		"0": "seen_any",
		"1": "seen_2_hrs_ago",
		"2": "seen_12_hrs_ago",
		"3": "seen_5_days_ago",
		"4": "seen_5_plus_days_ago",
	}
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
		_ = cmd.Help()
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

	Seems to be inspired by linux permissions. five digits; last two seem to be never used (yet).

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
	result := make([]string, csvSize)

	result[0] = e.ID
	result[1] = e.UserMessage
	result[2] = e.Message
	result[3] = asmMap[e.ASM]

	// Condition Split
	conditionVars := strings.Split(e.Condition, "")
	result[4] = attachmentMap[conditionVars[2]] // Attachment
	result[5] = daytimeMap[conditionVars[0]]    // Daytime
	result[6] = lastSeenMap[conditionVars[1]]   // Last seen

	result[7] = strconv.FormatBool(e.Deleted)
	result[8] = strings.Join(e.History, constants.CSVListSeparator)
	if len(result[8]) == 0 {
		result[8] = emptyColumn
	}
	result[9] = strings.Join(e.DuplicateIDs, constants.CSVListSeparator)
	if len(result[9]) == 0 {
		result[9] = emptyColumn
	}

	return result
}

/*
	FromCSV converts a String array read from a CSV file into a Dataset entry.

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
func (e *DatasetEntry) FromCSV(src []string) (DatasetEntry, error) {
	// Check if valid
	if len(src) != csvSize {
		return DatasetEntry{}, fmt.Errorf("invalid length, must be %v", csvSize)
	}

	// Replace condition values with proper keys
	var foundEmotion, foundAttachment, foundDaytime, foundLastSeen bool
	for key, val := range asmMap {
		if val == src[3] {
			src[3] = key
			foundEmotion = true
			break
		}
	}
	for key, val := range attachmentMap {
		if val == src[4] {
			src[4] = key
			foundAttachment = true
			break
		}
	}
	for key, val := range daytimeMap {
		if val == src[5] {
			src[5] = key
			foundDaytime = true
			break
		}
	}
	for key, val := range lastSeenMap {
		if val == src[6] {
			src[6] = key
			foundLastSeen = true
			break
		}
	}
	if !foundEmotion {
		fmt.Println(fmt.Sprintf("WARNING: Invalid emotional key %v for dataset entry '%v'!", src[3], src[0]))
	}
	if !foundAttachment {
		fmt.Println(fmt.Sprintf("WARNING: Invalid attachment key %v for dataset entry '%v'!", src[4], src[0]))
	}
	if !foundDaytime {
		fmt.Println(fmt.Sprintf("WARNING: Invalid daytime key %v for dataset entry '%v'!", src[5], src[0]))
	}
	if !foundLastSeen {
		fmt.Println(fmt.Sprintf("WARNING: Invalid last seen key %v for dataset entry '%v'!", src[6], src[0]))
	}

	// Convert from Array elements
	condition := strings.Join([]string{src[5], src[6], src[4], "00"}, "")
	deleted, errParse := strconv.ParseBool(src[7])
	if errParse != nil {
		fmt.Println("WARNING: Error Parsing Bool from String!")
	}

	var history, duplicateIDs []string

	if len(src[8]) > 0 && src[8] != emptyColumn {
		history = strings.Split(src[8], constants.CSVListSeparator)
	}
	if len(src[9]) > 0 && src[9] != emptyColumn {
		duplicateIDs = strings.Split(src[9], constants.CSVListSeparator)
	}

	// Create dataset entry
	entry := DatasetEntry{
		ID:           src[0],
		UserMessage:  src[1],
		Message:      src[2],
		ASM:          src[3],
		Condition:    condition,
		Deleted:      deleted,
		History:      history,
		DuplicateIDs: duplicateIDs,
	}
	return entry, nil
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

func (e *DatasetEntry) ToAITraining(index int) query.AITraining {
	// ASM check -> empty string if not set
	asm := ""
	if e.ASM != "none" {
		asm = e.ASM
	}

	// Condition Split
	conditionVars := strings.Split(e.Condition, "")
	daytime := conditionVars[0]
	lastSeen := conditionVars[1]
	attachment := conditionVars[2]

	// Convert conditions to Form submit String
	conditionSubmitString := fmt.Sprintf("%v##%v%v%v0##%v##0", asm, daytime, lastSeen, attachment, index)

	return query.AITraining{
		UserMessage: graphql.String(e.UserMessage),
		Message:     graphql.String(e.Message),
		Condition:   graphql.String(conditionSubmitString),
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

// readInDatasetEntry takes a dataset entry and adds it to the store if it passes defined conditions
func readInDatasetEntry(store []DatasetEntry, entry DatasetEntry) []DatasetEntry {
	// Check for Duplicates
	for i, compare := range store {
		// Compare via Ref value; pointer safety
		ref := &compare
		if ref.isDuplicate(&entry) {
			fmt.Println(fmt.Sprintf("Warning: Dataset Entries %v and %v are identical!", compare.ID, entry.ID))
			// Mark them as duplicates for each other
			store[i].DuplicateIDs = append(compare.DuplicateIDs, entry.ID)
			entry.DuplicateIDs = append(entry.DuplicateIDs, compare.ID)
		}
	}

	// Add to the list
	store = append(store, entry)
	return store
}

// writeCSV
func writeCSV(target string, entries []DatasetEntry) error {

	csvfile, err := os.Create(target)
	if err != nil {
		return err
	}

	// write lines - Convert and store as UTF-8
	inputWriter, err := charset.NewWriter("utf-8", csvfile)
	if err != nil {
		return err
	}
	csvwriter := csv.NewWriter(inputWriter)
	for _, entry := range entries {
		if err = csvwriter.Write(entry.ToCSV()); err != nil {
			return err
		}
	}

	csvwriter.Flush()

	if err = csvfile.Close(); err != nil {
		return err
	}

	return nil
}

// readCSV
func readCSV(source string) ([]DatasetEntry, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer func() {
		if errClose := f.Close(); errClose != nil {
			fmt.Println("Warn: Unable to close file handle")
		}
	}()

	// Read in lines - Expecting Input to be UTF-8
	inputReader, err := charset.NewReader("utf-8", f)
	if err != nil {
		return nil, err
	}
	lines, err := csv.NewReader(inputReader).ReadAll()
	if err != nil {
		return nil, err
	}

	// Create new Dataset store and read in entries
	converter := &DatasetEntry{}
	entries := make([]DatasetEntry, 0)
	for i, line := range lines {
		entry, errRead := converter.FromCSV(line)
		if errRead != nil {
			fmt.Println(fmt.Sprintf("Warn: error parsing CSV line %v: %v", i, errRead.Error()))
		}
		entries = readInDatasetEntry(entries, entry)
	}
	return entries, nil
}
