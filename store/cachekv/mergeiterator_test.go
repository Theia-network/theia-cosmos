package cachekv_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/cachekv"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/store/dbadapter"
)

func TestEmitEventMangerInParentIterator(t *testing.T) {
	// initiate mock kvstore
	mem := dbadapter.Store{DB: dbm.NewMemDB()}
	kvstore := cachekv.NewStore(mem, types.NewKVStoreKey("CacheKvTest"), types.DefaultCacheSizeLimit)
	value := randSlice(defaultValueSizeBz)
	startKey := randSlice(32)

	keys := generateSequentialKeys(startKey, 3)
	for _, k := range keys {
		kvstore.Set(k, value)
	}

	// initialize mock mergeIterator
	eventManager := sdktypes.NewEventManager()
	parent := kvstore.Iterator(keys[0], keys[2])
	cache := kvstore.Iterator(nil, nil)
	for ; cache.Valid(); cache.Next() {
	}
	iter := cachekv.NewCacheMergeIterator(parent, cache, true, eventManager, types.NewKVStoreKey("CacheKvTest"))
	
	// get the next value
	iter.Value()

	// assert the resource access is still emitted correctly when the cache store is unavailable
	require.Equal(t, "access_type", eventManager.Events()[0].Attributes[0].Key)
	require.Equal(t, "read", eventManager.Events()[0].Attributes[0].Value)
	require.Equal(t, "store_key", eventManager.Events()[0].Attributes[1].Key)
	require.Equal(t, "CacheKvTest", eventManager.Events()[0].Attributes[1].Value)
}