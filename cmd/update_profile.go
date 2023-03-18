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

	_, _, err = client.AthletesApi.UpdateLoggedInAthlete(auth, updateProfileOpts.weight)
	if err != nil {
		return err
	}

	fmt.Println("Updated profile weight")

	return nil
}
