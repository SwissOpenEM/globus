/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/SwissOpenEM/globus-transfer-request"
	"github.com/spf13/cobra"
)

// folderSyncCmd represents the folderSync command
var folderSyncCmd = &cobra.Command{
	Use:   "folderSync [flags]",
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
		authURL, _ := cmd.Flags().GetString("redirect-url")

		// getting transfer params
		srcEndpoint, _ := cmd.Flags().GetString("src-endpoint")
		srcPath, _ := cmd.Flags().GetString("src-path")
		destEndpoint, _ := cmd.Flags().GetString("dest-endpoint")
		destPath, _ := cmd.Flags().GetString("dest-path")

		// note: Globus has some non-standard extensions to Oauth2, meaning that it can give out
		// multiple tokens for different endpoints with the first one being the "default".
		// To get a client that works with transfers while using a standard oauth2 library, we
		// need to exclusively specify the scopes that are associated with the transfer api.
		scopes := globus.TransferDataAccessScopeCreator([]string{srcEndpoint, destEndpoint})

		// Authenticate
		client, err := login(authCodeGrant, clientID, clientSecret, authURL, scopes)
		if err != nil {
			log.Fatal(err)
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

	// transfer params
	folderSyncCmd.Flags().String("src-endpoint", "", "set source endpoint")
	folderSyncCmd.Flags().String("src-path", "", "path on source endpoint to sync")
	folderSyncCmd.Flags().String("dest-endpoint", "", "set destination endpoint")
	folderSyncCmd.Flags().String("dest-path", "", "path on destination endpoint to sync to")

	// mark flags as obligatory
	folderSyncCmd.MarkFlagRequired("src-endpoint")
	folderSyncCmd.MarkFlagRequired("src-path")
	folderSyncCmd.MarkFlagRequired("dest-endpoint")
	folderSyncCmd.MarkFlagRequired("dest-path")

}
