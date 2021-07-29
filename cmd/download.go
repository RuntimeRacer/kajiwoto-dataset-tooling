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
	"errors"
	"fmt"
	"github.com/runtimeracer/kajitool/constants"
	"github.com/runtimeracer/kajitool/query"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads a dataset from a specified source and stores it in a specified target file.",
	Long: `download fetches dataset content from the specified source and saves it into the specified target. 

param source: a full Kajiwoto dataset URL (including ID), or a Kajiwoto dataset ID.
param target: must be a local file. Data will be saved in csv format.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		// Sanity checks
		if source, err = validateSource(source); err != nil {
			return err
		}
		if target, err = validateTarget(target); err != nil {
			return err
		}

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
		if datasetInfo, err = client.GetAITrainerGroup(source, sessionKey); err != nil {
			return err
		}

		// Print some info on the Dataset
		fmt.Println(fmt.Sprintf("Dataset found: %v", datasetInfo.Name))
		fmt.Println(fmt.Sprintf("Indexed entries: %v", datasetInfo.Count))

		/*
			Safety mechanism: Only allow download of own datasets
			This is due to the some creators on kajiwoto selling complex datasets for coins and earning money from it.
			Being able to download the contents of a dataset and re-uploading it to an own dataset via kajitool, would
			make it very easy to bypass this. Of course the code switch is easily removed, but I want you to be aware
			of what you're doing here.

			Please don't be cheap. If you like a dataset, please respect the work put into it by the creator, and BUY IT!
		*/
		if datasetInfo.User.ID != userInfo.ID && datasetInfo.Price > 0 && !datasetInfo.Purchased {
			return errors.New("not your dataset! Please buy it to be able to download")
		}

		datasetContent := make([]query.AITrained, 0)

		// Fetch Dataset into result list
		// Continue as long as the result set size equals fetch limit, which means there must be another page
		var page = 0
		var datasetQueryResult []query.AITrained
		for limit := constants.FetchLimit; limit >= constants.FetchLimit; page++ {
			// Read subset of dataset
			datasetQueryResult, err = client.GetAITrainedList(string(datasetInfo.ID), "", sessionKey, limit, page)
			if err != nil {
				return err
			}

			datasetContent = append(datasetContent, datasetQueryResult...)
			limit = len(datasetQueryResult)

			if limit >= constants.FetchLimit {
				// Sleep 2 secs to not bombard the API
				time.Sleep(time.Second * 2)
			}
		}

		// Inform user on amount of fetch
		fmt.Println(fmt.Sprintf("Done. Fetched %v dataset entries.", len(datasetContent)))

		// Write to target file
		if err = writeCSV(target, datasetContent); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	datasetCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func validateSource(source string) (string, error) {
	if source == "" {
		return "", errors.New("empty source")
	}

	return source, nil
}

func validateTarget(target string) (string, error) {
	if target == "" {
		return "", errors.New("empty target")
	}

	return target, nil
}

func writeCSV(target string, entries []query.AITrained) error {

	csvfile, err := os.Create(target)
	if err != nil {
		return err
	}

	csvwriter := csv.NewWriter(csvfile)
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
