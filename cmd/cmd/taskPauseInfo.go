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
		authURL, _ := cmd.Flags().GetString("auth-url")

		if len(args) != 1 {
			log.Fatal("incorrect argument count")
		}
		taskId := args[0]

		scopes := []string{
			"urn:globus:auth:scope:transfer.api.globus.org:all",
		}

		client, err := login(authCodeGrant, clientID, clientSecret, authURL, scopes)
		if err != nil {
			log.Fatal(err)
		}

		// get task info by id
		info, err := globus.TransferGetTaskPauseInfo(client, taskId)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Result of request: %+v\n", info)
	},
}

func init() {
	rootCmd.AddCommand(taskPauseInfoCmd)

	taskPauseInfoCmd.Flags().BoolP("auth-code-grant", "a", false, "enable authorization code based OAuth2 authentication")
	taskPauseInfoCmd.Flags().String("client-id", "", "set client ID of application")
	taskPauseInfoCmd.Flags().String("client-secret", "", "set client secret of application")
	taskPauseInfoCmd.Flags().String("auth-url", "", "set auth url (only used in three-legged mode)")

	taskPauseInfoCmd.MarkFlagRequired("client-id")
	taskPauseInfoCmd.MarkFlagRequired("client-secret")
	taskPauseInfoCmd.MarkFlagsRequiredTogether("auth-code-grant", "auth-url")
}
