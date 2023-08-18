package resolver

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/rs/zerolog/log"
)

// ReceiveClass -
func (resolver *Resolver) ReceiveClass(ctx context.Context, class *storage.Class) error {
	rawClass, err := resolver.receiver.GetClass(ctx, encoding.EncodeHex(class.Hash))
	if err != nil {
		return err
	}

	if rawClass.RawAbi != nil {
		log.Debug().Hex("hash", class.Hash).Msg("class received")
		a, err := rawClass.GetAbi()
		if err != nil {
			return err
		}
		class.Abi = storage.Bytes(rawClass.RawAbi)
		interfaces, err := starknet.FindInterfaces(a)
		if err != nil {
			return err
		}
		class.Type = storage.NewClassType(interfaces...)
		resolver.cache.SetAbiByClassHash(*class, a)
		resolver.cache.SetClassByHash(*class)
	}

	switch rawClass.ClassVersion {
	case "0.1.0":
		class.Cairo = 1
	default:
		class.Cairo = 0
	}

	resolver.addClass(class)

	return nil
}

// FindClass -
func (resolver *Resolver) FindClass(ctx context.Context, class *storage.Class) error {
	classes := resolver.blockContext.Classes()
	if value, ok := classes[encoding.EncodeHex(class.Hash)]; ok {
		class.ID = value.ID
		class.Abi = value.Abi
		class.Type = value.Type
		return nil
	}
	generated, err := resolver.idGenerator.SetClassId(ctx, class)
	if err != nil {
		return err
	}
	if generated {
		resolver.addClass(class)
	}
	return nil
}

// FindClassByHash -
func (resolver *Resolver) FindClassByHash(ctx context.Context, hash data.Felt, height uint64) (*storage.Class, error) {
	if hash == "" {
		return nil, nil
	}
	class := &storage.Class{
		Hash:   hash.Bytes(),
		Height: height,
	}
	err := resolver.FindClass(ctx, class)
	return class, err
}

func (resolver *Resolver) addClass(class *storage.Class) {
	if len(class.Hash) == 0 {
		return
	}
	key := encoding.EncodeHex(class.Hash)
	classes := resolver.blockContext.Classes()
	if _, ok := classes[key]; !ok {
		classes[key] = class
	}
}

func (resolver *Resolver) parseClassFromFelt(ctx context.Context, classHash data.Felt, height uint64) (storage.Class, error) {
	return resolver.parseClass(ctx, classHash.Bytes(), height)
}

func (resolver *Resolver) parseClass(ctx context.Context, classHash []byte, height uint64) (storage.Class, error) {
	class := storage.Class{
		Hash:   classHash,
		Height: height,
	}
	if err := resolver.FindClass(ctx, &class); err != nil {
		return class, err
	}
	if class.Abi == nil {
		if err := resolver.ReceiveClass(ctx, &class); err != nil {
			log.Err(err).Msg("receive class error")
			classes := resolver.blockContext.Classes()
			delete(classes, encoding.EncodeHex(class.Hash))
			return class, err
		}
	}
	return class, nil
}

// FindClassByID -
func (resolver *Resolver) FindClassByID(ctx context.Context, id uint64) (*storage.Class, error) {
	classes := resolver.blockContext.Classes()
	for _, class := range classes {
		if class.ID == id {
			return class, nil
		}
	}
	return resolver.cache.GetClassById(ctx, id)
}
