package globus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

// retrieve the list of successfully transfered files of a task
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

// retrieve the list of paths that were skipped because of the skip_source_errors flag being set to true
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

// provides details about why a task is paused - includes pause rules on source and destination collections
// and per-task pause flags set by collection activity managers
func TransferGetTaskPauseInfo(client *http.Client, taskID string) (info PauseInfoLimited, err error) {
	resp, err := client.Get(transferBaseUrl + "/task/" + taskID + "/pause_info")
	if err != nil {
		return PauseInfoLimited{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return PauseInfoLimited{}, fmt.Errorf("Non-Succesful Status: %d - %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PauseInfoLimited{}, err
	}

	err = json.Unmarshal(body, &info)
	return info, err
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
