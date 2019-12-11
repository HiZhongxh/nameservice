package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/HiZhongxh/nameservice/x/nameservice/internal/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	CoinKeeper bank.Keeper
	storeKey  sdk.StoreKey // Unexposed key to access store from sdk.Context
	storeMarketKey  sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, storeMarketKey  sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		CoinKeeper: 		coinKeeper,
		storeKey:   		storeKey,
		storeMarketKey:		storeMarketKey,
		cdc:        		cdc,
	}
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetWhois(ctx sdk.Context, name string, whois types.Whois) {
	if whois.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
}

// Delete the entire Whois metadata struct for a name
func (k Keeper) DeleteWhois(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(name))
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetWhois(ctx sdk.Context, name string) types.Whois {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(name)) {
		return types.NewWhois()
	}
	bz := store.Get([]byte(name))
	var whois types.Whois
	k.cdc.MustUnmarshalBinaryBare(bz, &whois)
	return whois
}

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	return k.GetWhois(ctx, name).Value
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	whois := k.GetWhois(ctx, name)
	return !whois.Owner.Empty()
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhois(ctx, name).Owner
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.GetWhois(ctx, name)
	whois.Owner = owner
	k.SetWhois(ctx, name, whois)
}

// GetPrice - gets the current price of a name.  If price doesn't exist yet, set to 1nametoken.
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetWhois(ctx, name).Price
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.Price = price
	k.SetWhois(ctx, name, whois)
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}


// Sets the entire Auction metadata struct for a name
func (k Keeper) SetAuction(ctx sdk.Context, name string, auction types.Auction) {
	if auction.Auctor.Empty() {
		return
	}
	store := ctx.KVStore(k.storeMarketKey)
	// done to use protobuf to marshal Auction struct, although it need many coding
	//bz ,err := auction.Serialize()
	//if err != nil {
	//	return
	//}
	//store.Set([]byte(name), bz)

	// MustMarshalBinaryBare can not marshal map struct
	//store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(auction))

	// MustMarshalJSON can marshal map struct, although it is not tight
	store.Set([]byte(name), k.cdc.MustMarshalJSON(auction))
}

// Delete the entire Auction metadata struct for a name
func (k Keeper) DeleteAuction(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeMarketKey)
	store.Delete([]byte(name))
}

// Gets the entire Auction metadata struct for a name
func (k Keeper) GetAuction(ctx sdk.Context, name string) types.Auction {
	store := ctx.KVStore(k.storeMarketKey)
	if !store.Has([]byte(name)) {
		return types.NewAuction()
	}
	bz := store.Get([]byte(name))
	var auction types.Auction

	//k.cdc.MustUnmarshalBinaryBare(bz, &auction)

	//auction.Deserialize(bz)

	k.cdc.MustUnmarshalJSON(bz, &auction)
	return auction
}

func (k Keeper) NewAuction(ctx sdk.Context, name string, auctor sdk.AccAddress, startingPrice sdk.Coins, height int64) {
	auction := k.GetAuction(ctx, name)
	auction.Auctor = auctor
	auction.StartingPrice = startingPrice
	auction.DeadHeight = ctx.BlockHeight() + height
	auction.Bids = make(map[string] types.Bid)
	k.SetAuction(ctx, name, auction)
}

// Get an iterator over all names in which the keys are the names and the values are the auction
func (k Keeper) GetAuctionNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeMarketKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

func (k Keeper) HasAuctor(ctx sdk.Context, name string) bool {
	auction := k.GetAuction(ctx, name)
	return !auction.Auctor.Empty()
}

func (k Keeper) GetAuctor(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetAuction(ctx, name).Auctor
}

func (k Keeper) GetAuctionStartingPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetAuction(ctx, name).StartingPrice
}

func (k Keeper) GetValidateHeight(ctx sdk.Context, name string) int64 {
	return k.GetAuction(ctx, name).DeadHeight
}

func (k Keeper) GetAuctionResult(ctx sdk.Context, name string) (sdk.AccAddress, sdk.Coins) {
	auction := k.GetAuction(ctx, name)

	var higestBid sdk.Coins = types.MinNamePrice
	var winner sdk.AccAddress
	for acc, b := range auction.Bids {
		if b.Bid.IsAllGT(higestBid) {
			higestBid = b.Bid
			winner, _ = sdk.AccAddressFromBech32(acc)
		}
	}

	return winner, higestBid
}

func (k Keeper) DelAuctionBid(ctx sdk.Context, name string, bidder sdk.AccAddress) {
	auction := k.GetAuction(ctx, name)
	delete(auction.Bids, bidder.String())
	k.SetAuction(ctx, name, auction)
}

func (k Keeper) GetAuctionBid(ctx sdk.Context, name string, bidder sdk.AccAddress) *types.Bid {
	if v, ok := k.GetAuction(ctx, name).Bids[bidder.String()]; ok {
		return &v
	} else {
		return nil
	}
}

func (k Keeper) GetAuctionBids(ctx sdk.Context, name string) map[string]types.Bid {
	return k.GetAuction(ctx, name).Bids
}

func (k Keeper) SetAuctionBid(ctx sdk.Context, name string, bidder sdk.AccAddress, bid sdk.Coins) {
	auction := k.GetAuction(ctx, name)
	b := types.Bid{
		Bid:	bid,
	}
	auction.Bids[bidder.String()] = b
	k.SetAuction(ctx, name, auction)
}