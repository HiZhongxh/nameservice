package nameservice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/HiZhongxh/nameservice/x/nameservice/internal/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		//case MsgDeleteName:
		//	return handleMsgDeleteName(ctx, keeper, msg)
		case MsgAuctionName:
			return handleMsgAuctionName(ctx, keeper, msg)
		case MsgAuctionBid:
			return handleMsgAuctionBid(ctx, keeper, msg)
		case MsgAuctionReveal:
			return handleMsgAuctionReveal(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg types.MsgSetName) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
	return sdk.Result{} // return
}

// Handle a message to buy name
func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg types.MsgBuyName) sdk.Result {
	if keeper.GetPrice(ctx, msg.Name).IsAllGTE(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.Name) {
		//err := keeper.CoinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid) // If not, throw an error
		//if err != nil {
		//	return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		//}
		sdk.ErrUnauthorized("The name has owner").Result() // If not, throw an error
	} else {
		_, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return sdk.Result{}
}

// Handle a message to delete name
func handleMsgDeleteName(ctx sdk.Context, keeper Keeper, msg types.MsgDeleteName) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	keeper.DeleteWhois(ctx, msg.Name) // If so, delete the entire Whois metadata struct for a name
	return sdk.Result{} // return
}

// Handle a message to auction name
func handleMsgAuctionName(ctx sdk.Context, keeper Keeper, msg types.MsgAuctionName) sdk.Result {
	if !msg.Auctor.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	if keeper.HasAuctor(ctx, msg.Name) {
		return sdk.ErrUnauthorized("The name is aucting").Result() // If not, throw an error
	}

	keeper.NewAuction(ctx, msg.Name, msg.Auctor, msg.StartingPrice, msg.DeadHeight)
	return sdk.Result{} // return
}

// Handle a message to bid in auction
func handleMsgAuctionBid(ctx sdk.Context, keeper Keeper, msg types.MsgAuctionBid) sdk.Result {
	currentHeight := ctx.BlockHeight()
	validateHeight := keeper.GetValidateHeight(ctx, msg.Name)
	if currentHeight > validateHeight {
		return sdk.ErrUnauthorized("The auction is not existed or invalidated").Result() // If not, throw an error
	}

	if msg.Buyer.Equals(keeper.GetAuctor(ctx, msg.Name)) {
		return sdk.ErrUnauthorized("auctor can't bid in his auction").Result() // If not, throw an error
	}

	if keeper.GetAuctionStartingPrice(ctx, msg.Name).IsAllGTE(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
		return sdk.ErrInsufficientCoins("Bid is less than starting price").Result() // If not, throw an error
	}

	var hadPay sdk.Coins
	var err bool
	oldBid := keeper.GetAuctionBid(ctx, msg.Name, msg.Buyer)
	if oldBid == nil {
		hadPay = msg.Bid
	} else {
		if oldBid.Bid.IsAllGTE(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
			return sdk.ErrInsufficientCoins("Bid is less than previous price").Result() // If not, throw an error
		}

		hadPay, err = msg.Bid.SafeSub(oldBid.Bid)
		if err != false {
			return sdk.ErrInsufficientCoins("safe sub failed when calculate expense").Result() // If not, throw an error
		}
	}

	if keeper.HasAuctor(ctx, msg.Name) {
		_, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Buyer, hadPay) // If so, deduct the Bid amount from the sender
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		fmt.Println("The auction is not existed or invalidated")
		return sdk.ErrUnauthorized("The auction is not existed or invalidated").Result() // If not, throw an error
	}

	keeper.SetAuctionBid(ctx, msg.Name, msg.Buyer, msg.Bid)
	return sdk.Result{}
}

// Handle a message to reveal auction
func handleMsgAuctionReveal(ctx sdk.Context, keeper Keeper, msg types.MsgAuctionReveal) sdk.Result {
	auctor := keeper.GetAuctor(ctx, msg.Name)

	if !msg.Auctor.Equals(auctor) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	currentHeight := ctx.BlockHeight()
	validateHeight := keeper.GetValidateHeight(ctx, msg.Name)
	if currentHeight < validateHeight {
		return sdk.ErrUnauthorized("The auction is still aucting").Result() // If not, throw an error
	}
	winner, bid := keeper.GetAuctionResult(ctx, msg.Name)
	if !winner.Empty() {
		keeper.CoinKeeper.AddCoins(ctx, auctor, bid)
		bids := keeper.GetAuctionBids(ctx, msg.Name)
		for acc, b := range bids {
			bidder, _ := sdk.AccAddressFromBech32(acc)
			fmt.Println("bidder: ", bidder)
			if !bidder.Equals(winner) {
				keeper.CoinKeeper.AddCoins(ctx, bidder, b.Bid)
			}
		}
	}

	keeper.SetOwner(ctx, msg.Name, winner)
	keeper.SetPrice(ctx, msg.Name, bid)
	keeper.DeleteAuction(ctx, msg.Name)
	return sdk.Result{} // return
}