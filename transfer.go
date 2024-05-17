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

type FatalError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// TODO: confirm that this works correctly with replies
type Task struct {
	DataType                       string        `json:"DATA_TYPE"`
	TaskId                         string        `json:"task_id"`
	Type                           string        `json:"type"`
	Status                         string        `json:"status"`
	FatalError                     *FatalError   `json:"fatal_error,omitempty"`
	Label                          string        `json:"label"`
	OwnerId                        string        `json:"owner_id"`
	RequestTime                    string        `json:"request_time"`              // ISO8601
	CompletionTime                 *string       `json:"completion_time,omitempty"` // null if hasn't finished
	Deadline                       string        `json:"deadline"`
	SourceEndpointId               string        `json:"source_endpoint_id"`
	SourceEndpointDisplayName      string        `json:"source_endpoint_display_name"`
	DestinationEndpointId          *string       `json:"destination_endpoint_id,omitempty"` // null for delete tasks
	DestinationEndpointDisplayName *string       `json:"destination_endpoint_display_name,omitempty"`
	SyncLevel                      *int          `json:"sync_level,omitempty"`
	EncryptData                    bool          `json:"encrypt_data"`
	VerifyChecksum                 bool          `json:"verify_checksum"`
	DeleteDestinationExtra         bool          `json:"delete_destination_extra"`
	RecursiveSymlinks              *string       `json:"recursive_symlinks,omitempty"` // always null for delete tasks
	PreserveTimestamp              bool          `json:"preserve_timestamp"`
	SkipSourceErrors               bool          `json:"skip_source_errors"`
	FailOnQuotaErrors              bool          `json:"fail_on_quota_errors"`
	Command                        string        `json:"command"`
	HistoryDeleted                 bool          `json:"history_deleted"`
	Faults                         int           `json:"faults"`
	Files                          int           `json:"files"`       // no. of files affected by task (can grow w/ recursion)
	Directories                    int           `json:"directories"` // no. of directories affected by task (can grow w/ recursion)
	Symlinks                       int           `json:"symlinks"`    // no. of *kept* symlinks
	FilesSkipped                   *int          `json:"files_skipped,omitempty"`
	FilesTransferred               int           `json:"files_transferred"`
	SubtasksTotal                  int           `json:"subtasks_total"`
	SubtasksPending                int           `json:"subtasks_pending"`
	SubtasksRetrying               int           `json:"subtasks_retrying"`
	SubtasksSucceeded              int           `json:"subtasks_succeeded"`
	SubtasksExpired                int           `json:"subtasks_expired"`
	SubtasksCanceled               int           `json:"subtasks_canceled"`
	SubtasksFailed                 int           `json:"subtasks_failed"`
	SubtasksSkippedErrors          int           `json:"subtasks_skipped_errors"`
	BytesTransferred               int           `json:"bytes_transferred"`
	BytesChecksummed               int           `json:"bytes_checksummed"`
	EffectiveBytesPerSecond        int           `json:"effective_bytes_per_second"`
	NiceStatus                     *string       `json:"nice_status,omitempty"` // "OK" or "Queued" -> task is fine, otherwise some error
	NiceStatusShortDescription     string        `json:"nice_status_short_description"`
	NiceStatusExpiresIn            int           `json:"nice_status_expires_in"`
	CanceledByAdmin                *string       `json:"canceled_by_admin,omitempty"` // if the task was canceled by either collection's activity manager, otherwise null
	CanceledByAdminMessage         string        `json:"canceled_by_admin_message"`   // contains the message of cancelation set by activity manager
	IsPaused                       bool          `json:"is_paused"`
	FilterRules                    *[]FilterRule `json:"filter_rules,omitempty"` // can be null
	SourceLocalUser                *string       `json:"source_local_user,omitempty"`
	SourceLocalUserStatus          string        `json:"source_local_user_status"`
	DestinationLocalUser           *string       `json:"destination_local_user,omitempty"`
	DestinationLocalUserStatus     string        `json:"destination_local_user_status"`
	SourceBasePath                 *string       `json:"source_base_path,omitempty"`
	DestinationBasePath            *string       `json:"destination_base_path,omitempty"`

	// deprecated fields:
	//SourceEndpoint string `json:"source_endpoint`
	//Username string `json:"username"`
	//NiceStatusDetails string `json:"nice_status_details"`
}

type TaskList struct {
	DataType string `json:"DATA_TYPE"`
	Length   int    `json:"length"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Total    int    `json:"total"`
	Data     []Task `json:"Data"`
}

type Result struct {
	DataType  string `json:"DATA_TYPE"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
	Resource  string `json:"resource"`
}

type Event struct {
	DataType    string `json:"DATA_TYPE"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Details     string `json:"details"`
	IsError     bool   `json:"is_error"`
	Time        string `json:"time"`
}

type EventList struct {
	Data     []Event `json:"DATA"`
	DataType string  `json:"DATA_TYPE"`
	Limit    uint    `json:"limit"`
	Offset   uint    `json:"offset"`
	Total    uint    `json:"total"`
}

func getSubmissionId(client *http.Client) (submissionId string, err error) {
	if client == nil {
		return "", fmt.Errorf("client is nil")
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
func TransferPostTask(client *http.Client, transfer Transfer) (result TransferResult, err error) {
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

	err = json.Unmarshal(body, &result)
	return result, err
}

// fetches a list of transfer tasks from Globus Transfer API
// NOTE: the results are paginated using "offset" and "limit"
func TransferGetTaskList(client *http.Client, offset uint, limit uint) (taskList TaskList, err error) {
	req, err := http.NewRequest(http.MethodGet, transferBaseUrl+"/task_list", nil)
	if err != nil {
		return TaskList{}, err
	}

	q := req.URL.Query()
	q.Add("offset", fmt.Sprint(offset))
	q.Add("limit", fmt.Sprint(limit))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return TaskList{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return TaskList{}, fmt.Errorf("Non-Successful Status: %d - %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TaskList{}, err
	}

	err = json.Unmarshal(body, &taskList)
	return taskList, err
}

// fetches a specific transfer task from Globus Transfer API by its ID
func TransferGetTaskByID(client *http.Client, taskID string) (task Task, err error) {
	resp, err := client.Get(transferBaseUrl + "/task/" + taskID)
	if err != nil {
		return Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return Task{}, fmt.Errorf("Non-Successful Status: %d - %s", resp.StatusCode, resp.Status)
	}

	// read & return response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Task{}, err
	}

	err = json.Unmarshal(body, &task)
	return task, err
}

// cancels a task using its id
func TransferCancelTaskByID(client *http.Client, taskID string) (result Result, err error) {
	resp, err := client.Post(transferBaseUrl+"/task/"+taskID+"/cancel", "", nil)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return Result{}, fmt.Errorf("Non-Successful Status: %d - %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	err = json.Unmarshal(body, &result)
	return result, err
}

// removes a globus task
// NOTE: this can be only used under specific conditions: task must be associated with a
// a high assurance collection, must be either SUCCEEDED or FAILED.
func TransferRemoveTaskByID(client *http.Client, taskID string) (result Result, err error) {
	resp, err := client.Post(transferBaseUrl+"/task/"+taskID+"/remove", "", nil)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return Result{}, fmt.Errorf("Non-Successful Status: %d - %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	err = json.Unmarshal(body, &result)
	return result, err
}

func TransferGetTaskEventList(client *http.Client, taskID string, offset uint, limit uint) (eventList EventList, err error) {
	resp, err := client.Get(transferBaseUrl + "/task/" + taskID + "/event_list")
	if err != nil {
		return EventList{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return EventList{}, fmt.Errorf("Non-Successful Status: %d - %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return EventList{}, err
	}

	err = json.Unmarshal(body, &eventList)
	return eventList, err
}

// submits a transfer task to copy a folder recursively.
// NOTE: the transfer follows all default params (aside from recursivity)
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
	return TransferPostTask(client, transfer)
}

func TransferListTasks(client *http.Client) {
	client.Get(transferBaseUrl + "/task_list")
}

// creates a list of scopes to access data on the specified Globus endpoints
func TransferDataAccessScopeCreator(collectionIDs []string) (scopes []string) {
	for _, collectionID := range collectionIDs {
		if collectionID == "" {
			continue
		}
		scopes = append(scopes, "urn:globus:auth:scope:transfer.api.globus.org:all[*https://auth.globus.org/scopes/"+collectionID+"/data_access]")
	}

	return scopes
}
