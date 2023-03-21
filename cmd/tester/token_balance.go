package main

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/presets"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

// TokenBalanceTester -
type TokenBalanceTester struct {
	postgres postgres.Storage
	api      sequencer.API
}

// NewTokenBalanceTester -
func NewTokenBalanceTester(postgres postgres.Storage, api sequencer.API) TokenBalanceTester {
	return TokenBalanceTester{
		postgres: postgres,
		api:      api,
	}
}

// String -
func (tbt TokenBalanceTester) String() string {
	return "token balance tester"
}

// Test -
func (tbt TokenBalanceTester) Test(ctx context.Context) error {
	if err := tbt.testOnNegativeBalances(ctx); err != nil {
		return err
	}
	if err := tbt.testTotalSupplyForERC20(ctx); err != nil {
		return err
	}
	if err := tbt.testOwnerForERC271(ctx); err != nil {
		return err
	}

	return nil
}

// Close -
func (tbt TokenBalanceTester) Close() error {
	return nil
}

func (tbt TokenBalanceTester) testOnNegativeBalances(ctx context.Context) error {
	log.Info().Msg("test on negative balances...")

	nullAddress, err := tbt.postgres.Address.GetByHash(ctx, nullAddressHash)
	if err != nil {
		return err
	}

	exception1, err := tbt.postgres.Address.GetByHash(ctx, exceptNegativeTokenBalance1)
	if err != nil {
		return err
	}

	balances, err := tbt.postgres.TokenBalance.NegativeBalances(ctx)
	if err != nil {
		return err
	}

	if len(balances) == 0 {
		log.Info().Msg("negative balances are absent")
		return nil
	}

	for i := range balances {
		if balances[i].OwnerID == nullAddress.ID {
			continue
		}
		if balances[i].ContractID == exception1.ID {
			continue
		}

		log.Warn().
			Str("contract", encoding.EncodeHex(balances[i].Contract.Hash)).
			Str("owner", encoding.EncodeHex(balances[i].Owner.Hash)).
			Str("balance", balances[i].Balance.String()).
			Msg("negative balance")
	}

	return nil
}

func (tbt TokenBalanceTester) testTotalSupplyForERC20(ctx context.Context) error {
	log.Info().Msg("test total supply for ERC20...")

	var (
		offset = uint64(0)
		limit  = uint64(100)
		end    = false
	)

	last, err := tbt.postgres.Blocks.Last(ctx)
	if err != nil {
		return err
	}

	for !end {
		tokens, err := tbt.postgres.ERC20.List(ctx, limit, offset, storage.SortOrderAsc)
		if err != nil {
			return err
		}

		for i := range tokens {
			log.Info().Str("name", tokens[i].Name).Str("symbol", tokens[i].Symbol).Msg("test supply of ERC20")

			dbTotalSupply, err := tbt.postgres.TokenBalance.TotalSupply(ctx, tokens[i].ContractID, 0)
			if err != nil {
				return err
			}

			contract, err := tbt.postgres.Address.GetByID(ctx, tokens[i].ContractID)
			if err != nil {
				return err
			}

			token := presets.NewERC20(tbt.api, data.NewFeltFromBytes(contract.Hash))
			apiTotalSupply, err := token.TotalSupply(ctx, presets.WithBlockID(data.BlockID{
				Number: &last.Height,
			}))
			if err != nil {
				return err
			}

			if !dbTotalSupply.Equal(apiTotalSupply) {
				multiplier := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(tokens[i].Decimals)))
				diff := dbTotalSupply.Sub(apiTotalSupply).Div(multiplier)
				log.Info().
					Str("name", tokens[i].Name).
					Str("symbol", tokens[i].Symbol).
					Str("db_total_supply", dbTotalSupply.Div(multiplier).String()).
					Str("api_total_supply", apiTotalSupply.Div(multiplier).String()).
					Str("diff", diff.String()).
					Msg("total supplies are differ")
			}
		}

		end = len(tokens) < int(limit)
		offset += uint64(len(tokens))
	}

	return nil
}

func (tbt TokenBalanceTester) testOwnerForERC271(ctx context.Context) error {
	log.Info().Msg("test total supply for ERC20...")

	var (
		offset = uint64(0)
		limit  = uint64(100)
		end    = false
	)

	last, err := tbt.postgres.Blocks.Last(ctx)
	if err != nil {
		return err
	}

	for !end {
		tokens, err := tbt.postgres.ERC721.List(ctx, limit, offset, storage.SortOrderDesc)
		if err != nil {
			return err
		}

		for i := range tokens {
			log.Info().Str("name", tokens[i].Name).Str("symbol", tokens[i].Symbol).Msg("test owner of ERC721")

			nfts, err := tbt.postgres.TokenBalance.List(ctx, 5, 0, storage.SortOrderDesc)
			if err != nil {
				return err
			}

			for _, nft := range nfts {
				tokenBalance, err := tbt.postgres.TokenBalance.Owner(ctx, tokens[i].ContractID, nft.TokenID)
				if err != nil {
					return err
				}

				contract, err := tbt.postgres.Address.GetByID(ctx, tokens[i].ContractID)
				if err != nil {
					return err
				}

				token := presets.NewERC721(tbt.api, data.NewFeltFromBytes(contract.Hash))

				tokenId, err := data.NewUint256FromString(nft.TokenID.String())
				if err != nil {
					return err
				}

				apiOwner, err := token.OwnerOf(ctx, tokenId, presets.WithBlockID(data.BlockID{
					Number: &last.Height,
				}))
				if err != nil {
					return err
				}

				dbOwner := encoding.EncodeHex(tokenBalance.Owner.Hash)

				if dbOwner != apiOwner.String() {
					log.Info().
						Str("name", tokens[i].Name).
						Str("symbol", tokens[i].Symbol).
						Str("db_owner", dbOwner).
						Str("api_owner", apiOwner.String()).
						Msg("total supplies are differ")
				}
			}
		}

		end = len(tokens) < int(limit)
		offset += uint64(len(tokens))
	}
	return nil
}
