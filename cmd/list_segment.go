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

// listSegmentCmd represents the list segment command
var listSegmentCmd = &cobra.Command{
	Use:   "segment",
	Short: "List starred segment of the current user",
	RunE:  runListSegments,
}

func init() {
	listCmd.AddCommand(listSegmentCmd)
}

func runListSegments(_ *cobra.Command, _ []string) error {
	savedToken, err := authhelper.LoadTokenFromViper()
	if err != nil {
		return err
	}
	auth := context.WithValue(context.Background(), swagger.ContextAccessToken, savedToken.AccessToken)
	config := swagger.NewConfiguration()
	client := swagger.NewAPIClient(config)

	opts := &swagger.SegmentsApiGetLoggedInAthleteStarredSegmentsOpts{
		PerPage: optional.NewInt32(100),
		Page:    optional.NewInt32(1),
	}
	segments, _, err := client.SegmentsApi.GetLoggedInAthleteStarredSegments(auth, opts)
	if err != nil {
		return err
	}

	if listOpts.format == "json" {
		json, err := jsonhelper.GetJSONString(segments)
		if err != nil {
			return err
		}
		fmt.Println(json)
		return nil
	}

	var data [][]string
	for _, e := range segments {
		arr := []string{
			fmt.Sprintf("%d", e.Id),
			e.Country,
			e.Name,
			fmt.Sprintf("%.1f", e.Distance/1000.0),
			fmt.Sprintf("%.0f", e.ElevationHigh-e.ElevationLow),
			fmt.Sprintf("%.2f", (e.ElevationHigh-e.ElevationLow)/e.Distance*100),
		}
		data = append(data, arr)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Country", "Name", "Distance (km)", "Elevation (m)", "Average gradient (%)"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}
