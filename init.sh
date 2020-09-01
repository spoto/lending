rm -r ~/.lendingCLI
rm -r ~/.lendingD

% update to wherever is the gi bin folder in your machine
PATH="$PATH:$HOME/go/bin"

lendingD init mynode --chain-id lending

lendingCLI config keyring-backend test

lendingCLI keys add rich
lendingCLI keys add poor

lendingD add-genesis-account $(lendingCLI keys show rich -a) 1000foo,100000000stake
lendingD add-genesis-account $(lendingCLI keys show poor -a) 100foo

lendingCLI config chain-id lending
lendingCLI config output json
lendingCLI config indent true
lendingCLI config trust-node true
lendingD gentx --name rich --keyring-backend test
lendingD collect-gentxs
