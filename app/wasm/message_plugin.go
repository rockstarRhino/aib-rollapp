package wasm

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/AllInBetsCom/aib-rollapp/app/wasm/bindings"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

var mainnet = []string{""}
var testnet = []string{""}
var mainNetChainId = ""
var testNetChainId = ""
var moduleName = "aib"

func CustomMessageDecorator(bankkeeper bankkeeper.Keeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:    old,
			bankKeeper: bankkeeper,
		}
	}
}

type CustomMessenger struct {
	wrapped    wasmkeeper.Messenger
	bankKeeper bankkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {

		var aibMessages bindings.AibMessages
		if err := json.Unmarshal(msg.Custom, &aibMessages); err != nil {
			return nil, nil, sdkerrors.Wrap(err, "aib dispatch msg error")
		}

		if aibMessages.MsgMintToken != nil {

			return m.mintToken(ctx, contractAddr, aibMessages.MsgMintToken)
		}
		if aibMessages.MsgBurnToken != nil {

			return m.burnToken(ctx, contractAddr, aibMessages.MsgBurnToken)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) mintToken(ctx sdk.Context, contractAddr sdk.AccAddress, msgMintToken *bindings.MsgMintToken) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == mainNetChainId {
		if contractAddr.String() != mainnet[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == testNetChainId {
		if contractAddr.String() != testnet[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := mint(m.bankKeeper, ctx, contractAddr.String(), msgMintToken)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "error while minting tokens for aib")
	}
	return nil, nil, nil
}

func mint(bankKeeper bankkeeper.Keeper, ctx sdk.Context, contractAddr string,
	msgMintToken *bindings.MsgMintToken,
) error {

	mintCoin := sdk.NewCoins(sdk.NewCoin(msgMintToken.Denom, msgMintToken.Amount))
	err := bankKeeper.MintCoins(ctx, moduleName, mintCoin)
	if err != nil {
		return err
	}

	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleName, sdk.AccAddress(msgMintToken.MintToAddress), mintCoin)
	if err != nil {
		return err
	}

	return nil
}

func (m *CustomMessenger) burnToken(ctx sdk.Context, contractAddr sdk.AccAddress, msgBurnToken *bindings.MsgBurnToken) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == mainNetChainId {
		if contractAddr.String() != mainnet[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == testNetChainId {
		if contractAddr.String() != testnet[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := burn(m.bankKeeper, ctx, contractAddr.String(), msgBurnToken)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "error while burning tokens for aib")
	}
	return nil, nil, nil
}

func burn(bankKeeper bankkeeper.Keeper, ctx sdk.Context, contractAddr string,
	msgBurnToken *bindings.MsgBurnToken,
) error {

	burnCoin := sdk.NewCoins(sdk.NewCoin(msgBurnToken.Denom, msgBurnToken.Amount))
	err := bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(msgBurnToken.BurnFromAddress), moduleName, burnCoin)
	if err != nil {
		return err
	}
	err = bankKeeper.BurnCoins(ctx, moduleName, burnCoin)
	if err != nil {
		return err
	}

	return nil
}
