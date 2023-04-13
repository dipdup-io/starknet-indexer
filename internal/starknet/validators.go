package starknet

// HashValidator -
func HashValidator(hash []byte) bool {
	return len(hash) == 32
}
