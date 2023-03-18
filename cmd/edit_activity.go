package cmd

import (
	"context"
	"fmt"

	"github.com/alexhokl/helper/authhelper"
	"github.com/alexhokl/strava-cli/swagger"
	"github.com/alexhokl/strava-cli/ui"
	"github.com/antihax/optional"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type editActivityOptions struct {
	id int64
}

var editActivityOpts editActivityOptions

// editActivityCmd represents the edit activity command
var editActivityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Edit an activity",
	RunE:  runEditActivities,
}

func init() {
	editCmd.AddCommand(editActivityCmd)

	flags := editActivityCmd.Flags()
	flags.Int64Var(&editActivityOpts.id, "id", 0, "Activity ID")
	editActivityCmd.MarkFlagRequired("id")
}

func runEditActivities(_ *cobra.Command, _ []string) error {
	savedToken, err := authhelper.LoadTokenFromViper()
	if err != nil {
		return err
	}
	auth := context.WithValue(context.Background(), swagger.ContextAccessToken, savedToken.AccessToken)
	config := swagger.NewConfiguration()
	client := swagger.NewAPIClient(config)

	getOpts := &swagger.ActivitiesApiGetActivityByIdOpts{
		IncludeAllEfforts: optional.NewBool(false),
	}
	activity, _, err := client.ActivitiesApi.GetActivityById(auth, editActivityOpts.id, getOpts)
	if err != nil {
		return err
	}

	p := tea.NewProgram(ui.NewEditorModel(activity.Name, activity.Description))
	updatedModel, err := p.Run()
	if err != nil {
		return err
    }
	updatedEditorModel := updatedModel.(ui.EditorModel)
	if !updatedEditorModel.HasUpdate() {
		fmt.Println("No update made")
		return nil
	}

	activity.Name = updatedEditorModel.Name()
	activity.Description = updatedEditorModel.Description()

	opts := &swagger.ActivitiesApiUpdateActivityByIdOpts{
		Body: optional.NewInterface(swagger.UpdatableActivity{
			Commute:      activity.Commute,
			Trainer:      activity.Trainer,
			HideFromHome: activity.HideFromHome,
			Description:  activity.Description,
			Name:         activity.Name,
			Type_:        activity.Type_,
			SportType:    activity.SportType,
			GearId:       activity.GearId,
		}),
	}

	_, _, err = client.ActivitiesApi.UpdateActivityById(auth, editActivityOpts.id, opts)
	if err != nil {
		return err
	}

	fmt.Println("Updated completed.")

	return nil
}


