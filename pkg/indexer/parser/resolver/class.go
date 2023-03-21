package resolver

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

// ReceiveClass -
func (resolver *Resolver) ReceiveClass(ctx context.Context, class *storage.Class) error {
	rawClass, err := resolver.receiver.GetClass(ctx, encoding.EncodeHex(class.Hash))
	if err != nil {
		return err
	}

	class.Abi = storage.Bytes(rawClass.RawAbi)

	a, err := rawClass.GetAbi()
	if err != nil {
		return err
	}
	interfaces, err := starknet.FindInterfaces(a)
	if err != nil {
		return err
	}
	class.Type = storage.NewClassType(interfaces...)

	resolver.cache.SetAbiByClassHash(*class, a)
	resolver.addClass(class)

	return nil
}

// FindClass -
func (resolver *Resolver) FindClass(ctx context.Context, class *storage.Class) error {
	if value, ok := resolver.classes[encoding.EncodeHex(class.Hash)]; ok {
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
func (resolver *Resolver) FindClassByHash(ctx context.Context, hash data.Felt) (*storage.Class, error) {
	if hash == "" {
		return nil, nil
	}
	class := &storage.Class{
		Hash: hash.Bytes(),
	}
	err := resolver.FindClass(ctx, class)
	return class, err
}

func (resolver *Resolver) addClass(class *storage.Class) {
	if len(class.Hash) == 0 {
		return
	}
	key := encoding.EncodeHex(class.Hash)
	if _, ok := resolver.classes[key]; !ok {
		resolver.classes[key] = class
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
			return class, err
		}
	}
	return class, nil
}

// FindClassByID -
func (resolver *Resolver) FindClassByID(ctx context.Context, id uint64) (*storage.Class, error) {
	for _, class := range resolver.classes {
		if class.ID == id {
			return class, nil
		}
	}
	return resolver.cache.GetClassById(ctx, id)
}
