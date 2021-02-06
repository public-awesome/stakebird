#!/bin/sh

CHAINID=$1
GENACCT=$2

if [ -z "$1" ]; then
  echo "Need to input chain id..."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input genesis account address..."
  exit 1
fi

# Build genesis file incl account for passed address
coins="10000000000stake,100000000000samoleans"
starsd init --chain-id $CHAINID $CHAINID
starsd keys add validator --keyring-backend="test"
starsd add-genesis-account validator $coins --keyring-backend="test"
starsd add-genesis-account $GENACCT $coins --keyring-backend="test"
starsd gentx --name validator --keyring-backend="test"
starsd collect-gentxs

# Set proper defaults and change ports
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.starsd/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.starsd/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.starsd/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.starsd/config/config.toml

# Start the chain
starsd start --pruning=nothing