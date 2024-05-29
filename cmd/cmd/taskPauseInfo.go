/*
Copyright © 2024 The Swiss OpenEM Team
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// taskPauseInfoCmd represents the taskPauseInfo command
var taskPauseInfoCmd = &cobra.Command{
	Use:   "taskPauseInfo [flags] task_id",
	Short: "Get task's pause information",
	Long: `
This command receives the pause information of a task.
It will include all pause rules and pause messages.
Note: only the task owner can request this information.`,
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
			log.Fatal(err)
		}

		// get task info by id
		info, err := client.TransferGetTaskPauseInfo(taskId)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Result of request: %+v\n", info)
	},
}

func init() {
	rootCmd.AddCommand(taskPauseInfoCmd)
}
