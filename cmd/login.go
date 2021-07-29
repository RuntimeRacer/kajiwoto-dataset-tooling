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
	"github.com/spf13/viper"
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
3. Storing your session ID in your kajitool config file for reuse when executing further commands against the Kajiwoto API.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Init Client
		client := query.GetKajiwotoClient(endpoint)

		// Check whether there is a Session key defined
		loginResult := query.LoginResult{}
		var errLogin error
		if sessionKey != "" {
			fmt.Println(fmt.Sprintf("Performing login via Session key: %v", sessionKey))
			loginResult, errLogin = client.DoLoginAuthToken(sessionKey)
			if errLogin != nil {
				fmt.Println(fmt.Sprintf("Unable to login via auth token, trying with username / password. error: %v", errLogin))
				loginResult, errLogin = client.DoLoginUserPW(username, password)
			}
		} else {
			fmt.Println("Performing login via Username / Password combo")
			loginResult, errLogin = client.DoLoginUserPW(username, password)
		}

		// Check for error
		if errLogin != nil {
			return fmt.Errorf("unable to login, response: %q", errLogin)
		}

		// Validate response
		if loginResult.Login.AuthToken == "" {
			return errors.New("invalid response from server: Auth token empty")
		}

		// Seems like Login worked
		userInfo := &loginResult.Login.User
		fmt.Println(fmt.Sprintf("Login successful! Hello %v!", userInfo.DisplayName))

		// Update Auth token in config file
		sessionKey = loginResult.Login.AuthToken
		viper.Set("sessionkey", loginResult.Login.AuthToken)
		fmt.Println(fmt.Sprintf("Session key: %v", sessionKey))
		if err := viper.WriteConfig(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Flags for Login
	loginCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username, required for first login or if switching accounts.")
	loginCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password, required for first login or if switching accounts.")

}
