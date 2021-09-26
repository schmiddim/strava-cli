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
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/alexhokl/helper/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const port = 9990
const apiURL = "https://www.strava.com"

type loginOptions struct {
}

var (
	conf      *oauth2.Config
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

func runLogin(cmd *cobra.Command, args []string) error {
	clientID := viper.GetString("clientId")
	clientSecret := viper.GetString("clientSecret")
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("client_id or client_secret is not configured")
	}

	ctx = context.Background()
	conf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       getScopes(),
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth/authorize", apiURL),
			TokenURL: fmt.Sprintf("%s/oauth/token", apiURL),
		},
		RedirectURL: fmt.Sprintf("http://localhost:%d/callback", port),
	}

	// add transport for self-signed certificate to context
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	fmt.Println("You will now be taken to your browser for authentication")
	time.Sleep(1 * time.Second)
	cmdName, cmdArgs := cli.GetOpenCommand(url)
	_, errOpen := exec.Command(cmdName, cmdArgs...).Output()
	if errOpen != nil {
		return errOpen
	}
	time.Sleep(1 * time.Second)

	http.HandleFunc("/callback", callbackHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

	return nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	queryParts, _ := url.ParseQuery(r.URL.RawQuery)
	code := queryParts["code"][0]

	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	viper.Set("token", token.AccessToken)
	viper.WriteConfig()

	msg := "<p><strong>Success!</strong></p>"
	msg = msg + "<p>You are authenticated and can now return to the CLI.</p>"
	fmt.Fprintf(w, msg)

	fmt.Printf("Login has been completed successfully. Tokens are stored in [%s]\n", viper.ConfigFileUsed())

	go func() {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()
}

func getScopes() []string {
	return []string{"activity:write,profile:read_all,activity:read_all"}
}
