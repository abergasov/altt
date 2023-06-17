package entities

type Balance struct {
	Chain           Chain  `json:"chain"`
	ChainName       string `json:"chain_name"`
	Token           Token  `json:"token"`
	TokenBalance    string `json:"token_balance"`
	TokenBalanceWei string `json:"token_balance_wei"`
}
