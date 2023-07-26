# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/entity_filters.proto](#github-com_dipdup-io_starknet-indexer_pkg_grpc_proto_entity_filters-proto)
    - [AddressFilter](#proto-AddressFilter)
    - [DeclareFilters](#proto-DeclareFilters)
    - [DeployAccountFilters](#proto-DeployAccountFilters)
    - [DeployAccountFilters.ParsedCalldataEntry](#proto-DeployAccountFilters-ParsedCalldataEntry)
    - [DeployFilters](#proto-DeployFilters)
    - [DeployFilters.ParsedCalldataEntry](#proto-DeployFilters-ParsedCalldataEntry)
    - [EventFilter](#proto-EventFilter)
    - [EventFilter.ParsedDataEntry](#proto-EventFilter-ParsedDataEntry)
    - [FeeFilter](#proto-FeeFilter)
    - [FeeFilter.ParsedCalldataEntry](#proto-FeeFilter-ParsedCalldataEntry)
    - [InternalFilter](#proto-InternalFilter)
    - [InternalFilter.ParsedCalldataEntry](#proto-InternalFilter-ParsedCalldataEntry)
    - [InvokeFilters](#proto-InvokeFilters)
    - [InvokeFilters.ParsedCalldataEntry](#proto-InvokeFilters-ParsedCalldataEntry)
    - [L1HandlerFilter](#proto-L1HandlerFilter)
    - [L1HandlerFilter.ParsedCalldataEntry](#proto-L1HandlerFilter-ParsedCalldataEntry)
    - [MessageFilter](#proto-MessageFilter)
    - [StorageDiffFilter](#proto-StorageDiffFilter)
    - [TokenBalanceFilter](#proto-TokenBalanceFilter)
    - [TokenFilter](#proto-TokenFilter)
    - [TransferFilter](#proto-TransferFilter)
  
- [github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/enum.proto](#github-com_dipdup-io_starknet-indexer_pkg_grpc_proto_enum-proto)
    - [CallType](#proto-CallType)
    - [EntrypointType](#proto-EntrypointType)
    - [Status](#proto-Status)
  
- [github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/filters.proto](#github-com_dipdup-io_starknet-indexer_pkg_grpc_proto_filters-proto)
    - [BetweenInteger](#proto-BetweenInteger)
    - [BytesArray](#proto-BytesArray)
    - [BytesFilter](#proto-BytesFilter)
    - [EnumFilter](#proto-EnumFilter)
    - [EqualityFilter](#proto-EqualityFilter)
    - [EqualityIntegerFilter](#proto-EqualityIntegerFilter)
    - [IntegerArray](#proto-IntegerArray)
    - [IntegerFilter](#proto-IntegerFilter)
    - [StringArray](#proto-StringArray)
    - [StringFilter](#proto-StringFilter)
    - [TimeFilter](#proto-TimeFilter)
  
- [github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/indexer.proto](#github-com_dipdup-io_starknet-indexer_pkg_grpc_proto_indexer-proto)
    - [Bytes](#proto-Bytes)
    - [JsonSchema](#proto-JsonSchema)
    - [JsonSchemaItem](#proto-JsonSchemaItem)
    - [SubscribeRequest](#proto-SubscribeRequest)
    - [Subscription](#proto-Subscription)
  
    - [IndexerService](#proto-IndexerService)
  
- [github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/response.proto](#github-com_dipdup-io_starknet-indexer_pkg_grpc_proto_response-proto)
    - [Address](#proto-Address)
    - [Block](#proto-Block)
    - [Declare](#proto-Declare)
    - [Deploy](#proto-Deploy)
    - [DeployAccount](#proto-DeployAccount)
    - [EndOfBlock](#proto-EndOfBlock)
    - [Event](#proto-Event)
    - [Fee](#proto-Fee)
    - [Internal](#proto-Internal)
    - [Invoke](#proto-Invoke)
    - [L1Handler](#proto-L1Handler)
    - [StarknetMessage](#proto-StarknetMessage)
    - [StorageDiff](#proto-StorageDiff)
    - [Token](#proto-Token)
    - [TokenBalance](#proto-TokenBalance)
    - [Transfer](#proto-Transfer)
  
- [Scalar Value Types](#scalar-value-types)



<a name="github-com_dipdup-io_starknet-indexer_pkg_grpc_proto_entity_filters-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/entity_filters.proto



<a name="proto-AddressFilter"></a>

### AddressFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |
| only_starknet | [bool](#bool) |  |  |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-DeclareFilters"></a>

### DeclareFilters



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |
| time | [TimeFilter](#proto-TimeFilter) |  |  |
| status | [EnumFilter](#proto-EnumFilter) |  |  |
| version | [EnumFilter](#proto-EnumFilter) |  |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-DeployAccountFilters"></a>

### DeployAccountFilters



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |
| time | [TimeFilter](#proto-TimeFilter) |  |  |
| status | [EnumFilter](#proto-EnumFilter) |  |  |
| class | [BytesFilter](#proto-BytesFilter) |  |  |
| parsed_calldata | [DeployAccountFilters.ParsedCalldataEntry](#proto-DeployAccountFilters-ParsedCalldataEntry) | repeated |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-DeployAccountFilters-ParsedCalldataEntry"></a>

### DeployAccountFilters.ParsedCalldataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="proto-DeployFilters"></a>

### DeployFilters



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |
| time | [TimeFilter](#proto-TimeFilter) |  |  |
| status | [EnumFilter](#proto-EnumFilter) |  |  |
| class | [BytesFilter](#proto-BytesFilter) |  |  |
| parsed_calldata | [DeployFilters.ParsedCalldataEntry](#proto-DeployFilters-ParsedCalldataEntry) | repeated |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-DeployFilters-ParsedCalldataEntry"></a>

### DeployFilters.ParsedCalldataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="proto-EventFilter"></a>

### EventFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |
| time | [TimeFilter](#proto-TimeFilter) |  |  |
| contract | [BytesFilter](#proto-BytesFilter) |  |  |
| from | [BytesFilter](#proto-BytesFilter) |  |  |
| name | [StringFilter](#proto-StringFilter) |  |  |
| parsed_data | [EventFilter.ParsedDataEntry](#proto-EventFilter-ParsedDataEntry) | repeated |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-EventFilter-ParsedDataEntry"></a>

### EventFilter.ParsedDataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="proto-FeeFilter"></a>

### FeeFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |
| time | [TimeFilter](#proto-TimeFilter) |  |  |
| status | [EnumFilter](#proto-EnumFilter) |  |  |
| contract | [BytesFilter](#proto-BytesFilter) |  |  |
| caller | [BytesFilter](#proto-BytesFilter) |  |  |
| class | [BytesFilter](#proto-BytesFilter) |  |  |
| selector | [EqualityFilter](#proto-EqualityFilter) |  |  |
| entrypoint | [StringFilter](#proto-StringFilter) |  |  |
| entrypoint_type | [EnumFilter](#proto-EnumFilter) |  |  |
| call_type | [EnumFilter](#proto-EnumFilter) |  |  |
| parsed_calldata | [FeeFilter.ParsedCalldataEntry](#proto-FeeFilter-ParsedCalldataEntry) | repeated |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-FeeFilter-ParsedCalldataEntry"></a>

### FeeFilter.ParsedCalldataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="proto-InternalFilter"></a>

### InternalFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |
| time | [TimeFilter](#proto-TimeFilter) |  |  |
| status | [EnumFilter](#proto-EnumFilter) |  |  |
| contract | [BytesFilter](#proto-BytesFilter) |  |  |
| caller | [BytesFilter](#proto-BytesFilter) |  |  |
| class | [BytesFilter](#proto-BytesFilter) |  |  |
| selector | [EqualityFilter](#proto-EqualityFilter) |  |  |
| entrypoint | [StringFilter](#proto-StringFilter) |  |  |
| entrypoint_type | [EnumFilter](#proto-EnumFilter) |  |  |
| call_type | [EnumFilter](#proto-EnumFilter) |  |  |
| parsed_calldata | [InternalFilter.ParsedCalldataEntry](#proto-InternalFilter-ParsedCalldataEntry) | repeated |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-InternalFilter-ParsedCalldataEntry"></a>

### InternalFilter.ParsedCalldataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="proto-InvokeFilters"></a>

### InvokeFilters



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |
| time | [TimeFilter](#proto-TimeFilter) |  |  |
| status | [EnumFilter](#proto-EnumFilter) |  |  |
| version | [EnumFilter](#proto-EnumFilter) |  |  |
| contract | [BytesFilter](#proto-BytesFilter) |  |  |
| selector | [EqualityFilter](#proto-EqualityFilter) |  |  |
| entrypoint | [StringFilter](#proto-StringFilter) |  |  |
| parsed_calldata | [InvokeFilters.ParsedCalldataEntry](#proto-InvokeFilters-ParsedCalldataEntry) | repeated |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-InvokeFilters-ParsedCalldataEntry"></a>

### InvokeFilters.ParsedCalldataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="proto-L1HandlerFilter"></a>

### L1HandlerFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |
| time | [TimeFilter](#proto-TimeFilter) |  |  |
| status | [EnumFilter](#proto-EnumFilter) |  |  |
| contract | [BytesFilter](#proto-BytesFilter) |  |  |
| selector | [EqualityFilter](#proto-EqualityFilter) |  |  |
| entrypoint | [StringFilter](#proto-StringFilter) |  |  |
| parsed_calldata | [L1HandlerFilter.ParsedCalldataEntry](#proto-L1HandlerFilter-ParsedCalldataEntry) | repeated |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-L1HandlerFilter-ParsedCalldataEntry"></a>

### L1HandlerFilter.ParsedCalldataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="proto-MessageFilter"></a>

### MessageFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |
| time | [TimeFilter](#proto-TimeFilter) |  |  |
| contract | [BytesFilter](#proto-BytesFilter) |  |  |
| from | [BytesFilter](#proto-BytesFilter) |  |  |
| to | [BytesFilter](#proto-BytesFilter) |  |  |
| selector | [EqualityFilter](#proto-EqualityFilter) |  |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-StorageDiffFilter"></a>

### StorageDiffFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |
| contract | [BytesFilter](#proto-BytesFilter) |  |  |
| key | [EqualityFilter](#proto-EqualityFilter) |  |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-TokenBalanceFilter"></a>

### TokenBalanceFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner | [BytesFilter](#proto-BytesFilter) |  |  |
| contract | [BytesFilter](#proto-BytesFilter) |  |  |
| token_id | [StringFilter](#proto-StringFilter) |  |  |






<a name="proto-TokenFilter"></a>

### TokenFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| contract | [BytesFilter](#proto-BytesFilter) |  |  |
| owner | [BytesFilter](#proto-BytesFilter) |  |  |
| type | [EnumFilter](#proto-EnumFilter) |  |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |






<a name="proto-TransferFilter"></a>

### TransferFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [IntegerFilter](#proto-IntegerFilter) |  |  |
| time | [TimeFilter](#proto-TimeFilter) |  |  |
| contract | [BytesFilter](#proto-BytesFilter) |  |  |
| from | [BytesFilter](#proto-BytesFilter) |  |  |
| to | [BytesFilter](#proto-BytesFilter) |  |  |
| token_id | [StringFilter](#proto-StringFilter) |  |  |
| id | [IntegerFilter](#proto-IntegerFilter) |  |  |





 

 

 

 



<a name="github-com_dipdup-io_starknet-indexer_pkg_grpc_proto_enum-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/enum.proto


 


<a name="proto-CallType"></a>

### CallType
Call type of transactions

| Name | Number | Description |
| ---- | ------ | ----------- |
| CALL_TYPE_RESERVED | 0 | unused |
| CALL_TYPE_UNKNOWN | 1 | used only if entity has unknown call type for the system |
| CALL_TYPE_CALL | 2 | call |
| CALL_TYPE_DELEGATE | 3 | delegate call |



<a name="proto-EntrypointType"></a>

### EntrypointType
Entrypoint type of transactions

| Name | Number | Description |
| ---- | ------ | ----------- |
| ENTRYPOINT_TYPE_RESERVED | 0 | unused |
| ENTRYPOINT_TYPE_UNKNOWN | 1 | used only if entity has unknown entrypoint type for the system |
| ENTRYPOINT_TYPE_EXTERNAL | 2 | external entrypoint type |
| ENTRYPOINT_TYPE_CONSTRUCTOR | 3 | constructor entrypoint type |
| ENTRYPOINT_TYPE_L1_HANDLER | 4 | l1 handler entrypoint type |



<a name="proto-Status"></a>

### Status
Block status

| Name | Number | Description |
| ---- | ------ | ----------- |
| STATUS_RESERVED | 0 | unused |
| STATUS_UNKNOWN | 1 | used only if entity has unknown status for the system |
| STATUS_NOT_RECEIVED | 2 | not received |
| STATUS_RECEIVED | 3 | received |
| STATUS_PENDING | 4 | pending |
| STATUS_REJECTED | 5 | rejected |
| STATUS_ACCEPTED_ON_L2 | 6 | accepted on L2 |
| STATUS_ACCEPTED_ON_L1 | 7 | accepted on L1 |


 

 

 



<a name="github-com_dipdup-io_starknet-indexer_pkg_grpc_proto_filters-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/filters.proto



<a name="proto-BetweenInteger"></a>

### BetweenInteger
Between unsigned interger filter. Equals to SQL expression: `x BETWEEN from AND to`.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| from | [uint64](#uint64) |  | from value |
| to | [uint64](#uint64) |  | to value |






<a name="proto-BytesArray"></a>

### BytesArray
Wrapper over bytes array for using `repeated` option


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| arr | [bytes](#bytes) | repeated | array |






<a name="proto-BytesFilter"></a>

### BytesFilter
Set of bytes filters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| eq | [bytes](#bytes) |  | equals |
| in | [BytesArray](#proto-BytesArray) |  | check the value is in array `x IN (\x00, \x0010)` |






<a name="proto-EnumFilter"></a>

### EnumFilter
Set of filters for enumerations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| eq | [uint64](#uint64) |  | equals |
| neq | [uint64](#uint64) |  | not equals |
| in | [IntegerArray](#proto-IntegerArray) |  | check the value is in array `x IN (1,2,3,4)` |
| notin | [IntegerArray](#proto-IntegerArray) |  | check the value is not in array `x NOT IN (1,2,3,4)` |






<a name="proto-EqualityFilter"></a>

### EqualityFilter
Equality filters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| eq | [string](#string) |  | equals |
| neq | [string](#string) |  | not equals |






<a name="proto-EqualityIntegerFilter"></a>

### EqualityIntegerFilter
Equality filters for integer values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| eq | [uint64](#uint64) |  | equals |
| neq | [uint64](#uint64) |  | not equals |






<a name="proto-IntegerArray"></a>

### IntegerArray
Wrapper over integer array for using `repeated` option


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| arr | [uint64](#uint64) | repeated | array |






<a name="proto-IntegerFilter"></a>

### IntegerFilter
Set of integer filters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| gt | [uint64](#uint64) |  | greater than |
| gte | [uint64](#uint64) |  | greater than or equals |
| lt | [uint64](#uint64) |  | less than |
| lte | [uint64](#uint64) |  | less than or equals |
| eq | [uint64](#uint64) |  | equals |
| neq | [uint64](#uint64) |  | not equals |
| between | [BetweenInteger](#proto-BetweenInteger) |  | between |






<a name="proto-StringArray"></a>

### StringArray
Wrapper over string array for using `repeated` option


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| arr | [string](#string) | repeated | array |






<a name="proto-StringFilter"></a>

### StringFilter
Set of string filters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| eq | [string](#string) |  | equals |
| in | [StringArray](#proto-StringArray) |  | check the value is in array `x IN (a, abc)` |






<a name="proto-TimeFilter"></a>

### TimeFilter
Set of time filters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| gt | [uint64](#uint64) |  | greater than |
| gte | [uint64](#uint64) |  | greater than or equals |
| lt | [uint64](#uint64) |  | less than |
| lte | [uint64](#uint64) |  | less than or equals |
| between | [BetweenInteger](#proto-BetweenInteger) |  | between |





 

 

 

 



<a name="github-com_dipdup-io_starknet-indexer_pkg_grpc_proto_indexer-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/indexer.proto



<a name="proto-Bytes"></a>

### Bytes
Bytes array


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | array |






<a name="proto-JsonSchema"></a>

### JsonSchema
Json schema entity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| functions | [JsonSchemaItem](#proto-JsonSchemaItem) | repeated | list of functions json schema |
| l1_handlers | [JsonSchemaItem](#proto-JsonSchemaItem) | repeated | list of l1 handlers json schema |
| constructors | [JsonSchemaItem](#proto-JsonSchemaItem) | repeated | list of contructors json schema |
| events | [JsonSchemaItem](#proto-JsonSchemaItem) | repeated | list of events json schema |
| structs | [JsonSchemaItem](#proto-JsonSchemaItem) | repeated | list of declared structures json schema |






<a name="proto-JsonSchemaItem"></a>

### JsonSchemaItem
Json schema item


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name of json schema item |
| schema | [bytes](#bytes) |  | json schema |






<a name="proto-SubscribeRequest"></a>

### SubscribeRequest
List of requested subscriptions


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| head | [bool](#bool) |  |  |
| invokes | [InvokeFilters](#proto-InvokeFilters) | repeated |  |
| declares | [DeclareFilters](#proto-DeclareFilters) | repeated |  |
| deploys | [DeployFilters](#proto-DeployFilters) | repeated |  |
| deploy_accounts | [DeployAccountFilters](#proto-DeployAccountFilters) | repeated |  |
| l1_handlers | [L1HandlerFilter](#proto-L1HandlerFilter) | repeated |  |
| internals | [InternalFilter](#proto-InternalFilter) | repeated |  |
| fees | [FeeFilter](#proto-FeeFilter) | repeated |  |
| events | [EventFilter](#proto-EventFilter) | repeated |  |
| msgs | [MessageFilter](#proto-MessageFilter) | repeated |  |
| transfers | [TransferFilter](#proto-TransferFilter) | repeated |  |
| storage_diffs | [StorageDiffFilter](#proto-StorageDiffFilter) | repeated |  |
| token_balances | [TokenBalanceFilter](#proto-TokenBalanceFilter) | repeated |  |
| tokens | [TokenFilter](#proto-TokenFilter) | repeated |  |
| addresses | [AddressFilter](#proto-AddressFilter) | repeated |  |






<a name="proto-Subscription"></a>

### Subscription
Subscription entity. It contains subscription id and subscription&#39;s live notifications. It&#39;s response on `Subscribe` request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| response | [SubscribeResponse](#proto-SubscribeResponse) |  | message containing subscription id |
| block | [Block](#proto-Block) |  |  |
| declare | [Declare](#proto-Declare) |  |  |
| deploy | [Deploy](#proto-Deploy) |  |  |
| deploy_account | [DeployAccount](#proto-DeployAccount) |  |  |
| event | [Event](#proto-Event) |  |  |
| fee | [Fee](#proto-Fee) |  |  |
| internal | [Internal](#proto-Internal) |  |  |
| invoke | [Invoke](#proto-Invoke) |  |  |
| l1_handler | [L1Handler](#proto-L1Handler) |  |  |
| message | [StarknetMessage](#proto-StarknetMessage) |  |  |
| storage_diff | [StorageDiff](#proto-StorageDiff) |  |  |
| token_balance | [TokenBalance](#proto-TokenBalance) |  |  |
| transfer | [Transfer](#proto-Transfer) |  |  |
| token | [Token](#proto-Token) |  |  |
| address | [Address](#proto-Address) |  |  |
| end_of_block | [EndOfBlock](#proto-EndOfBlock) |  | message which signals about block data ends |





 

 

 


<a name="proto-IndexerService"></a>

### IndexerService
Desription of server interface

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Subscribe | [SubscribeRequest](#proto-SubscribeRequest) | [Subscription](#proto-Subscription) stream | Subscribe to live notification from indexer |
| Unsubscribe | [UnsubscribeRequest](#proto-UnsubscribeRequest) | [UnsubscribeResponse](#proto-UnsubscribeResponse) | Unsubscribe from live notification from indexer |
| JSONSchemaForClass | [Bytes](#proto-Bytes) | [Bytes](#proto-Bytes) | Receives JSON schema of class ABI by class hash |
| JSONSchemaForContract | [Bytes](#proto-Bytes) | [Bytes](#proto-Bytes) | Receives JSON schema of class ABI by contract hash |

 



<a name="github-com_dipdup-io_starknet-indexer_pkg_grpc_proto_response-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/dipdup-io/starknet-indexer/pkg/grpc/proto/response.proto



<a name="proto-Address"></a>

### Address



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| hash | [bytes](#bytes) |  |  |
| class_id | [uint64](#uint64) | optional |  |
| height | [uint64](#uint64) |  |  |






<a name="proto-Block"></a>

### Block



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| time | [uint64](#uint64) |  |  |
| version | [string](#string) |  |  |
| tx_count | [uint64](#uint64) |  |  |
| invokes_count | [uint64](#uint64) |  |  |
| declares_count | [uint64](#uint64) |  |  |
| deploys_count | [uint64](#uint64) |  |  |
| deploy_account_count | [uint64](#uint64) |  |  |
| l1_handlers_count | [uint64](#uint64) |  |  |
| storage_diffs_count | [uint64](#uint64) |  |  |
| status | [uint64](#uint64) |  |  |
| hash | [bytes](#bytes) |  |  |
| parent_hash | [bytes](#bytes) |  |  |
| new_root | [bytes](#bytes) |  |  |
| sequencer_address | [bytes](#bytes) |  |  |






<a name="proto-Declare"></a>

### Declare



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| time | [uint64](#uint64) |  |  |
| version | [uint64](#uint64) |  |  |
| position | [uint64](#uint64) |  |  |
| sender | [Address](#proto-Address) | optional |  |
| contract | [Address](#proto-Address) | optional |  |
| status | [uint64](#uint64) |  |  |
| class | [bytes](#bytes) |  |  |
| hash | [bytes](#bytes) |  |  |
| max_fee | [string](#string) |  |  |
| nonce | [string](#string) |  |  |






<a name="proto-Deploy"></a>

### Deploy



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| time | [uint64](#uint64) |  |  |
| position | [uint64](#uint64) |  |  |
| contract | [Address](#proto-Address) |  |  |
| status | [uint64](#uint64) |  |  |
| class | [bytes](#bytes) |  |  |
| hash | [bytes](#bytes) |  |  |
| salt | [bytes](#bytes) |  |  |
| calldata | [string](#string) | repeated |  |
| parsed_calldata | [bytes](#bytes) |  |  |






<a name="proto-DeployAccount"></a>

### DeployAccount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| time | [uint64](#uint64) |  |  |
| position | [uint64](#uint64) |  |  |
| contract | [Address](#proto-Address) |  |  |
| status | [uint64](#uint64) |  |  |
| class | [bytes](#bytes) |  |  |
| hash | [bytes](#bytes) |  |  |
| salt | [bytes](#bytes) |  |  |
| max_fee | [string](#string) |  |  |
| nonce | [string](#string) |  |  |
| calldata | [string](#string) | repeated |  |
| parsed_calldata | [bytes](#bytes) |  |  |






<a name="proto-EndOfBlock"></a>

### EndOfBlock



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| height | [uint64](#uint64) |  |  |






<a name="proto-Event"></a>

### Event



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| time | [uint64](#uint64) |  |  |
| order | [uint64](#uint64) |  |  |
| contract | [Address](#proto-Address) |  |  |
| from | [Address](#proto-Address) |  |  |
| keys | [string](#string) | repeated |  |
| data | [string](#string) | repeated |  |
| name | [string](#string) |  |  |
| parsed_data | [bytes](#bytes) |  |  |






<a name="proto-Fee"></a>

### Fee



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| time | [uint64](#uint64) |  |  |
| contract | [Address](#proto-Address) |  |  |
| caller | [Address](#proto-Address) |  |  |
| class | [bytes](#bytes) |  |  |
| selector | [bytes](#bytes) |  |  |
| entrypoint_type | [uint64](#uint64) |  |  |
| call_type | [uint64](#uint64) |  |  |
| calldata | [string](#string) | repeated |  |
| result | [string](#string) | repeated |  |
| entrypoint | [string](#string) |  |  |
| parsed_calldata | [bytes](#bytes) |  |  |






<a name="proto-Internal"></a>

### Internal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| time | [uint64](#uint64) |  |  |
| status | [uint64](#uint64) |  |  |
| hash | [bytes](#bytes) |  |  |
| contract | [Address](#proto-Address) |  |  |
| caller | [Address](#proto-Address) |  |  |
| class | [bytes](#bytes) |  |  |
| selector | [bytes](#bytes) |  |  |
| entrypoint_type | [uint64](#uint64) |  |  |
| call_type | [uint64](#uint64) |  |  |
| calldata | [string](#string) | repeated |  |
| result | [string](#string) | repeated |  |
| entrypoint | [string](#string) |  |  |
| parsed_calldata | [bytes](#bytes) |  |  |
| parsed_result | [bytes](#bytes) |  |  |






<a name="proto-Invoke"></a>

### Invoke



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| time | [uint64](#uint64) |  |  |
| status | [uint64](#uint64) |  |  |
| hash | [bytes](#bytes) |  |  |
| version | [uint64](#uint64) |  |  |
| position | [uint64](#uint64) |  |  |
| contract | [Address](#proto-Address) |  |  |
| selector | [bytes](#bytes) |  |  |
| max_fee | [string](#string) |  |  |
| nonce | [string](#string) |  |  |
| calldata | [string](#string) | repeated |  |
| entrypoint | [string](#string) |  |  |
| parsed_calldata | [bytes](#bytes) |  |  |






<a name="proto-L1Handler"></a>

### L1Handler



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| time | [uint64](#uint64) |  |  |
| status | [uint64](#uint64) |  |  |
| hash | [bytes](#bytes) |  |  |
| position | [uint64](#uint64) |  |  |
| contract | [Address](#proto-Address) |  |  |
| selector | [bytes](#bytes) |  |  |
| max_fee | [string](#string) |  |  |
| nonce | [string](#string) |  |  |
| calldata | [string](#string) | repeated |  |
| entrypoint | [string](#string) |  |  |
| parsed_calldata | [bytes](#bytes) |  |  |






<a name="proto-StarknetMessage"></a>

### StarknetMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| time | [uint64](#uint64) |  |  |
| contract | [Address](#proto-Address) |  |  |
| from | [Address](#proto-Address) |  |  |
| to | [Address](#proto-Address) |  |  |
| selector | [string](#string) |  |  |
| nonce | [string](#string) |  |  |
| payload | [string](#string) | repeated |  |






<a name="proto-StorageDiff"></a>

### StorageDiff



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| contract | [Address](#proto-Address) |  |  |
| key | [bytes](#bytes) |  |  |
| value | [bytes](#bytes) |  |  |






<a name="proto-Token"></a>

### Token



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| deploy_height | [uint64](#uint64) |  |  |
| deploy_time | [uint64](#uint64) |  |  |
| contract | [Address](#proto-Address) |  |  |
| owner | [Address](#proto-Address) |  |  |
| type | [int32](#int32) |  |  |
| metadata | [bytes](#bytes) |  |  |






<a name="proto-TokenBalance"></a>

### TokenBalance



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner | [Address](#proto-Address) |  |  |
| contract | [Address](#proto-Address) |  |  |
| token_id | [string](#string) |  |  |
| balance | [string](#string) |  |  |






<a name="proto-Transfer"></a>

### Transfer



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| height | [uint64](#uint64) |  |  |
| time | [uint64](#uint64) |  |  |
| contract | [Address](#proto-Address) |  |  |
| from | [Address](#proto-Address) |  |  |
| to | [Address](#proto-Address) |  |  |
| amount | [string](#string) |  |  |
| token_id | [string](#string) |  |  |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

