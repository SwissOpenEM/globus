package globus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const transferBaseUrl = "https://transfer.api.globusonline.org/v0.10"

// helper funcs.

func boolPointer(v bool) *bool { return &v }

//func stringPointer(v string) *string { return &v }

// request structures

type TransferItem struct {
	DataType        string `json:"DATA_TYPE"` // = "tranfer_item" OR "transfer_symlink_item"
	SourcePath      string `json:"source_path"`
	DestinationPath string `json:"destination_path"`
	// optionals
	Recursive         *bool   `json:"recursive,omitempty"`
	ExternalChecksum  *string `json:"external_checksum,omitempty"`
	ChecksumAlgorithm *string `json:"checksum_algorithm,omitempty"`
}

type DeleteItem struct {
	DataType string `json:"DATA_TYPE"` // always delete_item
	Path     string `json:"path"`
}

type FilterRule struct {
	DataType string `json:"DATA_TYPE"` // = filter_rule
	Method   string `json:"method"`
	Type     string `json:"type"`
	Name     string `json:"name"`
}

type CommonTransfer struct {
	DataType     string `json:"DATA_TYPE"` // = transfer OR delete
	SubmissionId string `json:"submission_id"`
	// optional fields
	Label               *string `json:"label,omitempty"`
	NotifyOnSucceeded   *bool   `json:"notify_on_succeeded,omitempty"`
	NotifyOnFailed      *bool   `json:"notify_on_failed,omitempty"`
	NotifyOnInactive    *bool   `json:"notify_on_inactive,omitempty"`
	SkipActivationCheck *bool   `json:"skip_activation_check,omitempty"`
	Deadline            *string `json:"deadline,omitempty"`
	StoreBasePathInfo   *bool   `json:"store_base_path_info,omitempty"`
}

type Transfer struct {
	CommonTransfer
	SourceEndpoint      string         `json:"source_endpoint"`
	DestinationEndpoint string         `json:"destination_endpoint"`
	Data                []TransferItem `json:"DATA"`
	// optionals
	FilterRules            *[]FilterRule `json:"filter_rules,omitempty"`
	EncryptData            *bool         `json:"encrypt_data,omitempty"`             // default: false
	SyncLevel              *int          `json:"sync_level,omitempty"`               //
	VerifyChecksum         *bool         `json:"verify_checksum,omitempty"`          // default: false
	PreserveTimestamp      *bool         `json:"preserve_timestamp,omitempty"`       // default: false
	DeleteDestinationExtra *bool         `json:"delete_destination_extra,omitempty"` // default: false
	SkipSourceErrors       *bool         `json:"skip_source_errors,omitempty"`       // default: ?
	FailOnQuotaErrors      *bool         `json:"fail_on_quota_errors,omitempty"`     // default: ?
	SourceLocalUser        *string       `json:"source_local_user,omitempty"`
	DestinationLocalUser   *string       `json:"destination_local_user,omitempty"`
	// some BETA or experimental optional fields that are omitted:
	//  - RecursiveSymlinks
	//  - perf_cc, perf_p, perf_pp
	//  - perf_udt
}

type Delete struct {
	CommonTransfer
	Endpoint string `json:"endpoint"`
	Data     []DeleteItem
	// optionals
	Recursive      *bool   `json:"recursive,omitempty"`       // default: false, required if any item is a directory
	IgnoreMissing  *bool   `json:"ignore_missing,omitempty"`  // default: false
	InterpretGlobs *bool   `json:"interpret_globs,omitempty"` // default: false
	LocalUser      *string `json:"local_user,omitempty"`
}

// response structures

type SubmissionId struct {
	DataType string `json:"DATA_TYPE"`
	Value    string `json:"value"`
}

type TransferResult struct {
	DataType     string `json:"DATA_TYPE"`
	TaskId       string `json:"task_id"`
	SubmissionId string `json:"submission_id"`
	Code         string `json:"code"`
	Message      string `json:"message"`
	Resource     string `json:"resource"`
	RequestId    string `json:"requst_id"`
}

func getSubmissionId(client *http.Client) (submissionId string, err error) {
	if client == nil {
		return "", fmt.Errorf("client is nil.")
	}

	resp, err := client.Get(transferBaseUrl + "/submission_id")
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
func TransferSubmitGenericTask(client *http.Client, transfer Transfer) (result TransferResult, err error) {
	// get submission id for submission
	submission_id, err := getSubmissionId(client)
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
	resp, err := client.Post(
		transferBaseUrl+"/transfer",
		"application/json",
		bytes.NewBuffer(transferJSON),
	)
	if err != nil {
		return TransferResult{}, err
	}
	defer resp.Body.Close()

	fmt.Printf("Transfer req - status: %s\n", resp.Status)

	// read & return response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TransferResult{}, err
	}

	json.Unmarshal(body, &result)

	return result, nil
}

// submits a transfer task to copy a folder recursively.
// NOTE: all
func TransferFolderSync(client *http.Client, sourceEndpoint string, sourcePath string, destEndpoint string, destPath string) (TransferResult, error) {
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
	return TransferSubmitGenericTask(client, transfer)
}

func TransferListTasks(client *http.Client) {
	client.Get(transferBaseUrl + "/task_list")
}
