package module

type AccountInfo struct {
	Account struct {
		Type               string `json:"@type"`
		BaseVestingAccount struct {
			BaseAccount struct {
				Address string `json:"address"`
				PubKey  struct {
					Type string `json:"@type"`
					Key  string `json:"key"`
				} `json:"pub_key"`
				AccountNumber string `json:"account_number"`
				Sequence      string `json:"sequence"`
			} `json:"base_account"`
			OriginalVesting []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"original_vesting"`
			DelegatedFree []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"delegated_free"`
			DelegatedVesting []interface{} `json:"delegated_vesting"`
			EndTime          string        `json:"end_time"`
		} `json:"base_vesting_account"`
		VestingPeriods []struct {
			StartTime string `json:"start_time"`
			Length    string `json:"length"`
			Amount    []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"amount"`
			ActionType int `json:"action_type"`
		} `json:"vesting_periods"`
	} `json:"account"`
}

type Account struct {
	Mnemonic string
	DelAddr  string
}
