package contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func getAbiArgs() abi.Arguments {
	return abi.Arguments{
		// {Name: "paymasterMode", Type: abiuint8},
		{Name: "validUntil", Type: uint48},
		{Name: "validAfter", Type: uint48},
		{Name: "erc20Token", Type: address},
		{Name: "postOpGas", Type: uint128},
		{Name: "exchangeRate", Type: uint256},
	}
}

func EncodePaymasterAndData(
	data *Data,
	signature []byte,
) ([]byte, error) {
	// args := getAbiArgs()
	// packed, err := args.Pack(data.ValidUntil, data.ValidAfter, data.ERC20Token, data.ExchangeRate)
	// packed, err := args.Pack(data.ValidUntil, data.ValidAfter, data.ERC20Token, data.PostOpGas, data.ExchangeRate)
	// if err != nil {
	// 	return nil, err
	// }

	concat := data.Paymaster.Bytes()
	concat = append(concat, data.PaymasterMode)
	validUntilBytes := [6]byte{}
	copy(validUntilBytes[2:], data.ValidUntil.Bytes())
	concat = append(concat, validUntilBytes[:]...)
	validAfterBytes := [6]byte{}
	copy(validAfterBytes[2:], data.ValidAfter.Bytes())
	concat = append(concat, validAfterBytes[:]...)
	if data.PaymasterMode == 1 {
		concat = append(concat, data.ERC20Token.Bytes()...)
		concat = append(concat, data.PostOpGas.Bytes()...)
		concat = append(concat, data.ExchangeRate.Bytes()...)
	}
	concat = append(concat, signature...)
	return concat, nil
}
