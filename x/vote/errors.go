package vote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 5

	CodeUnknownProposal       sdk.CodeType = 1
	CodeInactiveProposal      sdk.CodeType = 2
	CodeAlreadyActiveProposal sdk.CodeType = 3
)

//----------------------------------------
// Error constructors

func ErrUnknownProposal(codespace sdk.CodespaceType, proposalID uint64) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownProposal, fmt.Sprintf("Unknown proposal with id %d", proposalID))
}
