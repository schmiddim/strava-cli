package cmd

import (
	"context"
	"fmt"

	"github.com/alexhokl/helper/authhelper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const port = 9990
const apiURL = "https://www.strava.com"

type loginOptions struct {
}

var (
	ctx       context.Context
	loginOpts loginOptions
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log onto Strava API",
	RunE:  runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(_ *cobra.Command, _ []string) error {
	config, err := getOAuthConfigurationFromViper()
	if err != nil {
		return err
	}

	ctx = context.Background()
	token, err := authhelper.GetToken(ctx, config, false)
	if err != nil {
		return err
	}

	authhelper.SaveTokenToViper(token)
	viper.WriteConfig()
	fmt.Printf("Login has been completed successfully. Tokens are stored in [%s]\n", viper.ConfigFileUsed())

	return nil
}

func getScopes() []string {
	return []string{"activity:write,profile:read_all,activity:read_all,profile:write"}
}
