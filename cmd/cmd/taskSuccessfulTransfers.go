/*
Copyright © 2024 The Swiss OpenEM Team
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/SwissOpenEM/globus-transfer-request"
	"github.com/spf13/cobra"
)

// taskSuccessfulTransfersCmd represents the taskSuccessfulTransfers command
var taskSuccessfulTransfersCmd = &cobra.Command{
	Use:   "taskSuccessfulTransfers [flags] task_id",
	Short: "Retrieve a list of successful transfers related to a completed task",
	Long: `
Using a task id, it retrieves the associated task's list of successfully 
transfered files from the Globus API. It can only be used with completed tasks.`,
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
		transfers, err := globus.TransferGetTaskSuccessfulTransfers(client, taskId, marker)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("Result of request: \n")
		fmt.Printf("data_type: %s\n[\n", transfers.DataType)
		for _, transfer := range transfers.Data {
			fmt.Printf("\n%+v\n", transfer)
		}
		fmt.Print("]\n")
	},
}

func init() {
	rootCmd.AddCommand(taskSuccessfulTransfersCmd)
	taskSuccessfulTransfersCmd.Flags().Uint("marker", 0, "used to retreive the next page by following the 'next_marker' attribute")
}
