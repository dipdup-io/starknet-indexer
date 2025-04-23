package address

import (
	"context"
	"encoding/json"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type TransferInfo struct {
	TxHash          string    `json:"tx_hash"`
	Height          int64     `json:"block_height"`
	Time            time.Time `json:"timestamp"`
	FromAddress     string    `json:"from_address"`
	ToAddress       string    `json:"to_address"`
	ContractAddress string    `json:"contract_address"`
	TokenID         string    `json:"token_id,omitempty"`
	Amount          string    `json:"amount"`
	Type            string    `json:"tx_type"`
}

type InvocationInfo struct {
	TxHash          string    `json:"tx_hash"`
	Height          uint64    `json:"block_height"`
	Time            time.Time `json:"timestamp"`
	ContractAddress string    `json:"contract_address"`
	Entrypoint      string    `json:"entrypoint"`
	Status          string    `json:"status"`
	Error           string    `json:"error,omitempty"`
}

type EventInfo struct {
	TxHash          string    `json:"tx_hash"`
	Height          uint64    `json:"block_height"`
	Time            time.Time `json:"timestamp"`
	ContractAddress string    `json:"contract_address"`
}

type DeployInfo struct {
	TxHash          string    `json:"tx_hash"`
	Height          int64     `json:"block_height"`
	Time            time.Time `json:"timestamp"`
	ContractAddress string    `json:"contract_address"`
	ClassHash       string    `json:"class_hash"`
}

type Activity struct {
	Address           string                    `json:"address"`
	OutgoingTransfers []*TransferInfo           `json:"outgoing_transfers"`
	IncomingTransfers []*TransferInfo           `json:"incoming_transfers"`
	Invocations       []*InvocationInfo         `json:"invocations"`
	Events            []*EventInfo              `json:"events"`
	Deploys           []models.DeployedContract `json:"deploys"`
	TotalCountInfo    map[string]int            `json:"total_count_info"`
}

const LimitMaxValue = 100

// GetAddressActivity - returns complex chain activity address data and stats
func GetAddressActivity(storage postgres.Storage, ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var err error
	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return nil, errors.Errorf("address must be a string")
	}

	limit := 20
	if limitParam, ok := request.Params.Arguments["limit"].(string); ok && limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return nil, errors.Wrapf(err, "error casting limit to int")
		}
	}
	if limit > LimitMaxValue {
		return nil, errors.Errorf("limit max value  is %d", LimitMaxValue)
	}

	offset := 0
	if offsetParam, ok := request.Params.Arguments["offset"].(string); ok && offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			return nil, errors.Wrapf(err, "error casting offset to int")
		}
	}

	addressHash, err := models.HexToBytes(address)
	if err != nil {
		return nil, errors.Wrapf(err, "error converting address to bytes: %s", address)
	}

	activity := Activity{
		Address:        address,
		TotalCountInfo: make(map[string]int),
	}

	outTransfers, err := storage.Transfer.Filter(ctx,
		[]models.TransferFilter{
			{
				From: models.BytesFilter{
					Eq: addressHash,
				},
			},
		},
		models.WithLimitFilter(limit),
		models.WithOffsetFilter(offset),
		models.WithDescSortByIdFilter(),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error filtering outgoing transfers")
	}

	totalOutTransfersCount, err := storage.Transfer.Count(ctx,
		[]models.TransferFilter{
			{
				From: models.BytesFilter{
					Eq: addressHash,
				},
			},
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error counting outgoing transfers")
	}
	activity.TotalCountInfo["outgoing_transfers"] = int(totalOutTransfersCount)
	parsedTransfers, err := parseTransfers(ctx, storage, outTransfers)
	if err != nil {
		return nil, err
	}
	activity.OutgoingTransfers = parsedTransfers

	inTransfers, err := storage.Transfer.Filter(ctx,
		[]models.TransferFilter{
			{
				To: models.BytesFilter{
					Eq: addressHash,
				},
			},
		},
		models.WithLimitFilter(limit),
		models.WithOffsetFilter(offset),
		models.WithDescSortByIdFilter(),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error filtering incoming transfers")
	}

	totalInTransfersCount, err := storage.Transfer.Count(ctx,
		[]models.TransferFilter{
			{
				To: models.BytesFilter{
					Eq: addressHash,
				},
			},
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error counting incoming transfers")
	}
	activity.TotalCountInfo["incoming_transfers"] = int(totalInTransfersCount)
	parsedTransfers, err = parseTransfers(ctx, storage, inTransfers)
	if err != nil {
		return nil, err
	}
	activity.IncomingTransfers = parsedTransfers

	invocations, err := storage.Invoke.Filter(ctx,
		[]models.InvokeFilter{
			{
				Contract: models.BytesFilter{
					Eq: addressHash,
				},
			},
		},
		models.WithLimitFilter(limit),
		models.WithOffsetFilter(offset),
		models.WithDescSortByIdFilter(),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error filtering invocations")
	}

	invocationsCount, err := storage.Invoke.Count(ctx,
		[]models.InvokeFilter{
			{
				Contract: models.BytesFilter{
					Eq: addressHash,
				},
			},
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error counting invocations")
	}

	activity.TotalCountInfo["invocations"] = int(invocationsCount)

	parsedInvocations := make([]*InvocationInfo, len(invocations))
	for i := range invocations {
		storageContractAddress, err := storage.Address.GetByID(ctx, invocations[i].ContractID)
		if err != nil {
			return nil, errors.Wrapf(err, "can't get address by id %d", invocations[i].ContractID)
		}
		contractAddress := models.BytesToFormattedHex(storageContractAddress.Hash)

		txHash := ""
		if invocations[i].Hash != nil {
			txHash = models.BytesToFormattedHex(invocations[i].Hash)
		}

		parsedInvocations[i] = &InvocationInfo{
			TxHash:          txHash,
			Height:          invocations[i].Height,
			Time:            invocations[i].Time,
			ContractAddress: contractAddress,
			Entrypoint:      invocations[i].Entrypoint,
			Status:          statusToString(invocations[i].Status),
		}
		if invocations[i].Error != nil {
			parsedInvocations[i].Error = *invocations[i].Error
		}
	}
	activity.Invocations = parsedInvocations

	storageAddress, err := storage.Address.GetByHash(ctx, addressHash)
	if err != nil {
		return nil, errors.Wrapf(err, "can't get address %s", address)
	}

	events, err := storage.Event.Filter(ctx,
		[]models.EventFilter{
			{
				Contract: models.IdFilter{
					Eq: storageAddress.ID,
				},
			},
		},
		models.WithLimitFilter(limit),
		models.WithOffsetFilter(offset),
		models.WithDescSortByIdFilter(),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error filtering events with address %s", address)
	}

	eventsCount, err := storage.Event.Count(ctx,
		[]models.EventFilter{
			{
				Contract: models.IdFilter{
					Eq: storageAddress.ID,
				},
			},
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error counting events with address %s", address)
	}
	activity.TotalCountInfo["events"] = int(eventsCount)

	eventsList := make([]*EventInfo, len(events))
	for i, event := range events {
		storageContractAddress, err := storage.Address.GetByID(ctx, event.ContractID)
		if err != nil {
			return nil, errors.Wrapf(err, "can't get address by id %d", event.ContractID)
		}
		contractAddress := models.BytesToFormattedHex(storageContractAddress.Hash)

		txHash := ""
		switch {
		case events[i].InvokeID != nil:
			tx, err := storage.Invoke.GetByID(ctx, *events[i].InvokeID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching invoke with id %d", events[i].InvokeID)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case events[i].DeclareID != nil:
			tx, err := storage.Declare.GetByID(ctx, *events[i].DeclareID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching declare with id %d", events[i].DeclareID)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case events[i].DeployID != nil:
			tx, err := storage.Deploy.GetByID(ctx, *events[i].DeployID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching deploy with id %d", events[i].DeployID)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case events[i].DeployAccountID != nil:
			tx, err := storage.DeployAccount.GetByID(ctx, *events[i].DeployAccountID)
			if err != nil {
				return nil, errors.Wrapf(
					err,
					"error fetching deploy_account with id %d",
					events[i].DeployAccountID,
				)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case events[i].L1HandlerID != nil:
			tx, err := storage.Deploy.GetByID(ctx, *events[i].L1HandlerID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching l1_handler with id %d", events[i].L1HandlerID)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case events[i].InternalID != nil:
			tx, err := storage.Internal.GetByID(ctx, *events[i].InternalID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching internal tx with id %d", events[i].InternalID)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case events[i].FeeID != nil:
			fee, err := storage.Fee.GetByID(ctx, *events[i].FeeID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching fee with id %d", events[i].FeeID)
			}
			txHash, err = getTxHashFromFee(ctx, &storage, fee.ID)
			if err != nil {
				return nil, err
			}
		}

		eventsList[i] = &EventInfo{
			TxHash:          txHash,
			Height:          event.Height,
			Time:            event.Time,
			ContractAddress: contractAddress,
		}
	}
	activity.Events = eventsList

	deploys, err := storage.Internal.GetDeployedContracts(ctx, addressHash)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting deployed contracts")
	}
	activity.Deploys = deploys

	jsonResult, err := json.Marshal(activity)
	if err != nil {
		return nil, errors.Wrapf(err, "error marshalling address activity data")
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

func parseTransfers(ctx context.Context, storage postgres.Storage, transfers []models.Transfer) ([]*TransferInfo, error) {
	resultTransfers := make([]*TransferInfo, len(transfers))
	for i := range transfers {
		contractAddress := ""
		if transfers[i].Contract.Hash != nil {
			contractAddress = models.BytesToFormattedHex(transfers[i].Contract.Hash)
		}

		fromAddress := ""
		if transfers[i].From.Hash != nil {
			fromAddress = models.BytesToFormattedHex(transfers[i].From.Hash)
		}

		toAddress := ""
		if transfers[i].To.Hash != nil {
			toAddress = models.BytesToFormattedHex(transfers[i].To.Hash)
		}

		txHash := ""
		txType := "unknown"
		switch {
		case transfers[i].InvokeID != nil:
			txType = "invoke"
			tx, err := storage.Invoke.GetByID(ctx, *transfers[i].InvokeID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching invoke with id %d", transfers[i].InvokeID)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case transfers[i].DeclareID != nil:
			txType = "declare"
			tx, err := storage.Declare.GetByID(ctx, *transfers[i].DeclareID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching declare with id %d", transfers[i].DeclareID)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case transfers[i].DeployID != nil:
			txType = "deploy"
			tx, err := storage.Deploy.GetByID(ctx, *transfers[i].DeployID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching deploy with id %d", transfers[i].DeployID)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case transfers[i].DeployAccountID != nil:
			txType = "deploy_account"
			tx, err := storage.DeployAccount.GetByID(ctx, *transfers[i].DeployAccountID)
			if err != nil {
				return nil, errors.Wrapf(
					err,
					"error fetching deploy_account with id %d",
					transfers[i].DeployAccountID,
				)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case transfers[i].L1HandlerID != nil:
			txType = "l1_handler"
			tx, err := storage.Deploy.GetByID(ctx, *transfers[i].L1HandlerID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching l1_handler with id %d", transfers[i].L1HandlerID)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case transfers[i].InternalID != nil:
			txType = "internal"
			tx, err := storage.Internal.GetByID(ctx, *transfers[i].FeeID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching internal tx with id %d", transfers[i].InternalID)
			}
			txHash = models.BytesToFormattedHex(tx.Hash)
		case transfers[i].FeeID != nil:
			txType = "fee"
			fee, err := storage.Fee.GetByID(ctx, *transfers[i].FeeID)
			if err != nil {
				return nil, errors.Wrapf(err, "error fetching fee with id %d", transfers[i].FeeID)
			}
			txHash, err = getTxHashFromFee(ctx, &storage, fee.ID)
			if err != nil {
				return nil, err
			}
		}

		resultTransfers[i] = &TransferInfo{
			TxHash:          txHash,
			Height:          int64(transfers[i].Height),
			Time:            transfers[i].Time,
			FromAddress:     fromAddress,
			ToAddress:       toAddress,
			ContractAddress: contractAddress,
			TokenID:         transfers[i].TokenID.String(),
			Amount:          transfers[i].Amount.String(),
			Type:            txType,
		}
	}
	return resultTransfers, nil
}

func getTxHashFromFee(ctx context.Context, storage *postgres.Storage, feeID uint64) (string, error) {
	fee, err := storage.Fee.GetByID(ctx, feeID)
	if err != nil {
		return "", errors.Wrapf(err, "error fetching fee with id %d", feeID)
	}

	var txHash string
	var txErr error

	switch {
	case fee.InvokeID != nil:
		invoke, err := storage.Invoke.GetByID(ctx, *fee.InvokeID)
		if err != nil {
			txErr = errors.Wrapf(err, "error fetching invoke with id %d", *fee.InvokeID)
		} else {
			txHash = models.BytesToFormattedHex(invoke.Hash)
		}
	case fee.DeclareID != nil:
		declare, err := storage.Declare.GetByID(ctx, *fee.DeclareID)
		if err != nil {
			txErr = errors.Wrapf(err, "error fetching declare with id %d", *fee.DeclareID)
		} else {
			txHash = models.BytesToFormattedHex(declare.Hash)
		}
	case fee.DeployID != nil:
		deploy, err := storage.Deploy.GetByID(ctx, *fee.DeployID)
		if err != nil {
			txErr = errors.Wrapf(err, "error fetching deploy with id %d", *fee.DeployID)
		} else {
			txHash = models.BytesToFormattedHex(deploy.Hash)
		}
	case fee.DeployAccountID != nil:
		deployAccount, err := storage.DeployAccount.GetByID(ctx, *fee.DeployAccountID)
		if err != nil {
			txErr = errors.Wrapf(err, "error fetching deploy account with id %d", *fee.DeployAccountID)
		} else {
			txHash = models.BytesToFormattedHex(deployAccount.Hash)
		}
	case fee.L1HandlerID != nil:
		l1Handler, err := storage.L1Handler.GetByID(ctx, *fee.L1HandlerID)
		if err != nil {
			txErr = errors.Wrapf(err, "error fetching l1 handler with id %d", *fee.L1HandlerID)
		} else {
			txHash = models.BytesToFormattedHex(l1Handler.Hash)
		}
	default:
		return "", errors.Errorf("fee does not have any associated transaction ID")
	}

	if txErr != nil {
		return "", txErr
	}

	if txHash == "" {
		return "", errors.Errorf("transaction hash not found")
	}

	return txHash, nil
}

func statusToString(status models.Status) string {
	switch status {
	case models.StatusUnknown:
		return "UNKNOWN"
	case models.StatusNotReceived:
		return "NOT_RECEIVED"
	case models.StatusReceived:
		return "RECEIVED"
	case models.StatusPending:
		return "PENDING"
	case models.StatusRejected:
		return "REJECTED"
	case models.StatusAcceptedOnL2:
		return "ACCEPTED_ON_L2"
	case models.StatusAcceptedOnL1:
		return "ACCEPTED_ON_L1"
	case models.StatusReverted:
		return "REVERTED"
	default:
		return "INVALID_STATUS"
	}
}
