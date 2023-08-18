package starknet

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/goccy/go-json"
)

var interfaces map[string]abi.Abi

// Interfaces - receives predefined interfaces
func Interfaces(dir string) (map[string]abi.Abi, error) {
	if len(interfaces) != 0 {
		return interfaces, nil
	}
	interfaces = make(map[string]abi.Abi)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for i := range entries {
		info, err := entries[i].Info()
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			continue
		}

		ext := filepath.Ext(entries[i].Name())
		if ext != ".json" {
			continue
		}

		path := filepath.Join(dir, entries[i].Name())
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		var a abi.Abi
		if err := json.NewDecoder(f).Decode(&a); err != nil {
			return nil, err
		}

		name := strings.TrimSuffix(entries[i].Name(), ext)
		interfaces[name] = a

		if err := f.Close(); err != nil {
			return nil, err
		}
	}

	return interfaces, nil
}

// FindInterfaces - returns names of interfaces found in passed abi
func FindInterfaces(a abi.Abi) ([]string, error) {
	result := make([]string, 0)
	for name, i := range interfaces {
		if ok := checkInterface(a, i); ok {
			result = append(result, name)
		}
	}
	return result, nil
}

func checkInterface(a abi.Abi, i abi.Abi) bool {
	for name, iFunc := range i.Functions {
		aFunc, ok := a.Functions[name]
		if !ok {
			return false
		}

		if len(iFunc.Inputs) != len(aFunc.Inputs) {
			return false
		}

		if len(iFunc.Outputs) != len(aFunc.Outputs) {
			return false
		}

		for j := range iFunc.Inputs {
			if iFunc.Inputs[j].Type != aFunc.Inputs[j].Type {
				return false
			}
		}

		for j := range iFunc.Outputs {
			if iFunc.Outputs[j].Type != aFunc.Outputs[j].Type {
				return false
			}
		}
	}

	for name, iFunc := range i.L1Handlers {
		aFunc, ok := a.L1Handlers[name]
		if !ok {
			return false
		}

		if len(iFunc.Inputs) != len(aFunc.Inputs) {
			return false
		}

		if len(iFunc.Outputs) != len(aFunc.Outputs) {
			return false
		}

		for j := range iFunc.Inputs {
			if iFunc.Inputs[j].Type != aFunc.Inputs[j].Type {
				return false
			}
		}

		for j := range iFunc.Outputs {
			if iFunc.Outputs[j].Type != aFunc.Outputs[j].Type {
				return false
			}
		}
	}

	for name, iEvent := range i.Events {
		aEvent, ok := a.Events[name]
		if !ok {
			return false
		}

		if len(iEvent.Data) != len(aEvent.Data) {
			return false
		}

		for j := range iEvent.Data {
			if iEvent.Data[j].Type != aEvent.Data[j].Type {
				return false
			}
		}
	}
	return true
}
