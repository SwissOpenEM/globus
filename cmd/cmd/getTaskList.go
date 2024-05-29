/*
Copyright © 2024 The Swiss OpenEM Team
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/SwissOpenEM/globus"
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
		redirectURL, _ := cmd.Flags().GetString("redirect-url")
		limit, _ := cmd.Flags().GetUint("limit")

		if limit < 1 {
			log.Fatal(fmt.Errorf("limit can't be less than 1"))
		}

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

		client, err := login(authCodeGrant, clientID, clientSecret, redirectURL, scopes)
		if err != nil {
			log.Fatal(err)
		}

		// get task list
		transferList, err := globus.TransferGetTaskList(client, offset, limit)
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

	getTaskListCmd.Flags().Uint("offset", 0, "set the initial offset of the list for pagination (can't use with page)")
	getTaskListCmd.Flags().Uint("limit", 50, "set the max. size of the requested list")
	getTaskListCmd.Flags().Uint("page", 1, "set the page on the task list (can't use with offset)")
	getTaskListCmd.MarkFlagsMutuallyExclusive("offset", "page")
}
