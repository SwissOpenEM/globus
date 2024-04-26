package globus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const transferBaseUrl = "https://transfer.api.globusonline.org/v0.10"

type SubmissionId struct {
	DataType string `json:"DATA_TYPE"`
	Value    string `json:"value"`
}

type Link struct {
	DataType string `json:"DATA_TYPE"`
	Href     string `json:"href"`
	Rel      string `json:"rel"`
	// TODO: FINISH THIS!!!
}

type TransferResult struct {
	DataType     string `json:"DATA_TYPE"`
	Code         string `json:"code"`
	RequestId    string `json:"requst_id"`
	Resource     string `json:"resource"`
	SubmissionId string `json:"submission_id"`
	TaskId       string `json:"task_id"`
}

func TransferSubmitTask(client *http.Client, sourceEndpoint string, sourcePath string, destEndpoint string, destPath string) (err error) {
	// get submission id for submission
	resp, err := client.Get(transferBaseUrl + "/submission_id")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 || resp.Status != "OK" {
		return fmt.Errorf("unexpected status for submission id request: %d '%s' - %s", resp.StatusCode, resp.Status, string(body))
	}

	var result SubmissionId
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("could not parse response for submission id request")
	}
	if result.DataType != "submission_id" {
		return fmt.Errorf("incorrect value type returned for submission id request: %s", result.DataType)
	}

	submission_id := result.Value

	// submit request

	resp, err = client.Post(
		transferBaseUrl+"/transfer",
		"Content-Type: application/json",
		strings.NewReader(""),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TOOD: AND FINISH THIS!!!

	return nil
}

func TransferListTasks(client *http.Client) {
	client.Get(transferBaseUrl + "/task_list")
}
