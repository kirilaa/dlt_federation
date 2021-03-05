package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/example/blog/x/blog/types"
	"strconv"
)

// GetUserCount get the total number of user
func (k Keeper) GetUserCount(ctx sdk.Context) int64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserCountKey))
	byteKey := types.KeyPrefix(types.UserCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseInt(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to int64
		panic("cannot decode count")
	}

	return count
}

// SetUserCount set the total number of user
func (k Keeper) SetUserCount(ctx sdk.Context, count int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserCountKey))
	byteKey := types.KeyPrefix(types.UserCountKey)
	bz := []byte(strconv.FormatInt(count, 10))
	store.Set(byteKey, bz)
}

// CreateUser creates a user with a new id and update the count
func (k Keeper) CreateUser(ctx sdk.Context, msg types.MsgCreateUser) {
	// Create the user
	count := k.GetUserCount(ctx)
	var user = types.User{
		Creator: msg.Creator,
		Id:      strconv.FormatInt(count, 10),
		Name:    msg.Name,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserKey))
	key := types.KeyPrefix(types.UserKey + user.Id)
	value := k.cdc.MustMarshalBinaryBare(&user)
	store.Set(key, value)

	// Update user count
	k.SetUserCount(ctx, count+1)
}

// SetUser set a specific user in the store
func (k Keeper) SetUser(ctx sdk.Context, user types.User) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserKey))
	b := k.cdc.MustMarshalBinaryBare(&user)
	store.Set(types.KeyPrefix(types.UserKey+user.Id), b)
}

// GetUser returns a user from its id
func (k Keeper) GetUser(ctx sdk.Context, key string) types.User {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserKey))
	var user types.User
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.UserKey+key)), &user)
	return user
}

// HasUser checks if the user exists
func (k Keeper) HasUser(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserKey))
	return store.Has(types.KeyPrefix(types.UserKey + id))
}

// GetUserOwner returns the creator of the user
func (k Keeper) GetUserOwner(ctx sdk.Context, key string) string {
	return k.GetUser(ctx, key).Creator
}

// DeleteUser deletes a user
func (k Keeper) DeleteUser(ctx sdk.Context, key string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserKey))
	store.Delete(types.KeyPrefix(types.UserKey + key))
}

// GetAllUser returns all user
func (k Keeper) GetAllUser(ctx sdk.Context) (msgs []types.User) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.UserKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.User
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return
}
