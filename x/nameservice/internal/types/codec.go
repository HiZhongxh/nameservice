package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetName{}, "nameservice/SetName", nil)
	cdc.RegisterConcrete(MsgBuyName{}, "nameservice/BuyName", nil)
	cdc.RegisterConcrete(MsgDeleteName{}, "nameservice/DeleteName", nil)
	cdc.RegisterConcrete(MsgAuctionName{}, "nameservice/AuctionName", nil)
	cdc.RegisterConcrete(MsgAuctionBid{}, "nameservice/AuctionBid", nil)
	cdc.RegisterConcrete(MsgAuctionReveal{}, "nameservice/AuctionReveal", nil)
}
