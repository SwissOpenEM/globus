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

// listTaskEventsCmd represents the listTaskEvents command
var listTaskEventsCmd = &cobra.Command{
	Use:   "listTaskEvents [flags] task_id",
	Short: "Retrieve the list of events of a task",
	Long: `
It retrieves the list of events associated with a task,
the latter of which is specified through its id. The
event list is only kept up to 30 days after the completion
of the task, according to Globus docs.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		authCodeGrant, _ := cmd.Flags().GetBool("auth-code-grant")
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		authURL, _ := cmd.Flags().GetString("auth-url")
		limit, _ := cmd.Flags().GetUint("limit")

		if limit < 1 {
			log.Fatal(fmt.Errorf("limit can't be less than 1"))
		}

		if len(args) != 1 {
			log.Fatal("incorrect argument count")
		}
		taskId := args[0]

		// get offset - either by page number or directly specified offset
		var offset uint
		if cmd.Flags().Lookup("page").Changed {
			if cmd.Flags().Lookup("offset").Changed {
				log.Fatal(fmt.Errorf("both page and offset are specified at the same time"))
			}
			offset, _ = cmd.Flags().GetUint("page")
			if offset == 0 {
				offset += 1
			}
			offset = (offset - 1) * limit
		} else {
			offset, _ = cmd.Flags().GetUint("offset")
		}

		// login
		scopes := []string{
			"urn:globus:auth:scope:transfer.api.globus.org:all",
		}

		client, err := login(authCodeGrant, clientID, clientSecret, authURL, scopes)
		if err != nil {
			log.Fatal(err)
		}

		// get event list of task
		eventList, err := globus.TransferGetTaskEventList(client, taskId, offset, limit)
		if err != nil {
			log.Fatal(err)
		}

		// present results
		fmt.Print("Result of request: \n")
		for _, event := range eventList.Data {
			fmt.Printf("\n%+v\n", event)
		}
	},
}

func init() {
	rootCmd.AddCommand(listTaskEventsCmd)

	listTaskEventsCmd.Flags().BoolP("auth-code-grant", "a", false, "enable authorization code based OAuth2 authentication")
	listTaskEventsCmd.Flags().String("client-id", "", "set client ID of application")
	listTaskEventsCmd.Flags().String("client-secret", "", "set client secret of application")
	listTaskEventsCmd.Flags().String("auth-url", "", "set auth url (only used in three-legged mode)")

	listTaskEventsCmd.Flags().Uint("offset", 0, "set the initial offset of the list for pagination (can't use with page)")
	listTaskEventsCmd.Flags().Uint("limit", 50, "set the max. size of the requested list")
	listTaskEventsCmd.Flags().Uint("page", 1, "set the page on the task list (can't use with offset)")

	listTaskEventsCmd.MarkFlagRequired("client-id")
	listTaskEventsCmd.MarkFlagRequired("client-secret")
	listTaskEventsCmd.MarkFlagsRequiredTogether("auth-code-grant", "auth-url")
	listTaskEventsCmd.MarkFlagsMutuallyExclusive("offset", "page")
}
