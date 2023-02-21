# StarkNet generic indexer

This is a basic layer of the indexing stack called "DipDup Vertical" being created for the StarkNet ecosystem.  
With this open-source indexing solution we aim to enable developers with flexible APIs for accessing blockchain data, as well as with a powerful framework for creating custom APIs for particular dapps.

## What is DipDup

DipDup is a modular framework for creating selective indexers and featureful backends for decentralized applications. It eliminates the boilerplate and takes care of most of the indexing-specific things allowing developers to focus on the business logic, reducing time-to-market.

https://dipdup.io

DipDup is used in production for 1.5 years by 15+ teams in the Tezos ecosystem, and we are currently expanding to EVM, planning to support 70 compatible networks by Q2 2023.

### DipDup Verticals

DipDup Vertical is a group of indexing services working on top of each other and wrapped with a common API facade. In the result developers get access to all kinds of APIs, from generic blockchain account and transaction data, to NFT metadata and aggregated DEX quotes, all through a single GraphQL endpoint.

We currently have a single vertical implemented for Tezos chain, it includes token/contract metadata indexers, mempool indexer, domains indexer, DEX aggregator, analytical backend, search engine, and dapps listing.

## Key features

DipDup was developed to tackle the complexity of multi-domain indexing in an environment of rapidly growing data volumes and tight deadlines for providing new features. Its design allows to significantly reduce costs to develop and maintain custom API solutions.

### Fast feature-to-market

Forementioned DipDup Verticals are implemented as stateful microservices that form a hierarchy where one sub-indexer uses output of another sub-indexer as input.

This design allows us to handle complexity and to be able to ship new features in a short time without sacrificing code quality and reliability of the entire system, because:
- There's no need to re-sync the entire system, just the service that was updated and optionally its dependent services; or create a new service if feature is large and standalone enough
- We can use any tech stack and delegate tasks to different teams, as long as the common indexing flow and communication interface are implemented
- We can re-use existing indexing modules providing already processed data

In order to solve the problem of having too many APIs with such approach, we use API federation pattern, namely GraphQL federation which allows to combine multiple API endpoints under a single facade. Additionally it enables cross-schema relations, which makes the resulting API more coupled and functional.

### Deep customization

Unlike competing solutions like TheGraph, DipDup does not restrict developers in his choice of database engine, API gateway, or any other third-party integrations.

By default, it works with PostgreSQL and Hasura GQL engine, but you can also use PostgREST, send data to Kafka, write your own API endpoints, or just leave your backend headless. Similarly, you can choose to use TimescaleDB if you deal with time series, or any other DB that suits your needs.

You can also query any datasources, run background jobs, maybe have some user-generated content mixed with blockchain data, enable user authentication, and many more. Basically, you can do whatever you need to build a fully-fledged backend for your dapp.

### Enterprise grade

DipDup is open-source, written in Python, works with a variety of time-proven DB and API engines, and natively supports common observability and maintenance services like Prometheus and Sentry.

You also have a relatively low vendor lock-in risk with such a setup.

### Low development costs

It is very easy to start building on DipDup, as it's Python-based, accompanied with comprehensive docs and templates. You can also get help in our developer community.

DipDup is a framework, and it gently guides you towards the right path by restricting file organisation, by providing to-be-implemented stubs and other hints.

DipDup significantly reduces boilerplate associated with querying and decoding data, handling network and chain issues. What left is just core logic you need to implement.

## Project milestones

Our project has two major milestones:
- Implement a minimal DipDup vertical for the StarkNet rollup chain
- Add support for StarkNet to the DipDup framework

### Milestone 1

**ETA**: end of Q2
**Deliverables**: a generic indexer providing high-level data representation via GraphQL API

Requirements:
- open-source, documented self-deployment setup
- index blocks, handle reorgs (rollbacks)
- index contract classes and instances, handle proxies
- index transactions: declare, deploy, invoke, l1_handler; decode calldata
- index events; decode event payloads
- index ERC token transfers
- expose GraphQL API

### Milestone 2

**ETA**: circa mid-Q3
**Deliverables**: DipDup-compatible datasource and full support for Cairo contracts in DipDup framework

Requirements:
- expose DipDup-compatible API set (for code generation and selective querying) via gRPC
- add support for the new gRPC datasource to the DipDup framework
- provide `operations`, `events`, and `token_transfers` index kinds at the framework level
- allow to filter by contract class in addition to filtering by contract address
- provide a template that can be used via scaffolding CLI
- provide an extensive documentation and quickstart example
