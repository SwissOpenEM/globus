package globus

import "encoding/json"

func (pauseInfo *PauseRuleLimited) UnmarshalJSON(data []byte) error {
	type innerRule PauseRuleLimited
	inner := &innerRule{
		PauseLs:         true,
		PauseMkdir:      true,
		PauseSymlink:    true,
		PauseRename:     true,
		PauseTaskDelete: true,
	}

	if err := json.Unmarshal(data, inner); err != nil {
		return err
	}

	*pauseInfo = PauseRuleLimited(*inner)
	return nil
}
