package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	// "github.com/cosmos/cosmos-sdk/x/gov"
	"hashgard/x/vote"

	// "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagProposalID   = "proposal-id"
	flagTitle        = "title"
	flagDescription  = "description"
	flagProposalType = "type"
	flagDeposit      = "deposit"
	flagVoter        = "voter"
	flagOption       = "option"
	flagDepositer    = "depositer"
	flagStatus       = "status"
	flagNumLimit     = "limit"
	flagProposal     = "proposal"
)

// GetCmdVote implements creating a new vote command.
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "Vote for an active proposal, options: yes/no/no_with_veto/abstain",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			voterAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			proposalID := uint64(viper.GetInt64(flagProposalID))
			option := viper.GetString(flagOption)

			msg := vote.NewMsgVote(voterAddr, proposalID, option)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg}, false)
			}

			fmt.Printf("Vote[Voter:%s,ProposalID:%d,Option:%s]",
				voterAddr.String(), msg.ProposalID, msg.Option,
			)

			// Build and sign the transaction, then broadcast to a Tendermint
			// node.
			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal voting on")
	cmd.Flags().String(flagOption, "", "vote option {yes, no, no_with_veto, abstain}")

	return cmd
}
