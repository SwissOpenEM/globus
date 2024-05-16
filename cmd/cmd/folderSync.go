/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/SwissOpenEM/globus-transfer-request"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// folderSyncCmd represents the folderSync command
var folderSyncCmd = &cobra.Command{
	Use:   "folderSync",
	Short: "Syncs (copies) a folder between two Globus endpoints",
	Long: `
This command will copy all files that are in a source 
endpoint at a specified path to a destination endpoint
at its corresponding path. Files that already exist and
have the same checksum will not be copied.`,
	Run: func(cmd *cobra.Command, args []string) {
		// getting auth. params
		authCodeGrant, _ := cmd.Flags().GetBool("auth-code-grant")
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		authURL, _ := cmd.Flags().GetString("auth-url")

		// getting transfer params
		srcEndpoint, _ := cmd.Flags().GetString("src-endpoint")
		srcPath, _ := cmd.Flags().GetString("src-path")
		destEndpoint, _ := cmd.Flags().GetString("dest-endpoint")
		destPath, _ := cmd.Flags().GetString("dest-path")

		flag.Parse()

		// note: Globus has some non-standard extensions to Oauth2, meaning that it can give out
		// multiple tokens for different endpoints with the first one being the "default".
		// To get a client that works with transfers while using a standard oauth2 library, we
		// need to exclusively specify the scopes that are associated with the transfer api.
		scopes := globus.TransferDataAccessScopeCreator([]string{srcEndpoint, destEndpoint})

		// Authenticate
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

		// Transfer - Sync folders
		result, err := globus.TransferFolderSync(client, srcEndpoint, srcPath, destEndpoint, destPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Result of request: \n%+v\n", result)
	},
}

func init() {
	rootCmd.AddCommand(folderSyncCmd)

	// auth. params
	folderSyncCmd.Flags().BoolP("auth-code-grant", "a", false, "enable auth code grant mode (3-legged auth.)")
	folderSyncCmd.Flags().String("client-id", "", "set client ID of application")
	folderSyncCmd.Flags().String("client-secret", "", "set client secret of application")
	folderSyncCmd.Flags().String("auth-url", "", "set auth url (only used in three-legged mode)")

	// transfer params
	folderSyncCmd.Flags().String("src-endpoint", "", "set source endpoint")
	folderSyncCmd.Flags().String("src-path", "", "path on source endpoint to sync")
	folderSyncCmd.Flags().String("dest-endpoint", "", "set destination endpoint")
	folderSyncCmd.Flags().String("dest-path", "", "path on destination endpoint to sync to")

	// mark flags as obligatory
	folderSyncCmd.MarkFlagRequired("client-id")
	folderSyncCmd.MarkFlagRequired("client-secret")
	folderSyncCmd.MarkFlagRequired("src-endpoint")
	folderSyncCmd.MarkFlagRequired("src-path")
	folderSyncCmd.MarkFlagRequired("dest-endpoint")
	folderSyncCmd.MarkFlagRequired("dest-path")
	folderSyncCmd.MarkFlagsRequiredTogether("auth-code-grant", "auth-url")
}
