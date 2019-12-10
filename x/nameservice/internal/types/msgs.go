package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

// MsgBuyName defines the BuyName message
type MsgBuyName struct {
	Name 	string
	Bid 	sdk.Coins
	Buyer	sdk.AccAddress
}

// NewMsgBuyName is the constructor function for MsgBuyName
func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name: name,
		Bid:    bid,
		Buyer:  buyer,
	}
}

// Route should return the name of the module
func (msg MsgBuyName) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBuyName) Type() string { return "buy_name" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBuyName) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBuyName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}


// MsgSetName defines a SetName message
type MsgSetName struct {
	Name string
	Value  string
	Owner  sdk.AccAddress
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName{
		Name: name,
		Value:  value,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgSetName) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetName) Type() string { return "set_name"}

// ValidateBasic runs stateless checks on the message
func (msg MsgSetName) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Name and/or Value cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}


// MsgDeleteName defines a DeleteName message
type MsgDeleteName struct {
	Name  string         `json:"name"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgDeleteName is a constructor function for MsgDeleteName
func NewMsgDeleteName(name string, owner sdk.AccAddress) MsgDeleteName {
	return MsgDeleteName{
		Name:  name,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgDeleteName) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDeleteName) Type() string { return "delete_name" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDeleteName) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDeleteName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDeleteName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}


type MsgAuctionName struct {
	Name			string			`json:"name"`
	StartingPrice 	sdk.Coins		`json:"starting_price"`
	DeadHeight		int64			`json:"dead_height"`
	Auctor			sdk.AccAddress	`json:"auctor"`
}

func NewMsgAuctionName(name string, starting_price sdk.Coins, deadHeight int64, owner sdk.AccAddress) MsgAuctionName {
	return MsgAuctionName{
		Name:			name,
		StartingPrice:	starting_price,
		DeadHeight:		deadHeight,
		Auctor:			owner,
	}
}

// Route should return the name of the module
func (msg MsgAuctionName) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAuctionName) Type() string { return "auction_name" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAuctionName) ValidateBasic() sdk.Error {
	if msg.Auctor.Empty() {
		return sdk.ErrInvalidAddress(msg.Auctor.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}
	if !msg.StartingPrice.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Starting price must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAuctionName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAuctionName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Auctor}
}


// MsgAuctionBid defines the AuctionBid message
type MsgAuctionBid struct {
	Name 	string				`json:"name"`
	Bid 	sdk.Coins			`json:"bid"`
	Buyer	sdk.AccAddress		`json:"buyer"`
}

// NewMsgAuctionBid is the constructor function for AuctionBid
func NewMsgAuctionBid(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgAuctionBid {
	return MsgAuctionBid{
		Name: name,
		Bid:    bid,
		Buyer:  buyer,
	}
}

// Route should return the name of the module
func (msg MsgAuctionBid) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAuctionBid) Type() string { return "auction_bid" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAuctionBid) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAuctionBid) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAuctionBid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

type MsgAuctionReveal struct {
	Name 	string			`json:"name"`
	Auctor	sdk.AccAddress	`json:"auctor"`
}

func NewMsgAuctionReveal(name string, auctor sdk.AccAddress) MsgAuctionReveal {
	return MsgAuctionReveal{
		Name:	name,
		Auctor:	auctor,
	}
}

// Route should return the name of the module
func (msg MsgAuctionReveal) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAuctionReveal) Type() string { return "auction_reveal" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAuctionReveal) ValidateBasic() sdk.Error {
	if msg.Auctor.Empty() {
		return sdk.ErrInvalidAddress(msg.Auctor.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAuctionReveal) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAuctionReveal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Auctor}
}