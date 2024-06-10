package bindings

import "cosmossdk.io/math"

type RollAppMessages struct {
	MsgMintToken *MsgMintToken `json:"msg_mint_token,omitempty"`
	MsgBurnToken *MsgBurnToken `json:"msg_burn_token,omitempty"`
}

type MsgBurnToken struct {
	Denom           string   `json:"denom"`
	Amount          math.Int `json:"amount"`
	BurnFromAddress string   `json:"burn_from_address"`
}

type MsgMintToken struct {
	Denom         string   `json:"denom"`
	Amount        math.Int `json:"amount"`
	MintToAddress string   `json:"mint_to_address"`
}
