package cmd

import (
	"context"
	"fmt"

	"github.com/alexhokl/helper/authhelper"
	"github.com/alexhokl/strava-cli/swagger"
	"github.com/spf13/cobra"
)

// showProfileCmd represents the profile command
var showProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Show profile of the current user",
	RunE:  runShowProfile,
}

func init() {
	showCmd.AddCommand(showProfileCmd)
}

func runShowProfile(_ *cobra.Command, _ []string) error {
	savedToken, err := authhelper.LoadTokenFromViper()
	if err != nil {
		return err
	}
	auth := context.WithValue(context.Background(), swagger.ContextAccessToken, savedToken.AccessToken)
	config := swagger.NewConfiguration()
	client := swagger.NewAPIClient(config)
	profile, _, err := client.AthletesApi.GetLoggedInAthlete(auth)
	if err != nil {
		return err
	}
	fmt.Printf(
		"Hi %s!\nYour FTP is %dW with weight %fkg\n",
		profile.Firstname,
		profile.Ftp,
		profile.Weight,
	)
	return nil
}
