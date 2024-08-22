package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/SwissOpenEM/globus"
	"github.com/spf13/cobra"
)

// fileListSync represents the folderSync command
var fileListSyncCmd = &cobra.Command{
	Use:   "fileListSync [flags]",
	Short: "Syncs (copies) a list of files between two Globus endpoints",
	Long: `
This command will copy all files that are in a source 
endpoint at a specified path listed in a text file relative
to a destination endpoint at its corresponding path. 
Files that already exist and have the same checksum will 
not be copied. This command does *not* support symlinks.`,
	Run: func(cmd *cobra.Command, args []string) {
		// getting auth. params
		authCodeGrant, _ := cmd.Flags().GetBool("auth-code-grant")
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		redirectURL, _ := cmd.Flags().GetString("redirect-url")

		// getting transfer params
		srcEndpoint, _ := cmd.Flags().GetString("src-endpoint")
		srcPath, _ := cmd.Flags().GetString("src-path")
		destEndpoint, _ := cmd.Flags().GetString("dest-endpoint")
		destPath, _ := cmd.Flags().GetString("dest-path")
		fileListPath, _ := cmd.Flags().GetString("file-list")

		// reading filelist
		file, err := os.Open(fileListPath)
		if err != nil {
			log.Fatalf("Error occured when opening filelist: %v\n", err)
		}
		defer file.Close()

		var files []string
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			files = append(files, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error occured when reading filelist: %v\n", err)
		}

		// note: Globus has some non-standard extensions to Oauth2, meaning that it can give out
		// multiple tokens for different endpoints with the first one being the "default".
		// To get a client that works with transfers while using a standard oauth2 library, we
		// need to exclusively specify the scopes that are associated with the transfer api.
		scopes := globus.TransferDataAccessScopeCreator([]string{srcEndpoint, destEndpoint})

		// Authenticate
		client, err := login(authCodeGrant, clientID, clientSecret, redirectURL, scopes)
		if err != nil {
			log.Fatal(err)
		}

		// Transfer - Sync folders                                                                                                        )
		result, err := client.TransferFileList(srcEndpoint, srcPath, destEndpoint, destPath, files, []bool{}, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Result of request: \n%+v\n", result)
	},
}

func init() {
	rootCmd.AddCommand(fileListSyncCmd)

	// transfer params
	fileListSyncCmd.Flags().String("src-endpoint", "", "set source endpoint")
	fileListSyncCmd.Flags().String("src-path", "", "path on source endpoint to sync")
	fileListSyncCmd.Flags().String("dest-endpoint", "", "set destination endpoint")
	fileListSyncCmd.Flags().String("dest-path", "", "path on destination endpoint to sync to")
	fileListSyncCmd.Flags().String("file-list", "", "list of files to sync (relative to src-path)")

	// mark flags as obligatory
	fileListSyncCmd.MarkFlagRequired("src-endpoint")
	fileListSyncCmd.MarkFlagRequired("src-path")
	fileListSyncCmd.MarkFlagRequired("dest-endpoint")
	fileListSyncCmd.MarkFlagRequired("dest-path")
	fileListSyncCmd.MarkFlagRequired("file-list")
}
