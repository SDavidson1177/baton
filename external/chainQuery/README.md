# Query Chains

This is a module that is intended to help test querying Cosmos chains (tendermint) for routing information.

## Fees

The fees command allows you to query fee information about a cosmos chain. It requires you to supply both a tendermint RPC endpoint (HTTP) and a block height. This command works by first iterating over all the transactions committed at a particular block height. The command will then return the total gas used, the average gas used per transaction and the total amount of transactions. This is a very simple query, and it would likely need to be improved to be practical in a real world setting (as different types of transactions have different gas usage, so you would need to distinguish which transactions to look at). Nonetheless, this tool can be used to get some quick estimations on gas utilization.

### Intuition for Baton

In the multihop setting, as long as you are targeting a destination chain and not performing anycast of some sort, it does not make sense to check gas prices at the destination. Furthermore, there is only one type of transaction that will triggered on the intermediate chains when a message is sent from source to destination. The relayer(s) may submit an *"update client"* transaction to the intermediate chains. This fortunately simplifies logic, as the sender is likely most concerned with the cost of updating the intermediate clients. Therefore, you can check the gas usage for this particular type of transaction in the intermediate chains.

This will require doing the following:

1. Calling the intermediat chain's Blockchain API (gRPC) and querying the IBC module for client information.
2. Request the heights of when the client for a particular channel was updated.
3. Query the Tendermint RPC endpoint for the gas usage of these transactions (implemented here)

### Usage

Example command:

```Bash
go run main.go fees "http://0.0.0.0:26657" 2066
```

Example response:

```Bash
Total Gas Used: 40457
Average Gas Used: 40457
Number of Txs: 1
```