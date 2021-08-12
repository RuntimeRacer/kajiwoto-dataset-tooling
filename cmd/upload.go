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
	"errors"
	"fmt"
	"github.com/runtimeracer/kajitool/query"
	"github.com/spf13/cobra"
	"sort"
	"strings"
	"time"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Takes a training data from a specified source file and uploads it into a specified target dataset.",
	Long: `upload fetches training data from the specified source file and uploads it into the specified target dataset. 

param source: must be a local file. Data will be expected to be in csv format.
param target: a full Kajiwoto dataset URL (including ID), or a Kajiwoto dataset ID.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		// Sanity checks
		if source, err = validateUploadSource(source); err != nil {
			return err
		}
		if target, err = validateUploadTarget(target); err != nil {
			return err
		}

		// Read data from source file
		var trainingData []DatasetEntry
		if trainingData, err = readCSV(source); err != nil {
			return err
		}

		// Analyze data - Each entry / set of entries we upload creates an API request. Only send new ones
		qualified := make([]DatasetEntry, 0)
		for _, analyzed := range trainingData {
			if analyzed.ID == "" {
				qualified = append(qualified, analyzed)
			}
		}
		fmt.Println(fmt.Sprintf("Found %v new entries in source data", len(qualified)))

		// Init Client
		client := query.GetKajiwotoClient(endpoint)

		// Login via Session key
		loginResult := query.LoginResult{}
		if loginResult, err = client.DoLoginAuthToken(sessionKey); err != nil {
			return err
		}

		// Get User Info from Login result
		userInfo := &loginResult.Login.User

		// Get Info on the source Dataset
		datasetInfo := query.AITrainerGroup{}
		if datasetInfo, err = client.GetAITrainerGroup(target, sessionKey); err != nil {
			return err
		}

		// Print some info on the Dataset
		fmt.Println(fmt.Sprintf("Dataset found: %v", datasetInfo.Name))
		fmt.Println(fmt.Sprintf("Indexed entries: %v", datasetInfo.Count))

		if datasetInfo.User.ID != userInfo.ID {
			return errors.New("not your dataset! You cannot upload to foreign datasets")
		}

		// Perform Upload - FIXME: This currently can only do single / unrelated training upload
		for _, qEntry := range qualified {
			// Convert training information to a elements required by graphQL
			trainings := make([]query.AITraining, 0)

			if len(qEntry.History) > 0 {
				// Find matching history dataset entry
				bestMatch := findBestContextualMatch(qEntry, trainingData)

				// Match is Training 0
				trainings = append(trainings, bestMatch.ToAITraining(0))
				// Actual Entry is Training 1
				trainings = append(trainings, qEntry.ToAITraining(1))
			} else if len(qEntry.History) > 1 {
				return fmt.Errorf("error: Entry U: '%v' K: '%v' has too many history context items", qEntry.UserMessage, qEntry.Message)
			} else {
				// Default: No history
				trainings = append(trainings, qEntry.ToAITraining(0))
			}

			trainingResult := query.TrainDatasetResult{}
			if trainingResult, err = client.DoTrainDataset(string(datasetInfo.ID), sessionKey, trainings); err != nil {
				return err
			}
			fmt.Println(fmt.Sprintf("Training successful. New entry count: %v", trainingResult.Count))

			// Sleep 1s to not hammer the API too much.
			time.Sleep(time.Second)
		}

		return nil
	},
}

func init() {
	datasetCmd.AddCommand(uploadCmd)
}

func validateUploadSource(source string) (string, error) {
	if source == "" {
		return "", errors.New("empty source")
	}

	return source, nil
}

func validateUploadTarget(target string) (string, error) {
	if target == "" {
		return "", errors.New("empty target")
	}

	return target, nil
}

func findBestContextualMatch(contextEntry DatasetEntry, trainingData []DatasetEntry) DatasetEntry {
	// Calculate matching points based on Condition & ASM
	historyContent := contextEntry.History[0]
	matchKeys := make([]int, 0)
	matchMap := make(map[int]DatasetEntry)
	for _, match := range trainingData {
		if match.UserMessage == historyContent {
			matchPoints := 0

			// ASM
			if match.ASM == contextEntry.ASM {
				matchPoints++
			}

			// Condition
			conditionVars := strings.Split(contextEntry.Condition, "")
			matchConditionVars := strings.Split(match.Condition, "")
			for i, qCond := range conditionVars {
				if qCond == matchConditionVars[i] {
					matchPoints++
				}
			}

			// Append to map
			if _, ok := matchMap[matchPoints]; !ok {
				matchKeys = append(matchKeys, matchPoints)
			}
			matchMap[matchPoints] = match
		}
	}

	// Get best match
	sort.Ints(matchKeys)
	bestMatch, _ := matchMap[matchKeys[len(matchKeys)-1]]
	return bestMatch
}
