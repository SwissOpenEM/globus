package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/SwissOpenEM/globus"
	"github.com/spf13/cobra"
)

// fileListSync represents the folderSync command
var getRefreshToken = &cobra.Command{
	Use:   "getRefreshToken [flags]",
	Short: "Lets the user log in and returns the refresh token associated with the user",
	Long: `
This command will execute a 3-legged OAuth authentication, with the OfflineAccess flag enabled,
which will give the application a Refresh Token. This refresh token is then returned to the command line.
Be careful with the refresh token, it provides access to the API using the identity of the user and is 
valid by default for a long period of time.`,
	Run: func(cmd *cobra.Command, args []string) {
		// getting auth. params
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		redirectURL, _ := cmd.Flags().GetString("redirect-url")

		// getting transfer params
		srcEndpoint, _ := cmd.Flags().GetString("src-endpoint")
		destEndpoint, _ := cmd.Flags().GetString("dest-endpoint")

		// note: Globus has some non-standard extensions to Oauth2, meaning that it can give out
		// multiple tokens for different endpoints with the first one being the "default".
		// To get a client that works with transfers while using a standard oauth2 library, we
		// need to exclusively specify the scopes that are associated with the transfer api.
		endpoints := []string{}
		if srcEndpoint != "" {
			endpoints = append(endpoints, srcEndpoint)
		}
		if destEndpoint != "" {
			endpoints = append(endpoints, destEndpoint)
		}
		scopes := globus.TransferDataAccessScopeCreator(endpoints)

		// Authenticate
		// create config
		ctx := context.Background()
		conf := globus.AuthGenerateOauthClientConfig(ctx, clientID, clientSecret, redirectURL, scopes)
		token, err := getToken(ctx, clientID, clientSecret, redirectURL, scopes, conf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Your refresh token is: \"%s\"\n", token.RefreshToken)
	},
}

func init() {
	rootCmd.AddCommand(getRefreshToken)

	// transfer params
	getRefreshToken.Flags().String("src-endpoint", "", "set source endpoint")
	getRefreshToken.Flags().String("dest-endpoint", "", "set destination endpoint")
}
