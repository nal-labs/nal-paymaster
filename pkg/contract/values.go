package contract

import "github.com/ethereum/go-ethereum/accounts/abi"

var (
	uint256, _  = abi.NewType("uint256", "", nil)
	uint128, _  = abi.NewType("uint128", "", nil)
	uint48, _   = abi.NewType("uint48", "", nil)
	abiuint8, _ = abi.NewType("uint8", "", nil)
	address, _  = abi.NewType("address", "", nil)
)
