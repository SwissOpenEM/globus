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

// getTaskListCmd represents the getTaskList command
var getTaskListCmd = &cobra.Command{
	Use:   "getTaskList [flags]",
	Short: "Gets the current task list of a user or service account",
	Long: `
It requests the current transfer task list of the user
or service account that is provided. It will then print
out the results, with each task being printed out as
a raw struct.`,
	Run: func(cmd *cobra.Command, args []string) {
		authCodeGrant, _ := cmd.Flags().GetBool("auth-code-grant")
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		authURL, _ := cmd.Flags().GetString("auth-url")

		scopes := []string{
			"urn:globus:auth:scope:transfer.api.globus.org:all",
		}

		client, err := login(authCodeGrant, clientID, clientSecret, authURL, scopes)
		if err != nil {
			log.Fatal(err)
		}

		// get task list
		transferList, err := globus.TransferGetTaskList(client, 0, 50)
		if err != nil {
			log.Fatal(err)
		}

		// present results
		fmt.Print("Result of request: \n")
		for _, transfer := range transferList.Data {
			fmt.Printf("\n%+v\n", transfer)
		}
	},
}

func init() {
	rootCmd.AddCommand(getTaskListCmd)

	getTaskListCmd.Flags().BoolP("auth-code-grant", "a", false, "enable authorization code based OAuth2 authentication")
	getTaskListCmd.Flags().String("client-id", "", "set client ID of application")
	getTaskListCmd.Flags().String("client-secret", "", "set client secret of application")
	getTaskListCmd.Flags().String("auth-url", "", "set auth url (only used in three-legged mode)")

	getTaskListCmd.MarkFlagRequired("client-id")
	getTaskListCmd.MarkFlagRequired("client-secret")
	getTaskListCmd.MarkFlagsRequiredTogether("auth-code-grant", "auth-url")
}
