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
	"context"
	"fmt"
	"github.com/runtimeracer/go-graphql-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GraphQL for requests
type kajiwotoLoginMutation struct {
	Login struct {
		AuthToken graphql.String
		User      struct {
			ID          graphql.String
			Activated   graphql.Boolean
			Moderator   graphql.Boolean
			Username    graphql.String
			DisplayName graphql.String
			Plus        struct {
				ExpireAt  uint64
				Cancelled graphql.Boolean
				Icon      graphql.Int
				Coins     graphql.Int
				Type      graphql.String
			}
			Creator struct {
				AllowSubscriptions graphql.Boolean
				DatasetTags        []graphql.String
			}
			Profile struct {
				FirstName   graphql.String
				LastName    graphql.String
				Description graphql.String
				Gender      graphql.String
				Birthday    graphql.String
				PhotoUri    graphql.String
			}
			Email struct {
				Address  graphql.String
				Verified graphql.Boolean
			}
		}
		Settings struct {
			PersonalRoomOrder []graphql.String
			FavoriteRoomIds   []graphql.String
			FavoriteEmojis    []graphql.String
		}
	} `graphql:"login (usernameOrEmail: $usernameOrEmail, password: $password, deviceType: WEB)"`
	Welcome struct {
		WebVersion   graphql.String
		Announcement struct {
			Date      uint64
			Title     graphql.String
			Emojis    graphql.String
			Content   []graphql.String
			TextColor graphql.String
		}
	}
}

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
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check for flags properly set
		if username == "" || password == "" {
			return fmt.Errorf("invalid login credentials")
		}

		// Get GraphQL Client and execute login
		client := getGraphEndpointClient()

		vars := map[string]interface{}{
			"usernameOrEmail": graphql.String(username),
			"password":        graphql.String(password),
		}

		loginResult := kajiwotoLoginMutation{}
		if errLogin := client.Mutate(context.Background(), &loginResult, vars); errLogin != nil {
			fmt.Println(fmt.Sprintf("Unable to login, response: %v", errLogin))
			return nil
		}

		if loginResult.Login.AuthToken == "" {
			fmt.Println("Invalid response from server: Auth token empty.")
			return nil
		}

		// Seems like Login worked
		userInfo := &loginResult.Login.User
		fmt.Println(fmt.Sprintf("Login successful! Hello %v!", userInfo.DisplayName))

		// Update Auth token in config file
		// TODO

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Flags for Login
	loginCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username, required for first login or if switching accounts.")
	loginCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password, required for first login or if switching accounts.")

	// Lookup from Config
	viper.BindPFlag("username", loginCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", loginCmd.PersistentFlags().Lookup("password"))

}
