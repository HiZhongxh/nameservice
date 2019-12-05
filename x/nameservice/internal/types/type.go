package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"fmt"
)

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// MinNamePrice is Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// NewWhois returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}


type Auction struct {
	Auctor			sdk.AccAddress			`json:"auctor"`
	StartingPrice	sdk.Coins				`json:"starting_price"`
	DeadHeight		int64					`json:"dead_height"`
	Bids			map[string]sdk.Coins	`json:"bids"`
}

func NewAuction() Auction {
	return Auction{
		StartingPrice:	MinNamePrice,
		DeadHeight:		1,
	}
}

func (a Auction) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Auctor: %s StartingPrice: %s DeadHeight %d`,
		a.Auctor, a.StartingPrice, a.DeadHeight))
}