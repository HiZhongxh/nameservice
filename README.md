## nameservice
The project is develop with cosmos sdk v0.37.4, you can find origin from cosmos sdk-tutorials. The author add name auction function, and test ok in a local testnet contain of 2 nodes.

## support function:
1 buy/set name
2 auction/bid name

## prepare
1 one machine is ok, two is good(because we can test network and consensus)
2 go version > 1.13

I am developing in ubuntu 16.04

## 1 make
```
make
```

## 2 run
### 1 open one terminal in machine 2 (Optional)
```
nscli keys add hong
address: cosmos1wxct72sznuflutfwxrmcl5fhjnlapl5u0sn7pp
```

### 2  open one terminal in machine 1
```
# Initialize configuration files and genesis file
  # moniker is the name of your node
nsd init <moniker> --chain-id namechain


# Copy the `Address` output here and save it for later use
# [optional] add "--ledger" at the end to use a Ledger Nano S
nscli keys add jack

# Copy the `Address` output here and save it for later use
nscli keys add alice

# Copy the `Address` output here and save it for later use
nscli keys add bob

# Add both accounts, with coins to the genesis file
nsd add-genesis-account $(nscli keys show jack -a) 1000nametoken,100000000stake
nsd add-genesis-account $(nscli keys show alice -a) 1000nametoken,100000000stake
nsd add-genesis-account $(nscli keys show bob -a) 1000nametoken,100000000stake
// add hong to genesis account (Optional)
nsd add-genesis-account cosmos1wxct72sznuflutfwxrmcl5fhjnlapl5u0sn7pp 1000nametoken,100000000stake

# Configure your CLI to eliminate need for chain-id flag
nscli config chain-id namechain
nscli config output json
nscli config indent true
nscli config trust-node true

nsd gentx --name jack <or your key_name>

//After you have generated a genesis transaction, you will have to input the genTx into the genesis file, so that your nameservice chain is aware of the validators.
nsd collect-gentxs

// make sure your genesis file is correc
nsd validate-genesis

nsd start
```

### 3 open antoher terminal in mahcine 1
#### buy/set name
```
# First check the accounts to ensure they have funds
nscli query account $(nscli keys show jack -a)
nscli query account $(nscli keys show alice -a)

# Buy your first name using your coins from the genesis file
nscli tx nameservice buy-name jack.id 5nametoken --from jack

# Set the value for the name you just bought
nscli tx nameservice set-name jack.id 8.8.8.8 --from jack

# Try out a resolve query against the name you registered
nscli query nameservice resolve jack.id
# > 8.8.8.8

# Try out a whois query against the name you just registered
nscli query nameservice whois jack.id
# > {"value":"8.8.8.8","owner":"cosmos1l7k5tdt2qam0zecxrx78yuw447ga54dsmtpk2s","price":[{"denom":"nametoken","amount":"5"}]}

# Alice buys name from jack
nscli tx nameservice buy-name jack.id 10nametoken --from alice
```

#### auction/bid name
```
nscli tx nameservice auction-name jack.id 10nametoken 50 --from alice

nscli tx nameservice auction-bid jack.id 12nametoken --from jack
nscli query account $(nscli keys show jack -a)

nscli tx nameservice auction-bid jack.id 15nametoken --from bob
nscli query account $(nscli keys show bob -a)

nscli tx nameservice auction-bid jack.id 18nametoken --from jack
nscli query account $(nscli keys show jack -a)

nscli tx nameservice auction-bid jack.id 20nametoken --from bob
nscli query account $(nscli keys show bob -a)

// query auction struct
nscli query nameservice auction jack.id

// when auction is over
nscli tx nameservice auction-reveal jack.id --from alice

// query whois struct
nscli query nameservice whois jack.id

nscli query account $(nscli keys show jack -a)

nscli query account $(nscli keys show bob -a)
```

### Run second node on machine 2 (Optional)
Open terminal to run commands against that just created to install nsd and nscli

#### init use another moniker and same namechain
```
# Initialize configuration files and genesis file
  # moniker is the name of your node
nsd init <moniker-2> --chain-id namechain

# Configure your CLI to eliminate need for chain-id flag
nscli config chain-id namechain
nscli config output json
nscli config indent true
nscli config trust-node true
```
#### overwrite ~/.nsd/config/genesis.json with first node's genesis.json
#### change persistent_peers
```
vim /.nsd/config/config.toml
persistent_peers = "id@first_node_ip:26656"
run "nscli status" on first node to get id.
```
#### start this second node
```
nsd start

// register to be a validator
nscli tx staking create-validator --amount 100000000stake --moniker <moniker-2> --pubkey  $(nsd tendermint show-validator) --commission-max-change-rate 0.01 --commission-max-rate=0.2  --commission-rate=0.1 --min-self-delegation 1  --from hong
```
