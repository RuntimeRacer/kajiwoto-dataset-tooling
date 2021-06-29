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
var username, password string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Provide credentials to login to your kajiwoto account",
	Long: `login takes your kajiwoto login credentials as arguments to do the following things:

1. Login and verify correctness of credentials.
2. Storing your credentials in your kajitool config file ($HOMEDIR/.kajotool.yaml).
3. Storing your session ID in your kajitool config file for reuse when executing commands against the Kajiwoto API.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		fmt.Println("login called")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Flags for Login
	loginCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username, required for first login or if switching accounts.")
	loginCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password, required for first login or if switching accounts.")
}
