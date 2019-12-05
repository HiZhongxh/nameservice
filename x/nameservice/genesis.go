package nameservice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	WhoisRecords []Whois `json:"whois_records"`
	AuctionRecords	[]Auction	`json:"auction_records"`
}

func NewGenesisState(whoIsRecords []Whois) GenesisState {
	return GenesisState{WhoisRecords: nil, AuctionRecords: nil}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.WhoisRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Owner", record.Value)
		}
		if record.Value == "" {
			return fmt.Errorf("invalid WhoisRecord: Owner: %s. Error: Missing Value", record.Owner)
		}
		if record.Price == nil {
			return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Price", record.Value)
		}
	}

	for _, record := range data.AuctionRecords {
		if record.Auctor == nil {
			return fmt.Errorf("invalid AuctionRecords: Value: %s. Error: Missing Auctor", record.Auctor)
		}
		if record.StartingPrice == nil {
			return fmt.Errorf("invalid AuctionRecords: Value: %s. Error: Missing StartingPrice", record.StartingPrice)
		}
		if record.DeadHeight == 0 {
			return fmt.Errorf("invalid AuctionRecords: Value: %s. Error: Missing DeadHeight", record.DeadHeight)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		WhoisRecords: []Whois{},
		AuctionRecords:	[]Auction{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.WhoisRecords {
		keeper.SetWhois(ctx, record.Value, record)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Whois
	iterator := k.GetNamesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		name := string(iterator.Key())
		whois := k.GetWhois(ctx, name)
		records = append(records, whois)
	}

	var auctionRecords []Auction
	iterator2 := k.GetAuctionNamesIterator(ctx)
	for ; iterator2.Valid(); iterator2.Next() {
		name := string(iterator2.Key())
		auction := k.GetAuction(ctx, name)
		auctionRecords = append(auctionRecords, auction)
	}

	return GenesisState{WhoisRecords: records, AuctionRecords: auctionRecords}
}
