package decode

import (
	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/pkg/errors"
)

var (
	ErrUnknownSelector = errors.New("unknown selector")
)

// CalldataBySelector -
func CalldataBySelector(contractAbi abi.Abi, selector []byte, calldata []string) (map[string]any, string, error) {
	function, ok := contractAbi.GetFunctionBySelector(encoding.EncodeHex(selector))
	if !ok {
		return nil, "", errors.Wrapf(ErrUnknownSelector, "function %x", selector)
	}

	if len(calldata) == 0 {
		return nil, function.Name, nil
	}

	parsed, err := abi.DecodeFunctionCallData(calldata, *function, contractAbi.Structs, contractAbi.Enums)
	return parsed, function.Name, err
}

// CalldataForConstructor -
func CalldataForConstructor(classAbi abi.Abi, calldata []string) (map[string]any, error) {
	function, ok := classAbi.Constructor[encoding.ConstructorEntrypoint]
	if !ok {
		function, ok = classAbi.Functions[encoding.ConstructorEntrypoint]
		if !ok {
			return nil, errors.Errorf("unknown constructor")
		}
	}

	return abi.DecodeFunctionCallData(calldata, *function, classAbi.Structs, classAbi.Enums)
}

// CalldataForL1Handler -
func CalldataForL1Handler(contractAbi abi.Abi, selector []byte, calldata []string) (map[string]any, string, error) {
	function, ok := contractAbi.GetL1HandlerBySelector(encoding.EncodeHex(selector))
	if !ok {
		return CalldataBySelector(contractAbi, selector, calldata)
	}

	if len(calldata) == 0 {
		return nil, function.Name, nil
	}

	parsed, err := abi.DecodeFunctionCallData(calldata, *function, contractAbi.Structs, contractAbi.Enums)
	return parsed, function.Name, err
}

// Event -
func Event(contractAbi abi.Abi, keys []string, data []string) (map[string]any, string, error) {
	if len(keys) == 0 {
		return nil, "", nil
	}
	selector := encoding.EncodeHex(encoding.MustDecodeHex(keys[0]))
	event, ok := contractAbi.GetEventBySelector(selector)
	if !ok {
		return nil, "", nil
	}
	var values []string
	switch len(keys) {
	case 1:
		values = make([]string, len(data))
		copy(values, data)
	default:
		if len(event.Members) == 0 {
			values = make([]string, len(keys[1:]))
			copy(values, keys[1:])
			values = append(values, data...)
		} else {
			var dataIdx int
			var keysIdx = 1
			for i := range event.Members {
				switch event.Members[i].Kind {
				case "data":
					if len(data) > dataIdx {
						values = append(values, data[dataIdx])
						dataIdx++
					}
				case "key":
					if len(keys) > keysIdx {
						values = append(values, keys[keysIdx])
						keysIdx++
					}
				default:
					return nil, "", errors.Errorf("unknown event member kind: %s", event.Members[i].Kind)
				}
			}
		}
	}

	parsed, err := abi.DecodeEventData(values, *event, contractAbi.Structs, contractAbi.Enums)
	return parsed, event.Name, err
}

// ResultForFunction -
func ResultForFunction(contractAbi abi.Abi, data []string, selector []byte) (map[string]any, error) {
	if len(data) == 0 {
		return nil, nil
	}

	function, ok := contractAbi.GetFunctionBySelector(encoding.EncodeHex(selector))
	if !ok {
		return nil, errors.Wrapf(ErrUnknownSelector, "function %x", selector)
	}

	return abi.DecodeFunctionResult(data, *function, contractAbi.Structs, contractAbi.Enums)
}

// ResultForL1Handler -
func ResultForL1Handler(contractAbi abi.Abi, data []string, selector []byte) (map[string]any, error) {
	if len(data) == 0 {
		return nil, nil
	}

	function, ok := contractAbi.GetL1HandlerBySelector(encoding.EncodeHex(selector))
	if !ok {
		return nil, errors.Wrapf(ErrUnknownSelector, "l1_handler %x", selector)
	}

	return abi.DecodeFunctionResult(data, *function, contractAbi.Structs, contractAbi.Enums)
}

// ResultForConstructor -
func ResultForConstructor(contractAbi abi.Abi, data []string, selector []byte) (map[string]any, error) {
	if len(data) == 0 {
		return nil, nil
	}

	function, ok := contractAbi.GetConstructorBySelector(encoding.EncodeHex(selector))
	if !ok {
		return nil, errors.Wrapf(ErrUnknownSelector, "constructor %x", selector)
	}

	return abi.DecodeFunctionResult(data, *function, contractAbi.Structs, contractAbi.Enums)
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
