package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"fmt"
	"sort"
	"errors"
	"github.com/HiZhongxh/nameservice/x/nameservice/internal/types/pb"
	"github.com/gogo/protobuf/proto"
)

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value	string			`json:"value"`
	Owner	sdk.AccAddress	`json:"owner"`
	Price	sdk.Coins		`json:"price"`
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

type Bid struct {
	Bid 	sdk.Coins		`json:"bid"`
}

type Auction struct {
	Auctor			sdk.AccAddress			`json:"auctor"`
	StartingPrice	sdk.Coins				`json:"starting_price"`
	DeadHeight		int64					`json:"dead_height"`
	Bids			map[string]Bid			`json:"bids"`
}

func NewAuction() Auction {
	return Auction{
		StartingPrice:	MinNamePrice,
		DeadHeight:		1,
	}
}

func (a Auction) proto() (pb.Auction, error) {
	var pbAuction pb.Auction
	// map is stored randomly, if consistency is needed(eg: clone state), we should sort firstly
	// but we don't need yet
	var keysBid []string
	for k := range a.Bids {
		keysBid = append(keysBid, k)
	}
	sort.Strings(keysBid)
	for _, k := range keysBid {
		bid := pb.Bid{
			Bidder:	k,
			Bid:	a.Bids[k].Bid.String(),
		}
		pbAuction.Bids = append(pbAuction.Bids, &bid)
	}
	//pbAuction.Bids = make(map[string]*pb.Bid)
	//for k, v := range a.Bids {
	//	b := pb.Bid{
	//		Bid: v.Bid.String(),
	//	}
	//	pbAuction.Bids[k] = &b
	//}

	pbAuction.Auctor = a.Auctor
	pbAuction.StartingPrice = a.StartingPrice.String()
	pbAuction.DeadHeight = a.DeadHeight

	return pbAuction, nil
}

func (a Auction) Serialize() ([]byte, error) {
	p, err := a.proto()
	if err != nil {
		return nil, err
	}
	data, err := proto.Marshal(&p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a *Auction) Deserialize(data []byte) error {
	if len(data) == 0 {
		return errors.New("can not deserialize empty byte array")
	}

	var pbAuction pb.Auction
	err := proto.Unmarshal(data, &pbAuction)
	if err != nil {
		return err
	}

	a.Auctor = pbAuction.Auctor
	a.StartingPrice, err = sdk.ParseCoins(pbAuction.StartingPrice)
	if err != nil {
		return err
	}
	a.DeadHeight = pbAuction.DeadHeight
	a.Bids = make(map[string]Bid)
	for _, b := range pbAuction.Bids {
		var bid Bid
		bid.Bid, err = sdk.ParseCoins(b.Bid)
		if err != nil {
			return err
		}

		a.Bids[b.Bidder] = bid
	}
	//for k, v := range pbAuction.Bids {
	//	var bid Bid
	//	bid.Bid, err = sdk.ParseCoins(v.Bid)
	//	if err != nil {
	//		return err
	//	}
	//	a.Bids[k] = bid
	//}

	return nil
}

func (a Auction) String() string {
	bids := ModuleCdc.MustMarshalJSON(a.Bids)
	return strings.TrimSpace(fmt.Sprintf(`Auctor: %s
StartingPrice: %s
DeadHeight %d
Bids %s`, a.Auctor, a.StartingPrice, a.DeadHeight, string(bids)))

	//return string(ModuleCdc.MustMarshalJSON(a))
}