package globus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c GlobusClient) getSubmissionId() (submissionId string, err error) {
	if c.client == nil {
		return "", fmt.Errorf("client is nil")
	}

	resp, err := c.client.Get(transferBaseUrl + "/submission_id")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 || resp.Status != "200 OK" {
		return "", fmt.Errorf("unexpected status for submission id request: %d '%s' - %s", resp.StatusCode, resp.Status, string(body))
	}

	var result SubmissionId
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("could not parse response for submission id request")
	}
	if result.DataType != "submission_id" {
		return "", fmt.Errorf("incorrect value type returned for submission id request: %s", result.DataType)
	}

	return result.Value, nil
}

// Submits a generic transfer request using a Transfer struct.
// This function doesn't check whether the transfer struct is valid.
// You don't need to set the submission id of the transfer, this function does that for you.
func (c GlobusClient) TransferPostTask(transfer Transfer) (result TransferResult, err error) {
	// get submission id for submission
	submission_id, err := c.getSubmissionId()
	if err != nil {
		return TransferResult{}, err
	}

	// formulate request
	transfer.CommonTransfer.SubmissionId = submission_id

	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return TransferResult{}, err
	}

	// send request
	resp, err := c.client.Post(
		transferBaseUrl+"/transfer",
		"application/json",
		bytes.NewBuffer(transferJSON),
	)
	if err != nil {
		return TransferResult{}, err
	}
	defer resp.Body.Close()

	// read & return response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TransferResult{}, err
	}

	fmt.Printf("Transfer req - status: %s\n", resp.Status)
	if resp.StatusCode == 403 {
		var consent ConsentRequired
		err = json.Unmarshal(body, &consent)
		if err != nil {
			return TransferResult{}, fmt.Errorf("unknown 403 forbidden error - status: %s, body: %s", resp.Status, body)
		}
		return TransferResult{}, fmt.Errorf("consent is required: %+v", consent)
	}

	err = json.Unmarshal(body, &result)
	return result, err
}

func (c GlobusClient) TransferCopyFile(client *http.Client, sourceEndpoint string, sourceFile string, destEndpoint string, destFile string) (TransferResult, error) {
	// formulate request
	transfer := Transfer{
		CommonTransfer: CommonTransfer{
			DataType:     "transfer",
			SubmissionId: "",
		},
		SourceEndpoint:      sourceEndpoint,
		DestinationEndpoint: destEndpoint,
		Data: []TransferItem{
			{
				DataType:        "transfer_item",
				SourcePath:      sourceFile,
				DestinationPath: destFile,
			},
		},
	}

	// submit request
	return c.TransferPostTask(transfer)
}

// submits a transfer task to copy a folder recursively.
// NOTE: the transfer follows all default params (aside from recursivity)
func (c GlobusClient) TransferFolderSync(sourceEndpoint string, sourcePath string, destEndpoint string, destPath string) (TransferResult, error) {
	// formulate request
	transfer := Transfer{
		CommonTransfer: CommonTransfer{
			DataType:     "transfer",
			SubmissionId: "",
		},
		SourceEndpoint:      sourceEndpoint,
		DestinationEndpoint: destEndpoint,
		Data: []TransferItem{
			{
				DataType:        "transfer_item",
				SourcePath:      sourcePath,
				DestinationPath: destPath,
				Recursive:       boolPointer(true),
			},
		},
	}

	// submit request
	return c.TransferPostTask(transfer)
}

func (c GlobusClient) TransferFileList(sourceEndpoint string, sourcePath string, destEndpoint string, destPath string, fileList []string, isSymlink []bool) (TransferResult, error) {
	if len(isSymlink) > 0 && len(fileList) != len(isSymlink) {
		return TransferResult{}, errors.New("isSymlink list is defined and is not the same length as fileList")
	}
	var tItems []TransferItem
	for i, file := range fileList {
		itemType := "transfer_item"
		if len(isSymlink) > 0 && isSymlink[i] {
			itemType = "transfer_symlink_item"
		}
		tItems = append(tItems, TransferItem{
			DataType:        itemType,
			SourcePath:      sourcePath + file,
			DestinationPath: destPath + file,
		})
	}

	transfer := Transfer{
		CommonTransfer: CommonTransfer{
			DataType:     "transfer",
			SubmissionId: "",
		},
		SourceEndpoint:      sourceEndpoint,
		DestinationEndpoint: destEndpoint,
		Data:                tItems,
	}

	return c.TransferPostTask(transfer)
}
