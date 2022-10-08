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
	"os"
	"time"

	"github.com/alexhokl/strava-cli/swagger"
	"github.com/antihax/optional"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listSegmentEffortCmd = &cobra.Command{
	Use:   "segment-effort",
	Short: "List efforts of a segment of the current user order by quickest time",
	RunE:  runListSegmentEfforts,
}

type listSegmentEffortOptions struct {
	id int32
}

var listSegmentEffortOpts listSegmentEffortOptions

func init() {
	listCmd.AddCommand(listSegmentEffortCmd)

	flags := listSegmentEffortCmd.Flags()
	flags.Int32Var(&listSegmentEffortOpts.id, "id", 0, "Segment ID")
	listSegmentEffortCmd.MarkFlagRequired("id")
}

func runListSegmentEfforts(_ *cobra.Command, _ []string) error {
	accessToken := viper.GetString("token")
	auth := context.WithValue(context.Background(), swagger.ContextAccessToken, accessToken)
	config := swagger.NewConfiguration()
	client := swagger.NewAPIClient(config)

	opts := &swagger.SegmentEffortsApiGetEffortsBySegmentIdOpts{
		StartDateLocal: optional.NewTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
		EndDateLocal: optional.NewTime(time.Now()),
		PerPage: optional.NewInt32(100),
	}
	efforts, _, err := client.SegmentEffortsApi.GetEffortsBySegmentId(auth, listSegmentEffortOpts.id, opts)
	if err != nil {
		return err
	}

	var data [][]string
	for _, e := range efforts {
		duration, _ := time.ParseDuration(fmt.Sprintf("%ds", e.ElapsedTime))
		arr := []string{
			e.StartDate.Format("2006-01-02"),
			duration.String(),
			fmt.Sprintf("%.0f", e.AverageWatts),
			fmt.Sprintf("%.0f", e.AverageCadence),
			fmt.Sprintf("%.0f", e.AverageHeartrate),
			fmt.Sprintf("%.0f", e.MaxHeartrate),
		}
		data = append(data, arr)
	}

	if len(efforts) > 0 {
		fmt.Printf("Segment: %s\n", efforts[0].Segment.Name)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Duration", "Power (W)", "Cadence", "Heart rate", "Max heart rate"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}
