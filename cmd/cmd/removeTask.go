/*
Copyright © 2024 The Swiss OpenEM Team
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// removeTaskCmd represents the removeTask command
var removeTaskCmd = &cobra.Command{
	Use:   "removeTask [flags] task_id",
	Short: "Removes a Globus task by its id",
	Long: `
It removes a task in Globus. It is *not* equivalent to 
CANCELING the task. According to Globus "Only tasks 
against High Assurance collection are eligible for
removal." The task's state must be either SUCCEEDED or
FAILED.`,
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

		// remove task
		result, err := client.TransferRemoveTaskByID(taskId)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Result of request: %+v\n", result)
	},
}

func init() {
	rootCmd.AddCommand(removeTaskCmd)
}
