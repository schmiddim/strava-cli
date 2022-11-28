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

	"github.com/alexhokl/helper/authhelper"
	"github.com/alexhokl/strava-cli/swagger"
	"github.com/spf13/cobra"
)

// updateProfileCmd represents the update weight command
var updateProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Update profile of the current user",
	RunE:  runUpdateProfile,
}

type updateProfileOptions struct {
	weight float32
}

var updateProfileOpts updateProfileOptions


func init() {
	updateCmd.AddCommand(updateProfileCmd)

	flags := updateProfileCmd.Flags()
	flags.Float32VarP(&updateProfileOpts.weight, "weight", "w", 0, "Weight in kg")
	updateProfileCmd.MarkFlagRequired("weight")
}

func runUpdateProfile(_ *cobra.Command, _ []string) error {
	savedToken, err := authhelper.LoadTokenFromViper()
	if err != nil {
		return err
	}
	auth := context.WithValue(context.Background(), swagger.ContextAccessToken, savedToken.AccessToken)
	config := swagger.NewConfiguration()
	client := swagger.NewAPIClient(config)

	fmt.Printf("About to update to %f kg", updateProfileOpts.weight)
	athlete, _, err := client.AthletesApi.UpdateLoggedInAthlete(auth, updateProfileOpts.weight)
	if err != nil {
		return err
	}

	fmt.Printf("Weight updated to %f kg", athlete.Weight)

	return nil
}
