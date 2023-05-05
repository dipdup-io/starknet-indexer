# Starknet indexer
This is an indexing layer for Starknet written in Golang that operates on top of the (Feeder) Gateway API and stores data in a Postgres database.  

It can be used in multiple ways and for various purposes:
- As a base component for the [DipDup Vertical](https://dipdup.io) — GraphQL federation providing a wide range of APIs for accessing both on-chain and off-chain data/metadata
- As a datasource for the [DipDup Framework](https://dipdup.io) — a Python SDK for building custom API backends for dapps
- As a standalone service — in a headless mode or with gRPC interface exposed
- As a library (go module) — most on-chain model definitions and indexing primitives can be reused

## Features
- Blocks, events, and all operation types
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

Also check out the [cmd/rpc_tester](https://github.com/dipdup-io/starknet-indexer/tree/master/cmd/rpc_tester) folder with simple events indexers for StarknetID and Loot Survivor.

## Public instances
Public deployments with reasonable rate limits are available for testing and prototyping:
- [Starknet mainnet](https://play.dipdup.io/?endpoint=https://starknet-mainnet-gql.dipdup.net/v1/graphql)

## About
Project is supported by Starkware and Starknet Foundation via [OnlyDust platform](https://app.onlydust.xyz/projects/e1b6d080-7f15-4531-9259-10c3dae26848)
