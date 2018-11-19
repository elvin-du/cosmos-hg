package vote

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgVote{}, "hashgard/MsgVote", nil)
}

var msgCdc = codec.New()
