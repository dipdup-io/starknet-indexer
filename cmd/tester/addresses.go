package main

import "github.com/dipdup-io/starknet-go-api/pkg/encoding"

var (
	nullAddressHash             = encoding.MustDecodeHex("0x0000000000000000000000000000000000000000000000000000000000000000")
	exceptNegativeTokenBalance1 = encoding.MustDecodeHex("0x06a09ccb1caaecf3d9683efe335a667b2169a409d19c589ba1eb771cd210af75")
)
