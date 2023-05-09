package v0

import (
	"context"

	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/resolver"
	"github.com/pkg/errors"
)

var (
	errInvalidTokenAddress = errors.New("invalid token address")
)

// TokenParser -
type TokenParser struct {
	resolver resolver.Resolver
	cache    *cache.Cache
}

// NewTokenParser -
func NewTokenParser(cache *cache.Cache, resolver resolver.Resolver) TokenParser {
	return TokenParser{
		cache:    cache,
		resolver: resolver,
	}
}

// Parse -
func (parser TokenParser) Parse(ctx context.Context, txCtx data.TxContext, contract storage.Address, classType storage.ClassType, constructorData map[string]any) (*storage.Token, error) {
	switch classType {
	case storage.ClassTypeERC20:
		return parser.getErc20(ctx, txCtx, contract.ID, constructorData)
	case storage.ClassTypeERC721:
		return parser.getErc721(ctx, txCtx, contract.ID, constructorData)
	case storage.ClassTypeERC1155:
		return parser.getErc1155(ctx, txCtx, contract.ID, constructorData)
	}

	return nil, nil
}

func (parser TokenParser) getErc20(ctx context.Context, txCtx data.TxContext, contractId uint64, constructorData map[string]any) (*storage.Token, error) {
	token := storage.Token{
		DeployHeight: txCtx.Height,
		DeployTime:   txCtx.Time,
		ContractID:   contractId,
		Type:         storage.TokenTypeERC20,
		Metadata:     map[string]any{},
	}

	token.Metadata["name"] = getStringFromParsedData(constructorData, "name")
	token.Metadata["symbol"] = getStringFromParsedData(constructorData, "symbol")
	decimals, err := getUint64FromParsedData(constructorData, "decimals")
	if err != nil {
		return nil, err
	}
	token.Metadata["decimals"] = decimals
	return &token, nil
}

func (parser TokenParser) getErc721(ctx context.Context, txCtx data.TxContext, contractId uint64, constructorData map[string]any) (*storage.Token, error) {
	token := storage.Token{
		DeployHeight: txCtx.Height,
		DeployTime:   txCtx.Time,
		ContractID:   contractId,
		Type:         storage.TokenTypeERC721,
		Metadata:     map[string]any{},
	}
	var err error

	token.Metadata["name"] = getStringFromParsedData(constructorData, "name")
	token.Metadata["symbol"] = getStringFromParsedData(constructorData, "symbol")

	token.OwnerID, err = getTokenAddress(ctx, parser.resolver, constructorData, token.DeployHeight, "owner")
	if errors.Is(err, errInvalidTokenAddress) {
		return &token, nil
	}
	return &token, err
}

func (parser TokenParser) getErc1155(ctx context.Context, txCtx data.TxContext, contractId uint64, constructorData map[string]any) (*storage.Token, error) {
	token := storage.Token{
		DeployHeight: txCtx.Height,
		DeployTime:   txCtx.Time,
		ContractID:   contractId,
		Type:         storage.TokenTypeERC1155,
		Metadata:     map[string]any{},
	}
	var err error

	token.Metadata["uri"] = getStringFromParsedData(constructorData, "uri")
	token.OwnerID, err = getTokenAddress(ctx, parser.resolver, constructorData, token.DeployHeight, "owner")
	if errors.Is(err, errInvalidTokenAddress) {
		return &token, nil
	}
	return &token, err
}

func getStringFromParsedData(constructorData map[string]any, key string) string {
	if value, ok := constructorData[key]; ok {
		if str, ok := value.(string); ok {
			felt := starknetData.Felt(str)
			return felt.ToAsciiString()
		}
	}
	return ""
}

func getUint64FromParsedData(constructorData map[string]any, key string) (uint64, error) {
	if value, ok := constructorData[key]; ok {
		if str, ok := value.(string); ok {
			felt := starknetData.Felt(str)
			return felt.Uint64()
		}
	}
	return 0, nil
}

func getTokenAddress(ctx context.Context, resolver resolver.Resolver, constructorData map[string]any, height uint64, key string) (uint64, error) {
	if value, ok := constructorData[key]; ok {
		if sValue, ok := value.(string); ok {
			address := storage.Address{
				Hash:   encoding.MustDecodeHex(sValue),
				Height: height,
			}
			if err := resolver.FindAddress(ctx, &address); err != nil {
				return 0, err
			}
			return address.ID, nil
		}
	}
	return 0, errInvalidTokenAddress
}
