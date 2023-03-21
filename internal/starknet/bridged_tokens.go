package starknet

import (
	"os"
	"sync"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
)

// BridgedToken -
type BridgedToken struct {
	Name            string    `json:"name"`
	Symbol          string    `json:"symbol"`
	Decimals        uint64    `json:"decimals"`
	L1TokenAddress  string    `json:"l1_token_address"`
	L2TokenAddress  data.Felt `json:"l2_token_address"`
	L1BridgeAddress string    `json:"l1_bridge_address"`
	L2BridgeAddress data.Felt `json:"l2_bridge_address"`
}

var bridgedTokens = make([]BridgedToken, 0)
var loadOnce sync.Once

// LoadBridgedTokens -
func LoadBridgedTokens(filename string) (retErr error) {
	loadOnce.Do(func() {
		if filename == "" {
			return
		}
		f, err := os.Open(filename)
		if err != nil {
			retErr = err
			return
		}
		defer f.Close()

		retErr = json.NewDecoder(f).Decode(&bridgedTokens)
	})
	return retErr
}

// BridgedTokens -
func BridgedTokens() []BridgedToken {
	return bridgedTokens
}
