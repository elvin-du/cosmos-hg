package vote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Key for getting a specific vote from the store
func KeyVote(proposalID uint64, voterAddr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("votes:%d:%d", proposalID, voterAddr))
}
