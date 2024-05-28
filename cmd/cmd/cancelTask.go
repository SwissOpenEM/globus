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

// cancelTaskCmd represents the cancelTask command
var cancelTaskCmd = &cobra.Command{
	Use:   "cancelTask [flags] task_id",
	Short: "Cancels a Globus task by its id",
	Long: `
This command cancels a task in Globus. By specifying the id
of the task that is to be cancelled, a request will be sent
to globus to cancel it. Requires management privileges (the
user requested it, collection manager... etc.) over the task 
for this to be successful.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		authCodeGrant, _ := cmd.Flags().GetBool("auth-code-grant")
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		redirectURL, _ := cmd.Flags().GetString("redirect-url")

		if len(args) != 1 {
			log.Fatal("incorrect argument count")
		}
		taskId := args[0]

		scopes := []string{
			"urn:globus:auth:scope:transfer.api.globus.org:all",
		}

		client, err := login(authCodeGrant, clientID, clientSecret, redirectURL, scopes)
		if err != nil {
			log.Fatal()
		}

		// cancel task
		result, err := globus.TransferCancelTaskByID(client, taskId)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Result of request: %+v\n", result)
	},
}

func init() {
	rootCmd.AddCommand(cancelTaskCmd)
}
