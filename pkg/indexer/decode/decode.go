package decode

import (
	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/pkg/errors"
)

// CalldataBySelector -
func CalldataBySelector(contractAbi abi.Abi, selector []byte, calldata []string) (map[string]any, string, error) {
	function, ok := contractAbi.GetFunctionBySelector(encoding.EncodeHex(selector))
	if !ok {
		return nil, "", errors.Errorf("unknown selector: %x", selector)
	}

	if len(calldata) == 0 {
		return nil, function.Name, nil
	}

	parsed, err := abi.DecodeFunctionCallData(calldata, *function, contractAbi.Structs)
	return parsed, function.Name, err
}

// CalldataForConstructor -
func CalldataForConstructor(classAbi abi.Abi, calldata []string) (map[string]any, error) {
	function, ok := classAbi.Constructor[encoding.ConstructorEntrypoint]
	if !ok {
		return nil, errors.Errorf("unknown constructor")
	}

	return abi.DecodeFunctionCallData(calldata, *function, classAbi.Structs)
}

// CalldataForL1Handler -
func CalldataForL1Handler(contractAbi abi.Abi, selector []byte, calldata []string) (map[string]any, string, error) {
	function, ok := contractAbi.GetL1HandlerBySelector(encoding.EncodeHex(selector))
	if !ok {
		return nil, "", errors.Errorf("unknown selector: %x", selector)
	}

	if len(calldata) == 0 {
		return nil, function.Name, nil
	}

	parsed, err := abi.DecodeFunctionCallData(calldata, *function, contractAbi.Structs)
	return parsed, function.Name, err
}

// Event -
func Event(contractAbi abi.Abi, keys []string, data []string) (map[string]any, string, error) {
	if len(keys) != 1 {
		return nil, "", nil
	}
	selector := encoding.EncodeHex(encoding.MustDecodeHex(keys[0]))
	event, ok := contractAbi.GetEventBySelector(selector)
	if !ok {
		return nil, "", nil
	}
	parsed, err := abi.DecodeEventData(data, *event, contractAbi.Structs)
	return parsed, event.Name, err
}

// func parseEvents(cache *cache.Cache, contractAbi abi.Abi, events []storage.Event) error {
// 	for j := range events {
// 		parsed, name, err := Event(cache, contractAbi, events[j].Keys, events[j].Data)
// 		if err != nil {
// 			return err
// 		}
// 		ptr := &events[j]
// 		ptr.Name = name
// 		ptr.ParsedData = parsed
// 	}
// 	return nil
// }
