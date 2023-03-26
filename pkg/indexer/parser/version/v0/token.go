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
func (parser TokenParser) Parse(ctx context.Context, txCtx data.TxContext, contract storage.Address, classType storage.ClassType, constructorData map[string]any) (data.Token, error) {
	var token data.Token

	switch classType {
	case storage.ClassTypeERC20:
		erc20, err := parser.getErc20(ctx, txCtx, contract.ID, constructorData)
		if err != nil {
			return token, err
		}
		token.ERC20 = erc20
	case storage.ClassTypeERC721:
		erc721, err := parser.getErc721(ctx, txCtx, contract.ID, constructorData)
		if err != nil {
			return token, err
		}
		token.ERC721 = erc721
	case storage.ClassTypeERC1155:
		erc1155, err := parser.getErc1155(ctx, txCtx, contract.ID, constructorData)
		if err != nil {
			return token, err
		}
		token.ERC1155 = erc1155
	}

	return token, nil
}

func (parser TokenParser) getErc20(ctx context.Context, txCtx data.TxContext, contractId uint64, constructorData map[string]any) (*storage.ERC20, error) {
	token := storage.ERC20{
		DeployHeight: txCtx.Height,
		DeployTime:   txCtx.Time,
		ContractID:   contractId,
	}
	var err error

	token.Name = getStringFromParsedData(constructorData, "name")
	token.Symbol = getStringFromParsedData(constructorData, "symbol")
	token.Decimals, err = getUint64FromParsedData(constructorData, "decimals")
	return &token, err
}

func (parser TokenParser) getErc721(ctx context.Context, txCtx data.TxContext, contractId uint64, constructorData map[string]any) (*storage.ERC721, error) {
	token := storage.ERC721{
		DeployHeight: txCtx.Height,
		DeployTime:   txCtx.Time,
		ContractID:   contractId,
	}
	var err error

	token.Name = getStringFromParsedData(constructorData, "name")
	token.Symbol = getStringFromParsedData(constructorData, "symbol")

	token.OwnerID, err = getTokenAddress(ctx, parser.resolver, constructorData, token.DeployHeight, "owner")
	if errors.Is(err, errInvalidTokenAddress) {
		return &token, nil
	}
	return &token, err
}

func (parser TokenParser) getErc1155(ctx context.Context, txCtx data.TxContext, contractId uint64, constructorData map[string]any) (*storage.ERC1155, error) {
	token := storage.ERC1155{
		DeployHeight: txCtx.Height,
		DeployTime:   txCtx.Time,
		ContractID:   contractId,
	}
	var err error

	token.TokenUri = getStringFromParsedData(constructorData, "uri")
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
