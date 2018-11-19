##Cosmos-SDK开发指南

* 拷贝cosmos-sdk/examples/basecoin,并修改成你项目的名称
> `cp -rf $GOPATH/src/github.com/cosmos/cosmos-sdk/examples/basecoin $GOPATH/src/hashgard`

* 修改项目文件名和函数名等
>Mac环境下的命令：
>`cd $GOPATH/src/hashgard
mv cmd/basecli cmd/hashgardcli
mv cmd/basecoind cmd/hashgardd
grep -rl "github.com/cosmos/cosmos-sdk/examples/basecoin" .|xargs sed -i "" 's/github.com\/cosmos\/cosmos-sdk\/examples\/basecoin/hashgard/g'
grep -rl "basecoin" .|xargs sed -i "" 's/basecoin/hashgard/g'
grep -rl "Basecoin" .|xargs sed -i "" 's/Basecoin/Hashgard/g'
grep -rl "basecli" .|xargs sed -i "" 's/basecli/hashgardcli/g'`

### 设计并实现节点对消息的处理

* 创建自定义的x模块
>`mkdir -p x/vote x/vote/client/cli x/vote/client/rest
cd x/vote
touch handler.go msgs.go errors.go keeper.go keeper_keys.go types.go codec.go`

* 设计msg模型，并按照接口规范实现以下接口：
> msg.go Tx内的msg结构
    ```
    消息体接口规范：
    type Msg interface {
        // Return the message Route.
        // Must be alphanumeric or empty.
        // Must correspond to name of message handler (XXX).
        Route() string
    
        // ValidateBasic does a simple validation check that
        // doesn't require access to any other information.
        ValidateBasic() error
    
        // Get the canonical byte representation of the Msg.
        // This is what is signed.
        GetSignBytes() []byte
    
    	// GetSigners returns the addrs of signers that must sign.
        // CONTRACT: All signatures must be present to be valid.
        // CONTRACT: Returns addrs in some deterministic order.
        GetSigners() []AccAddress
    }
    ```

* 按照Error接口规范实现自定义的错误信息：
> errors.go 自定义错误信息
    ```
    // sdk Error type
    type Error interface {
    	// Implements cmn.Error
    	// Error() string
    	// Stacktrace() cmn.Error
    	// Trace(offset int, format string, args ...interface{}) cmn.Error
    	// Data() interface{}
    	cmnError
    
    	// convenience
    	TraceSDK(format string, args ...interface{}) Error
    
    	// set codespace
    	WithDefaultCodespace(CodespaceType) Error
    
    	Code() CodeType
    	Codespace() CodespaceType
    	ABCILog() string
    	ABCICode() ABCICodeType
    	Result() Result
    	QueryResult() abci.ResponseQuery
    }
    ```

* 实现消息处理程序
> handler.go 对消息进行路由
    ```
    // Handle all "gov" type messages.
    func NewHandler(keeper Keeper) sdk.Handler {
    	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
    		switch msg := msg.(type) {
    		case MsgVote:
    			return handleMsgVote(ctx, keeper, msg)
    		default:
    			errMsg := "Unrecognized gov msg type"
    			return sdk.ErrUnknownRequest(errMsg).Result()
    		}
    	}
    }
    
    func handleMsgVote(ctx sdk.Context, keeper Keeper, msg MsgVote) sdk.Result {
    	err := keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
    	if err != nil {
    		return err.Result()
    	}
    
    	proposalIDBytes := keeper.cdc.MustMarshalBinaryBare(msg.ProposalID)
    	resTags := sdk.NewTags(
    		tags.Action, tags.ActionVote,
    		tags.Voter, []byte(msg.Voter.String()),
    		tags.ProposalID, proposalIDBytes,
    	)
    	return sdk.Result{
    		Tags: resTags, //自定义事件
    	}
    }
    ```

> keeper.go 逻辑处理
```
// Governance Keeper
type Keeper struct {
	// The reference to the CoinKeeper to modify balances
	ck bank.Keeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Reserved codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ck bank.Keeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: codespace,
	}
}

// Adds a vote on a specific proposal
func (keeper Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, option string) sdk.Error {
	vote := Vote{
		ProposalID: proposalID,
		Voter:      voterAddr,
		Option:     option,
	}
	keeper.setVote(ctx, proposalID, voterAddr, vote)

	return nil
}

func (keeper Keeper) setVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, vote Vote) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(KeyVote(proposalID, voterAddr), bz)
}
```

* 注册消息handler

> app.go: 添加数据库，注册消息handler等
```
func NewHashgardApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *HashgardApp {
	// create your application type
	var app = &HashgardApp{
		cdc:        cdc,
		BaseApp:    bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...),
		keyMain:    sdk.NewKVStoreKey("main"),
		keyAccount: sdk.NewKVStoreKey("acc"),
		keyIBC:     sdk.NewKVStoreKey("ibc"),
		keyVote:    sdk.NewKVStoreKey("vote"),
	}
	
	.....
	
	app.voteKeeper = vote.NewKeeper(app.cdc, app.keyVote, app.bankKeeper, app.RegisterCodespace(vote.DefaultCodespace))
	// register message routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.bankKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.bankKeeper)).
		AddRoute("gov", vote.NewHandler(app.voteKeeper))
    
    
    app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC,app.keyVote)

}

// MakeCodec creates a new codec codec and registers all the necessary types
// with the codec.
func MakeCodec() *codec.Codec {
	cdc := codec.New()
	vote.RegisterCodec(cdc)
	cdc.Seal()
	return cdc
}
```

### 客户端
####cli客户端
* 产生交易
> mkdir -p x/client/cli 
> touch tx.go
```
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

```
把产生vote交易的cmd添加到根cmd中
> cmd/hashgardcli/main.go
```
func main() {
	rootCmd.AddCommand(
		client.PostCommands(
			votecmd.GetCmdVote(cdc),
		)...)
}
```

* 查询交易
> query.go
```
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

```

把查询交易的cmd添加到根cmd中
> cmd/hashgardcli/main.go
```
func main() {
// add query/post commands (custom to binary)
	rootCmd.AddCommand(
		client.GetCommands(
			votecmd.GetCmdQueryVote("vote", cdc),
		)...)
}
```
