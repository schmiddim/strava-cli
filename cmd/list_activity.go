package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alexhokl/helper/authhelper"
	"github.com/alexhokl/helper/jsonhelper"
	"github.com/alexhokl/strava-cli/swagger"
	"github.com/antihax/optional"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listActivityCmd represents the list activity command
var listActivityCmd = &cobra.Command{
	Use:   "activity",
	Short: "List recent activities of the current user",
	RunE:  runListActivities,
}

func init() {
	listCmd.AddCommand(listActivityCmd)
}

func runListActivities(_ *cobra.Command, _ []string) error {
	savedToken, err := authhelper.LoadTokenFromViper()
	if err != nil {
		return err
	}
	auth := context.WithValue(context.Background(), swagger.ContextAccessToken, savedToken.AccessToken)
	config := swagger.NewConfiguration()
	client := swagger.NewAPIClient(config)

	opts := &swagger.ActivitiesApiGetLoggedInAthleteActivitiesOpts{
		PerPage: optional.NewInt32(10),
		Page:    optional.NewInt32(1),
	}
	activities, _, err := client.ActivitiesApi.GetLoggedInAthleteActivities(auth, opts)
	if err != nil {
		return err
	}

	if listOpts.format == "json" {
		json, err := jsonhelper.GetJSONString(activities)
		if err != nil {
			return err
		}
		fmt.Println(json)
		return nil
	}

	var data [][]string
	for _, e := range activities {
		arr := []string{
			fmt.Sprintf("%d", e.Id),
			e.StartDate.Local().String(),
			e.Name,
		}
		data = append(data, arr)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Date", "Activity"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}
