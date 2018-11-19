package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Vote
type Vote struct {
	Voter      sdk.AccAddress `json:"voter"`       //  address of the voter
	ProposalID uint64         `json:"proposal_id"` //  proposalID of the proposal
	Option     string         `json:"option"`      //  option from OptionSet chosen by the voter
}
