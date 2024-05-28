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

// taskSkippedErrorsCmd represents the taskSkippedErrors command
var taskSkippedErrorsCmd = &cobra.Command{
	Use:   "taskSkippedErrors [flags] task_id",
	Short: "Retrieve discovered paths that were skipped due to \"skip_source_errors\" flag being set",
	Long: `
For completed tasks, retrieve a list of paths that were discovered
but skipped due to the "skip_source_errors" flag being set to true. 
The list will contain enough information to create new transfer_items
for retrying their transfer.`,
	Run: func(cmd *cobra.Command, args []string) {
		authCodeGrant, _ := cmd.Flags().GetBool("auth-code-grant")
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		redirectURL, _ := cmd.Flags().GetString("redirect-url")
		marker, _ := cmd.Flags().GetUint("marker")

		if len(args) != 1 {
			log.Fatal("incorrect argument count")
		}
		taskId := args[0]

		scopes := []string{
			"urn:globus:auth:scope:transfer.api.globus.org:all",
		}

		client, err := login(authCodeGrant, clientID, clientSecret, redirectURL, scopes)
		if err != nil {
			log.Fatal(err)
		}

		// get task by id
		skips, err := globus.TransferGetTaskSkippedErrors(client, taskId, marker)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("Result of request: \n")
		fmt.Printf("data_type: %s\n[\n", skips.DataType)
		for _, transfer := range skips.Data {
			fmt.Printf("\n%+v\n", transfer)
		}
		fmt.Print("]\n")
	},
}

func init() {
	rootCmd.AddCommand(taskSkippedErrorsCmd)
	taskSkippedErrorsCmd.Flags().Uint("marker", 0, "used to retreive the next page by following the 'next_marker' attribute")
}
