/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	"github.com/alexhokl/strava-cli/swagger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Show profile of the current user",
	RunE:  runProfile,
}

func init() {
	rootCmd.AddCommand(profileCmd)
}

func runProfile(cmd *cobra.Command, args []string) error {
	accessToken := viper.GetString("token")
	auth := context.WithValue(context.Background(), swagger.ContextAccessToken, accessToken)
	config := swagger.NewConfiguration()
	client := swagger.NewAPIClient(config)
	profile, _, err := client.AthletesApi.GetLoggedInAthlete(auth)
	if err != nil {
		return err
	}
	fmt.Printf(
		"Hi %s!\nYour user ID is %d\nYour FTP is %dW with weight %fkg\n",
		profile.Firstname,
		profile.Id,
		profile.Ftp,
		profile.Weight,
	)
	return nil
}
