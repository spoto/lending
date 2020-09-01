rm -r ~/.lendingCLI
rm -r ~/.lendingD

% update to wherever is the gi bin folder in your machine
PATH="$PATH:$HOME/go/bin"

lendingD init mynode --chain-id lending

lendingCLI config keyring-backend test

lendingCLI keys add me
lendingCLI keys add you

lendingD add-genesis-account $(lendingCLI keys show me -a) 1000foo,100000000stake
lendingD add-genesis-account $(lendingCLI keys show you -a) 1foo

lendingCLI config chain-id lending
lendingCLI config output json
lendingCLI config indent true
lendingCLI config trust-node true
lendingD gentx --name me --keyring-backend test
lendingD collect-gentxs
