# Starknet indexer
This is an indexing layer for Starknet written in Golang that operates on top of the (Feeder) Gateway API and stores data in a Postgres database.  

It can be used in multiple ways and for various purposes:
- As a base component for the [DipDup Vertical](https://dipdup.io) — GraphQL federation providing a wide range of APIs for accessing both on-chain and off-chain data/metadata
- As a datasource for the [DipDup Framework](https://dipdup.io) — a Python SDK for building custom API backends for dapps
- As a standalone service — in a headless mode or with gRPC interface exposed
- As a library (go module) — most on-chain model definitions and indexing primitives can be reused

What you can build with DipDup:
- A custom API for your dapp to enable rich user interface & to save on RPC calls
- Any kind of aggregator: for DEXes, NFT marketplaces, etc
- Portfolio tracking tool
- A generic / specialized chain explorer
- An archive node for your L3 chain

## Features
- Blocks, events, and all operation types
- Internal operations are also indexed
- ERC20/ERC721/ERC1155 token transfers and balances (legacy tokens supported)
- Proxy and wallet contracts supported (known implementations from ArgentX, Braavos)
- Transaction calldata and event logs are pre-decoded (if ABI is provided by the node)
- Rollbacks are handled
- Database is partitioned for better performance
- Optional diagnostic mode for consistency checks
- gRPC interface and Hasura GQL engine integration

## Documentation
Check out the repository wiki to learn more about the indexer internals:
- [Build and run](https://github.com/dipdup-io/starknet-indexer/wiki/Configuration-and-building#building)
- [Configuration](https://github.com/dipdup-io/starknet-indexer/wiki/Configuration-and-building)
- [Database schema](https://github.com/dipdup-io/starknet-indexer/wiki/Database-structure)
- [gRPC protocol](https://github.com/dipdup-io/starknet-indexer/wiki/gRPC-protocol)

Also check out the [cmd/rpc_tester](https://github.com/dipdup-io/starknet-indexer/tree/master/cmd/rpc_tester) folder with simple events indexers for Starknet.ID and Loot Survivor.

## Public instances
Public deployments with reasonable rate limits are available for testing and prototyping:
- [Starknet mainnet](https://ide.dipdup.io/?resource=mainnet.starknet.dipdup.net/v1/graphql) `https://starknet-mainnet-gql.dipdup.net/v1/graphql`

## Notes
- Indexer works on top of the API provided by the sequencer node — it contains the most comprehensive data set, in particular classes and ABIs ordinary nodes do not always have; It's possible though to outsource several request types to the node API to reduce the load on the sequencer and speed up the indexing process, there's an option in the config for that.
- Currently pending blocks are not handled, therefore depending on the L2 block time (which in turn depends on the tx rate) you may have long delays in data updates.
- The API is currently in developer preview, request interface/response layout might change.  
- The underlying DB is not yet tuned for best performance, some queries might take a while to execute: we are working on improving that.

## Example queries

### Get token balances

Querying token balances for a given account.

```graphql
query GetTokenBalances {
  token_balance(
    where: {owner: {hash: {_eq: "\\x06ac597f8116f886fa1c97a23fa4e08299975ecaf6b598873ca6792b9bbfb678"}}}
  ) {
    owner_id
    balance
    token {
      metadata
      id
      type
      contract {
        hash
      }
    }
  }
}
```

Few things to note:
- Contract (wallet) addresses are not primary keys (for the sake of performance we use integers) so you need to either use filters or query by `owner_id`
- Hex strings are prefixed with `\\x` (instead of usual `0x`)
- Token type is enum (integer is used under the hood to save on storage), all possible values are listed in the docstring (check out the docs panel in the playground) or [wiki](https://github.com/dipdup-io/starknet-indexer/wiki/Database-structure#token) `1 - ERC20 | 2 - ERC721 | 3 - ERC1155`

### Get L1<>L2 messages

Querying messages and `l1_handler` operations for the StarkGate contract.

```graphql
query GetStarkGateMessages {
  message(
    where: {contract: {hash: {_eq: "\\x073314940630fd6dcda0d772d4c972c4e0a9946bef9dabf4ef84eda8ef542b82"}}}
    order_by: {id: desc}
    limit: 5
  ) {
    payload
    time
    to {
      hash
    }
  }
  l1_handler(
    where: {contract: {hash: {_eq: "\\x073314940630fd6dcda0d772d4c972c4e0a9946bef9dabf4ef84eda8ef542b82"}}}
    order_by: {id: desc}
    limit: 5
  ) {
    parsed_calldata
    time
    status
    entrypoint
  }
}
```

Notes on the response:
- You might noticed zero-prefixed addresses in the `message.to.hash` field – those are Ethereum L1 addresses
- `l1_handler.parsed_calldata` is the original calldata (also available) decoded according with the ABI provided by the sequencer node API
- `l1_handler.status` is also an enum, check out the [docstrings/wiki](https://github.com/dipdup-io/starknet-indexer/wiki/Database-structure#l1_handler) for details `unknown - 1 , not received - 2 , received - 3 , pending - 4 , rejected - 5 , accepted on l2 - 6 , accepted on l1 - 7`

### Get event logs

Querying Starknet.ID events.

```graphql
query GetStarknetIDs {
  event(
    where: {contract: {hash: {_eq: "\\x06ac597f8116f886fa1c97a23fa4e08299975ecaf6b598873ca6792b9bbfb678"}}, name: {_eq: "domain_to_addr_update"}}
    limit: 20
    order_by: {id: desc}
  ) {
    parsed_data
    time
  }
}

```

### Get token transfers

```graphql
query GetSithSwapTransfers {
  transfer(
    where: {to: {class: {hash: {_eq: "\\x07eb597ad7d9ba28ea1db162cdb99e265fe22bcb00e9b690e188c2203de9e005"}}}}
    limit: 50
    order_by: {id: desc}
  ) {
    amount
    from {
      hash
    }
    to {
      hash
    }
    token_id
    contract {
      hash
    }
    time
  }
}
```

### Get internal transactions

Querying execution trace for a SithSwap swap.

```graphql
query GetExecutionTrace {
  internal_tx(
    where: {hash: {_eq: "\\x07563ba09f924376edfdaf94b11941680867993d9caf271ae791ff4e89740177"}}
    order_by: {id: asc}
  ) {
    parsed_calldata
    parsed_result
    entrypoint
    caller {
      hash
    }
    call_type
    contract {
      hash
    }
  }
}
```

## About

DipDup Vertical for Starknet is a federated API including the following services:
- [x] Generic Starknet indexer
- [x] Starknet.ID indexer
- [ ] NFT metadata resolver

Project is supported by Starkware and Starknet Foundation via [OnlyDust platform](https://app.onlydust.xyz/projects/e1b6d080-7f15-4531-9259-10c3dae26848)
