package vote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// name to idetify transaction types
const MsgRoute = "gov"

//-----------------------------------------------------------
// MsgVote
type MsgVote struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Voter      sdk.AccAddress `json:"voter"`       //  address of the voter
	Option     string         `json:"option"`      //  option from OptionSet chosen by the voter
}

func NewMsgVote(voter sdk.AccAddress, proposalID uint64, option string) MsgVote {
	return MsgVote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
	}
}

// Implements Msg.
// nolint
func (msg MsgVote) Route() string { return MsgRoute }

// Implements Msg.
func (msg MsgVote) ValidateBasic() sdk.Error {
	if msg.ProposalID < 0 {
		return ErrUnknownProposal(DefaultCodespace, msg.ProposalID)
	}

	return nil
}

// Implements Msg.
func (msg MsgVote) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

func (msg MsgVote) String() string {
	return fmt.Sprintf("MsgVote{%v - %s}", msg.ProposalID, msg.Option)
}

// Implements Msg.
func (msg MsgVote) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg MsgVote) Type() string { return "vote" }
