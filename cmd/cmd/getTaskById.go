/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/SwissOpenEM/globus-transfer-request"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// getTaskByIdCmd represents the getTaskById command
var getTaskByIdCmd = &cobra.Command{
	Use:   "getTaskById [flags] task_id",
	Short: "retrieve a task's details by its id",
	Long: `
	This command retrieves the task struct of a task by
	its id. It can only request tasks to which the 
	authenticated user has access to.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		authCodeGrant, _ := cmd.Flags().GetBool("auth-code-grant")
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		authURL, _ := cmd.Flags().GetString("auth-url")

		if len(args) != 1 {
			log.Fatal("incorrect argument count")
		}
		taskId := args[0]

		scopes := []string{
			"urn:globus:auth:scope:transfer.api.globus.org:all",
		}

		ctx := context.Background()
		var client *http.Client = nil
		if authCodeGrant {
			// 3-legged OAuth2 authentication - client software authenticates as a user

			// on https://app.globus.org/settings/developers/, you must select "Register a thick client or
			// script that will be installed and run by users on their devices" as registration type

			// create config
			conf := globus.AuthGenerateOauthClientConfig(ctx, clientID, clientSecret, authURL, scopes)

			// PKCE verifier
			verifier := oauth2.GenerateVerifier()

			// redirect user to consent page to ask for permission and obtain the code
			url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
			fmt.Printf("Visit the URL for the auth dialog: %v", url)

			// read-in and exchange code for token
			var code string
			if _, err := fmt.Scan(&code); err != nil {
				log.Fatal(err)
			}
			tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
			if err != nil {
				log.Fatal(err)
			}

			// create client
			client = conf.Client(ctx, tok)
		} else {
			// 2-legged OAuth2 authentication - client software authenticates as itself

			// on https://app.globus.org/settings/developers/, you must select "Register a service account or
			// application for automation" as registration type

			// create client
			var err error
			client, err = globus.AuthCreateServiceClient(ctx, clientID, clientSecret, scopes)
			if err != nil {
				log.Fatal(err)
			}
			if client == nil {
				log.Fatal("AUTH error: Client is nil\n")
			}
		}
		// get task by id
		transfer, err := globus.TransferGetTaskByID(client, taskId)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Result of request: %+v\n", transfer)
	},
}

func init() {
	rootCmd.AddCommand(getTaskByIdCmd)

	getTaskByIdCmd.Flags().BoolP("auth-code-grant", "a", false, "enable authorization code based OAuth2 authentication")
	getTaskByIdCmd.Flags().String("client-id", "", "set client ID of application")
	getTaskByIdCmd.Flags().String("client-secret", "", "set client secret of application")
	getTaskByIdCmd.Flags().String("auth-url", "", "set auth url (only used in three-legged mode)")

	getTaskByIdCmd.MarkFlagRequired("client-id")
	getTaskByIdCmd.MarkFlagRequired("client-secret")
	getTaskByIdCmd.MarkFlagsRequiredTogether("auth-code-grant", "auth-url")
}
