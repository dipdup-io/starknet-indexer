# API

## Types

* `IntegerFilter`

```yaml
gt: Integer
gte: Integer
lt: Integer
lte: Integer
eq: Integer
neq: Integer
between:
    from: Integer # excluded
    to: Integer   # included
```


* `TimeFilter`

```yaml
gt: UNIX_TS
gte: UNIX_TS
lt: UNIX_TS
lte: UNIX_TS
between:
    from: UNIX_TS # excluded
    to: UNIX_TS   # included
```


* `Status` - integer with following enumeration

```yaml
Unknown: 1
NotReceived: 2
Received: 3
Pending: 4
Rejected: 5
AcceptedOnL2: 6
AcceptedOnL1: 7
```

* `StatusFilter`

```yaml
eq: Status
neq: Status
in: [Status]
notin: [Status]
```

* `VersionFilter`

```yaml
eq: Integer
neq: Integer
in: [Integer]
notin: [Integer]
```

* `StringFilter`

```yaml
eq: String
in: [String]
```

* `EqualityFilter`

```yaml
eq: String
neq: String
```

* `EntrypointType` - integer with following enumeration

```yaml
Unknown: 1
External: 2
Constructor: 3
L1Handler: 4
```

* `EntrypointTypeFilter`

```yaml
eq: EntrypointType
neq: EntrypointType
in: [EntrypointType]
notin: [EntrypointType]
```

* `CallType` - integer with following enumeration

```yaml
Unknown: 1
Call: 2
Delegate: 3
```

* `CallTypeFilter`

```yaml
eq: CallTypeFilter
neq: CallTypeFilter
in: [CallTypeFilter]
notin: [CallTypeFilter]
```

## Subscriptions

message, transfers, 
token_balance_update, storage_diff, class

### Blocks

`head` - subscribe on new block
`list` - receive all blocks (maybe with filter)

### Invoke

`get` - receives invokes by filters

Filters:

* `height` - filter by height type of `IntegerFilter` 

* `time` - filter by block time type of `TimeFilter`

* `status` - filter by block status type of `StatusFilter`

* `version` - filter by invoke version type of `VersionFilter`

* `contract` - filter by calling contract type of `StringFilter`

* `selector` - filter by invoked contract selector type of `EqualityFilter`

* `entrypoint` - filter by invoked entrypoint name type of `StringFilter`

TODO:
* `calldata` - filter by used calldata typeof `?`

TODO:
* `parsed_calldata` - filter by parsed calldata typeof `?`

### Declare

`get` - receives declares by filters

Filters:

* `height` - filter by height type of `IntegerFilter` 

* `time` - filter by block time type of `TimeFilter`

* `status` - filter by block status type of `StatusFilter`

* `version` - filter by invoke version type of `VersionFilter`

TODO: additional fields

### Deploy

`get` - receives deploys by filters

Filters:

* `height` - filter by height type of `IntegerFilter` 

* `time` - filter by block time type of `TimeFilter`

* `status` - filter by block status type of `StatusFilter`

* `class` - filter by class which contract was deployed type of `StringFilter`

TODO:
* `calldata` - filter by used calldata typeof `?`

TODO:
* `parsed_calldata` - filter by parsed calldata typeof `?`

### Deploy accounts

`get` - receives deploy account by filters

Filters:

* `height` - filter by height type of `IntegerFilter` 

* `time` - filter by block time type of `TimeFilter`

* `status` - filter by block status type of `StatusFilter`

* `class` - filter by class which account was deployed type of `StringFilter`

TODO:
* `calldata` - filter by used calldata typeof `?`

TODO:
* `parsed_calldata` - filter by parsed calldata typeof `?`


### L1 handlers

`get` - receives l1 handlers by filters

Filters:

* `height` - filter by height type of `IntegerFilter` 

* `time` - filter by block time type of `TimeFilter`

* `status` - filter by block status type of `StatusFilter`

* `contract` - filter by contract type of `StringFilter`

* `selector` - filter by invoked contract selector type of `EqualityFilter`

* `entrypoint` - filter by invoked entrypoint name type of `StringFilter`

TODO:
* `calldata` - filter by used calldata typeof `?`

TODO:
* `parsed_calldata` - filter by parsed calldata typeof `?`


### Internals

`get` - receives internal transactions by filters

Filters:

* `height` - filter by height type of `IntegerFilter` 

* `time` - filter by block time type of `TimeFilter`

* `status` - filter by block status type of `StatusFilter`

* `contract` - filter by contract type of `StringFilter`

* `caller` - filter by caller type of `StringFilter`

* `class` - filter by class type of `StringFilter`

* `selector` - filter by invoked contract selector type of `EqualityFilter`

* `entrypoint` - filter by invoked entrypoint name type of `StringFilter`

* `entrypoint_type` - filter by entrypoint type `EntrypointTypeFilter`

* `call_type` - filter by call type `CallTypeFilter`

TODO:
* `calldata` - filter by used calldata typeof `?`

TODO:
* `parsed_calldata` - filter by parsed calldata typeof `?`

### Fee

`get` - receives fee invocations by filters

Filters:

* `height` - filter by height type of `IntegerFilter` 

* `time` - filter by block time type of `TimeFilter`

* `status` - filter by block status type of `StatusFilter`

* `contract` - filter by contract type of `StringFilter`

* `from` - filter by from address of `StringFilter`

* `to` - filter by to address of `StringFilter`

### Events

`get` - receives events by filters

Filters:

* `height` - filter by height type of `IntegerFilter` 

* `time` - filter by block time type of `TimeFilter`

* `contract` - filter by contract type of `StringFilter`

* `from` - filter by from addressof `StringFilter`

* `name` - filter by event name of `EqualityFilter`

TODO:
* `keys` - filters by keys type of `?`

TODO:
* `parsed_data` - filter by parsed data typeof `?`

### Messages

`get` - receives messages by filters

Filters:

* `height` - filter by height type of `IntegerFilter` 

* `time` - filter by block time type of `TimeFilter`

* `contract` - filter by contract type of `StringFilter`

* `from` - filter by from address type of `StringFilter`

* `to` - filter by to address of `StringFilter`

* `selector` - filter by selector type of `EqualityFilter`



