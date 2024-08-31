package contract

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Data struct {
	Paymaster     common.Address
	PaymasterMode uint8
	ValidUntil    *big.Int
	ValidAfter    *big.Int
	ERC20Token    common.Address
	PostOpGas     *big.Int
	ExchangeRate  *big.Int
}

func NewData(pm common.Address, token common.Address, exchangeRate *big.Int, pmMode uint8, postOpGas *big.Int) *Data {
	return &Data{
		Paymaster:     pm,
		PaymasterMode: pmMode,
		ValidUntil:    big.NewInt(int64(time.Now().Add(time.Hour).Unix())),
		ValidAfter:    big.NewInt(0),
		ERC20Token:    token,
		PostOpGas:     postOpGas,
		ExchangeRate:  exchangeRate,
	}
}
