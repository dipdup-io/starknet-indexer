package tools

import (
	"context"
	"encoding/json"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
)

func GetTxByHash(storage postgres.Storage, ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hash, ok := request.Params.Arguments["hash"].(string)
	if !ok {
		return nil, errors.Errorf("hash must be a string")
	}

	hashBytes, err := models.HexToBytes(hash)
	if err != nil {
		return nil, err
	}

	// invokes
	invokes, err := storage.Invoke.Filter(ctx,
		[]models.InvokeFilter{
			{
				Hash: models.BytesFilter{
					Eq: hashBytes,
				},
			},
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error during executing filter on invoke")
	}
	if len(invokes) == 1 {
		jsonInvokes, err := json.MarshalIndent(invokes[0], "", "  ")
		if err != nil {
			return nil, errors.Wrapf(err, "error marshalling invoke")
		}
		return mcp.NewToolResultText(string(jsonInvokes)), nil
	}

	// deploys
	deploys, err := storage.Deploy.Filter(ctx,
		[]models.DeployFilter{
			{
				Hash: models.BytesFilter{
					Eq: hashBytes,
				},
			},
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error during executing filter on deploy")
	}
	if len(deploys) == 1 {
		jsonDeploys, err := json.MarshalIndent(deploys[0], "", "  ")
		if err != nil {
			return nil, errors.Wrapf(err, "error marshalling deploy")
		}
		return mcp.NewToolResultText(string(jsonDeploys)), nil
	}

	// declares
	declares, err := storage.Declare.Filter(ctx,
		[]models.DeclareFilter{
			{
				Hash: models.BytesFilter{
					Eq: hashBytes,
				},
			},
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error during executing filter on declare")
	}
	if len(declares) == 1 {
		jsonDeclares, err := json.MarshalIndent(declares[0], "", "  ")
		if err != nil {
			return nil, errors.Wrapf(err, "error marshalling declare")
		}
		return mcp.NewToolResultText(string(jsonDeclares)), nil
	}

	// deployAccount
	accountDeploys, err := storage.DeployAccount.Filter(ctx,
		[]models.DeployAccountFilter{
			{
				Hash: models.BytesFilter{
					Eq: hashBytes,
				},
			},
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error during executing filter on deploy account")
	}
	if len(accountDeploys) == 1 {
		jsonAccountDeploys, err := json.MarshalIndent(accountDeploys[0], "", "  ")
		if err != nil {
			return nil, errors.Wrapf(err, "error marshalling account deploy")
		}
		return mcp.NewToolResultText(string(jsonAccountDeploys)), nil
	}

	// l1_handler
	l1Handlers, err := storage.L1Handler.Filter(ctx,
		[]models.L1HandlerFilter{
			{
				Hash: models.BytesFilter{
					Eq: hashBytes,
				},
			},
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error during executing filter on l1_handler")
	}
	if len(l1Handlers) == 1 {
		jsonL1Handlers, err := json.MarshalIndent(l1Handlers[0], "", "  ")
		if err != nil {
			return nil, errors.Wrapf(err, "error marshalling la_handler")
		}
		return mcp.NewToolResultText(string(jsonL1Handlers)), nil
	}

	return mcp.NewToolResultText(""), nil
}
