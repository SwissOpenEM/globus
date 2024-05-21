package globus

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

type MarkerPaging struct {
	Marker     uint  `json:"marker"`
	NextMarker *uint `json:"next_marker,omitempty"`
}

// this could be represented by a TransferItem too, technically
type SuccessfulTransfer struct {
	DataType        string `json:"DATA_TYPE"`
	SourcePath      string `json:"source_path"`
	DestinationPath string `json:"destination_path"`
}

type SuccessfulTransfers struct {
	DataType string `json:"DATA_TYPE"`
	MarkerPaging
	Data []SuccessfulTransfer `json:"DATA"`
}

type SkippedError struct {
	TransferItem                    // Recursive will always be null
	ErrorCode                string `json:"error_code"`
	ErrorDetails             string `json:"error_details"`
	IsDirectory              bool   `json:"is_directory"`
	IsSymlink                bool   `json:"is_symlink"`
	IsDeleteDestinationExtra *bool  `json:"is_delete_destination_extra,omitempty"`
}

type SkippedErrors struct {
	DataType string `json:"DATA_TYPE"`
	MarkerPaging
	Data []SkippedError `json:"DATA"`
}

// TODO: finish up this part
type PauseRuleLimited struct {
	DataType               string  `json:"DATA_TYPE"`
	Id                     string  `json:"id"`
	Message                string  `json:"message"`
	StartTime              string  `json:"start_time"`
	EndpointId             string  `json:"endpoint_id"`
	EndpointDisplayName    string  `json:"endpoint_display_name"`
	IdentityId             *string `json:"identity_id,omitempty"`
	ModifiedTime           string  `json:"modified_time"`
	PauseLs                bool    `json:"pause_ls"`
	PauseMkdir             bool    `json:"pause_mkdir"`
	PauseSymlink           bool    `json:"pause_symlink"`
	PauseRename            bool    `json:"pause_rename"`
	PauseTaskDelete        bool    `json:"pause_task_delete"`
	PauseTaskTransferWrite bool    `json:"pause_task_transfer_write"`
	PauseTaskTransferRead  bool    `json:"pause_task_transfer_read"`
}

type PauseInfoLimited struct {
	DataType                     string             `json:"DATA_TYPE"`
	PauseRules                   []PauseRuleLimited `json:"pause_rules"`
	SourcePauseMessage           *string            `json:"source_pause_message,omitempty"`
	DestinationPauseMessage      *string            `json:"destination_pause_message,omitempty"`
	SourcePauseMessageShare      *string            `json:"source_pause_message_share,omitempty"`
	DestinationPauseMessageShare *string            `json:"destination_pause_message_share,omitempty"`
}
