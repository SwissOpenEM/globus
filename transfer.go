package globus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const transferBaseUrl = "https://transfer.api.globusonline.org/v0.10"

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

	if resp.StatusCode == 409 {
		return Result{}, fmt.Errorf("history was deleted")
	}

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

// lists task's events
// NOTE: the history gets deleted after 30 days
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

// TODO: test!
func TransferGetTaskSuccessfulTransfers(client *http.Client, taskID string, marker uint) (transfers SuccessfulTransfers, err error) {
	req, err := http.NewRequest(http.MethodGet, transferBaseUrl+"/task/"+taskID+"/successful_transfers", nil)
	if err != nil {
		return SuccessfulTransfers{}, err
	}

	q := req.URL.Query()
	q.Add("marker", fmt.Sprint(marker))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return SuccessfulTransfers{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return SuccessfulTransfers{}, fmt.Errorf("Non-Successful Status: %d - %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SuccessfulTransfers{}, err
	}

	err = json.Unmarshal(body, &transfers)
	return transfers, err
}

// TODO: test!
func TransferGetTaskSkippedErrors(client *http.Client, taskID string, marker uint) (skips SkippedErrors, err error) {
	req, err := http.NewRequest(http.MethodGet, transferBaseUrl+"/task/"+taskID+"/skipped_errors", nil)
	if err != nil {
		return SkippedErrors{}, err
	}

	q := req.URL.Query()
	q.Add("marker", fmt.Sprint(marker))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return SkippedErrors{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return SkippedErrors{}, fmt.Errorf("Non-Successful Status: %d - %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SkippedErrors{}, err
	}

	err = json.Unmarshal(body, &skips)
	return skips, err
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
