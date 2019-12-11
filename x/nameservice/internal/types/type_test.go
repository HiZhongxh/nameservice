package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
	"fmt"
)

func TestEnDecode(t *testing.T) {
	jack := sdk.AccAddress("cosmos13pnn6qmhms2e08kajtpv3qzjjvwq3kkyg2r6y7")
	alice := sdk.AccAddress("cosmos19dn5f6ewn404umhxlf58y0wnruu38zmwnwzztt")
	bob := sdk.AccAddress("cosmos15gy4j083cq9y4z7mcsg7rdm33zd482sdfqnmta")
	bids := make(map[string]Bid)
	bids[string(alice)] = Bid{
		Bid:	sdk.Coins{sdk.NewCoin("test", sdk.NewInt(15))},
	}
	bids[string(bob)] = Bid{
		Bid:	sdk.Coins{sdk.NewCoin("test", sdk.NewInt(25))},
	}
	auction := Auction{
		Auctor:			jack,
		StartingPrice:	sdk.Coins{sdk.NewCoin("test", sdk.NewInt(10))},
		DeadHeight:		100,
		Bids:			bids,
	}

	bz, err := auction.Serialize()
	if err != nil {
		t.Error("auction serialize failed")
	}

	var dauction Auction
	err = dauction.Deserialize(bz)
	if err != nil {
		t.Error("auction deserialize failed")
	}

	fmt.Println(dauction.String())
	//require.Equal(t, auction, dauction)
}


type Students struct {
	students 	map[int]string
}

func (s Students) GetMap() map[int]string {
	return s.students
}

func TestMap(t *testing.T) {
	stu := Students{
		students:	make(map[int]string),
	}

	stu.students[0] = "xiao"
	stu.students[1] = "hong"

	stu1 := stu.GetMap()
	for _, v := range stu1 {
		fmt.Println("",  v)
	}
}