package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/spoto/lending/x/lending/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

func SetupTestInput() (cdc *codec.Codec, ctx sdk.Context, ak auth.AccountKeeper, bk bank.Keeper, k Keeper) {
	cdc = testCodec()

	keys := sdk.NewKVStoreKeys(
		auth.StoreKey,
		params.StoreKey,
		types.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(params.TStoreKey)
	memDB := db.NewMemDB()
	ms := store.NewCommitMultiStore(memDB)

	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	}
	for _, tKey := range tKeys {
		ms.MountStoreWithDB(tKey, sdk.StoreTypeTransient, memDB)
	}
	_ = ms.LoadLatestVersion()

	ctx = sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	ak = auth.NewAccountKeeper(cdc, keys[auth.StoreKey], pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk = bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), nil)
	k = NewKeeper(keys[types.StoreKey], bk, cdc)

	return
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	cdc.Seal()

	return cdc
}
