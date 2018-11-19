package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"hashgard/x/vote"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetCmdQueryVote(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-info",
		Short: "Query details of a single vote",
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID := uint64(viper.GetInt64(flagProposalID))
			voterAddr, err := sdk.AccAddressFromBech32(viper.GetString(flagVoter))
			if err != nil {

				return err
			}
			key := vote.KeyVote(proposalID, voterAddr)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			}

			var v vote.Vote
			err = cdc.UnmarshalBinaryLengthPrefixed(res, &v)
			if err != nil {
				return err
			}

			fmt.Println(v)

			return nil
		},
	}
	cmd.Flags().String(flagProposalID, "", "proposalID of proposal voting on")
	cmd.Flags().String(flagVoter, "", "bech32 voter address")

	return cmd
}
