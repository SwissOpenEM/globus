package globus

import "encoding/json"

func (pauseRule *PauseRuleLimited) UnmarshalJSON(data []byte) error {
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

	*pauseRule = PauseRuleLimited(*inner)
	return nil
}
