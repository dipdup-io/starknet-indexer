package txs

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
	"time"
)

type TransactionResult struct {
	TxType              string           `bun:"tx_type" json:"tx_type"`
	Hash                types.Hex        `bun:"hash" json:"hash"`
	BlockHeight         int64            `bun:"block_height" json:"block_height"`
	Timestamp           time.Time        `bun:"timestamp" json:"timestamp"`
	Status              int64            `bun:"status" json:"status"`
	Entrypoint          *string          `bun:"entrypoint" json:"entrypoint,omitempty"`
	Contract            types.Hex        `bun:"contract" json:"contract"`
	MaxFee              *string          `bun:"max_fee" json:"max_fee,omitempty"`
	ClassID             *int64           `bun:"class_id" json:"class_id,omitempty"`
	ContractAddressSalt *string          `bun:"contract_address_salt" json:"contract_address_salt,omitempty"`
	Calldata            *string          `bun:"calldata" json:"calldata,omitempty"`
	ParsedCalldata      *json.RawMessage `bun:"parsed_calldata" json:"parsed_calldata,omitempty"`
}

// GetTxByHash -
func GetTxByHash(ctx context.Context, s postgres.Storage, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hash, ok := request.Params.Arguments["hash"].(string)
	if !ok {
		return nil, errors.Errorf("hash must be a string")
	}

	txHashBytes, err := types.HexFromString(hash)
	if err != nil {
		return nil, err
	}
	query := `
		WITH hash_param AS (
			SELECT ?::bytea AS hash
		)
		
		SELECT
			'invoke' AS tx_type,
			i.hash,
			i.height AS block_height,
			i.time AS timestamp,
			i.status,
			i.entrypoint,
			a.hash AS contract,
			i.max_fee::text AS max_fee,
			NULL::bigint AS class_id,
			NULL::text AS contract_address_salt,
			array_to_string(i.call_data, ',') AS calldata,
			i.parsed_calldata
		FROM invoke i
		JOIN address a ON i.contract_id = a.id
		WHERE i.hash = (SELECT hash FROM hash_param)
		
		UNION ALL
		
		SELECT
			'deploy' AS tx_type,
			d.hash,
			d.height AS block_height,
			d.time AS timestamp,
			d.status,
			NULL AS entrypoint,
			a.hash AS contract,
			NULL::text AS max_fee,
			d.class_id,
			encode(d.contract_address_salt, 'hex') AS contract_address_salt,
			array_to_string(d.constructor_calldata, ',') AS calldata,
			NULL::jsonb AS parsed_calldata
		FROM deploy d
		JOIN address a ON d.contract_id = a.id
		WHERE d.hash = (SELECT hash FROM hash_param)
		
		UNION ALL
		
		SELECT
			'declare' AS tx_type,
			d.hash,
			d.height AS block_height,
			d.time AS timestamp,
			d.status,
			NULL AS entrypoint,
			a.hash AS contract,
			d.max_fee::text AS max_fee,
			d.class_id,
			NULL::text AS contract_address_salt,
			NULL::text AS calldata,
			NULL::jsonb AS parsed_calldata
		FROM declare d
		JOIN address a ON d.sender_id = a.id
		WHERE d.hash = (SELECT hash FROM hash_param)
		
		UNION ALL
		
		SELECT
			'deploy_account' AS tx_type,
			da.hash,
			da.height AS block_height,
			da.time AS timestamp,
			da.status,
			NULL AS entrypoint,
			a.hash AS contract,
			da.max_fee::text AS max_fee,
			da.class_id,
			encode(da.contract_address_salt, 'hex') AS contract_address_salt,
			array_to_string(da.constructor_calldata, ',') AS calldata,
			NULL::jsonb AS parsed_calldata
		FROM deploy_account da
		JOIN address a ON da.contract_id = a.id
		WHERE da.hash = (SELECT hash FROM hash_param)
		
		UNION ALL
		
		SELECT
			'l1_handler' AS tx_type,
			l.hash,
			l.height AS block_height,
			l.time AS timestamp,
			l.status,
			l.entrypoint,
			a.hash AS contract,
			l.max_fee::text AS max_fee,
			NULL AS class_id,
			NULL::text AS contract_address_salt,
			array_to_string(l.call_data, ',') AS calldata,
			l.parsed_calldata
		FROM l1_handler l
		JOIN address a ON l.contract_id = a.id
		WHERE l.hash = (SELECT hash FROM hash_param)`

	var result TransactionResult
	err = s.Connection().DB().NewRaw(query, txHashBytes).Scan(ctx, &result)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Errorf("transaction with hash %x not found", txHashBytes)
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to execute query, tx hash %x", txHashBytes)
	}
	jsonTx, err := json.Marshal(result)
	if err != nil {
		return nil, errors.Wrapf(err, "error marshalling tx json")
	}
	return mcp.NewToolResultText(string(jsonTx)), nil
}
