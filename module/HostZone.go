package module

type HostZone struct {
	HostZone []struct {
		ChainID           string `json:"chain_id"`
		ConnectionID      string `json:"connection_id"`
		Bech32Prefix      string `json:"bech32prefix"`
		TransferChannelID string `json:"transfer_channel_id"`
		Validators        []struct {
			Name                 string      `json:"name"`
			Address              string      `json:"address"`
			Status               string      `json:"status"`
			CommissionRate       string      `json:"commission_rate"`
			DelegationAmt        string      `json:"delegation_amt"`
			Weight               string      `json:"weight"`
			InternalExchangeRate interface{} `json:"internal_exchange_rate"`
		} `json:"validators"`
		BlacklistedValidators []interface{} `json:"blacklisted_validators"`
		WithdrawalAccount     struct {
			Address     string        `json:"address"`
			Delegations []interface{} `json:"delegations"`
			Target      string        `json:"target"`
		} `json:"withdrawal_account"`
		FeeAccount struct {
			Address     string        `json:"address"`
			Delegations []interface{} `json:"delegations"`
			Target      string        `json:"target"`
		} `json:"fee_account"`
		DelegationAccount struct {
			Address     string        `json:"address"`
			Delegations []interface{} `json:"delegations"`
			Target      string        `json:"target"`
		} `json:"delegation_account"`
		RedemptionAccount struct {
			Address     string        `json:"address"`
			Delegations []interface{} `json:"delegations"`
			Target      string        `json:"target"`
		} `json:"redemption_account"`
		IbcDenom           string `json:"ibc_denom"`
		HostDenom          string `json:"host_denom"`
		LastRedemptionRate string `json:"last_redemption_rate"`
		RedemptionRate     string `json:"redemption_rate"`
		UnbondingFrequency string `json:"unbonding_frequency"`
		StakedBal          string `json:"staked_bal"`
		Address            string `json:"address"`
	} `json:"host_zone"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}
