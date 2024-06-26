package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/AllInBetsCom/aib-rollapp/x/gasless/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := genState.Validate(); err != nil {
		panic(err)
	}
	k.SetParams(ctx, genState.Params)

	for _, txToTankIDs := range genState.TxToGasTankIds {
		k.SetTxGTIDs(ctx, txToTankIDs)
	}

	k.SetLastGasTankID(ctx, genState.LastGasTankId)

	for _, tank := range genState.GasTanks {
		k.SetGasTank(ctx, tank)
	}

	for _, consumer := range genState.GasConsumers {
		k.SetGasConsumer(ctx, consumer)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:         k.GetParams(ctx),
		TxToGasTankIds: k.GetAllTxGTIDs(ctx),
		LastGasTankId:  k.GetLastGasTankID(ctx),
		GasTanks:       k.GetAllGasTanks(ctx),
		GasConsumers:   k.GetAllGasConsumers(ctx),
	}
}
