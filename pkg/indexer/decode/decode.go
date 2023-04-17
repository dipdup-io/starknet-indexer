package decode

import (
	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
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

// ResultForFunction -
func ResultForFunction(contractAbi abi.Abi, data []string, selector []byte) (map[string]any, error) {
	if len(data) == 0 {
		return nil, nil
	}

	function, ok := contractAbi.GetFunctionBySelector(encoding.EncodeHex(selector))
	if !ok {
		return nil, errors.Errorf("unknown selector: %x", selector)
	}

	return abi.DecodeFunctionResult(data, *function, contractAbi.Structs)
}

// ResultForL1Handler -
func ResultForL1Handler(contractAbi abi.Abi, data []string, selector []byte) (map[string]any, error) {
	if len(data) == 0 {
		return nil, nil
	}

	function, ok := contractAbi.GetL1HandlerBySelector(encoding.EncodeHex(selector))
	if !ok {
		return nil, errors.Errorf("unknown selector: %x", selector)
	}

	return abi.DecodeFunctionResult(data, *function, contractAbi.Structs)
}

// ResultForConstructor -
func ResultForConstructor(contractAbi abi.Abi, data []string, selector []byte) (map[string]any, error) {
	if len(data) == 0 {
		return nil, nil
	}

	function, ok := contractAbi.GetConstructorBySelector(encoding.EncodeHex(selector))
	if !ok {
		return nil, errors.Errorf("unknown selector: %x", selector)
	}

	return abi.DecodeFunctionResult(data, *function, contractAbi.Structs)
}

// Result -
func Result(contractAbi abi.Abi, data []string, selector []byte, entrypointType storage.EntrypointType) (map[string]any, error) {
	switch entrypointType {
	case storage.EntrypointTypeConstructor:
		return ResultForConstructor(contractAbi, data, selector)
	case storage.EntrypointTypeExternal:
		return ResultForFunction(contractAbi, data, selector)
	case storage.EntrypointTypeL1Handler:
		return ResultForL1Handler(contractAbi, data, selector)
	default:
		return nil, errors.Errorf("unknown entrypoint type in result decoder: %s", entrypointType)
	}
}

// InternalCalldata -
func InternalCalldata(contractAbi abi.Abi, selector []byte, calldata []string, entrypointType storage.EntrypointType) (map[string]any, string, error) {
	switch entrypointType {
	case storage.EntrypointTypeExternal:
		return CalldataBySelector(contractAbi, selector, calldata)
	case storage.EntrypointTypeConstructor:
		data, err := CalldataForConstructor(contractAbi, calldata)
		return data, "", err
	case storage.EntrypointTypeL1Handler:
		return CalldataForL1Handler(contractAbi, selector, calldata)
	default:
		return nil, "", errors.Errorf("unknown entrypoint type in internal calldata decoder: %s", entrypointType)
	}
}
