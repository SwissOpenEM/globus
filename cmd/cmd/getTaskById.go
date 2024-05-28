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

// getTaskByIdCmd represents the getTaskById command
var getTaskByIdCmd = &cobra.Command{
	Use:   "getTaskById [flags] task_id",
	Short: "retrieve a task's details by its id",
	Long: `
This command retrieves the task struct of a task by
its id. It can only request tasks to which the 
authenticated user has access to.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		authCodeGrant, _ := cmd.Flags().GetBool("auth-code-grant")
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		authURL, _ := cmd.Flags().GetString("redirect-url")

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

		// get task by id
		transfer, err := globus.TransferGetTaskByID(client, taskId)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Result of request: %+v\n", transfer)
	},
}

func init() {
	rootCmd.AddCommand(getTaskByIdCmd)
	getTaskByIdCmd.Flags().BoolP("auth-code-grant", "a", false, "enable authorization code based OAuth2 authentication")
}
