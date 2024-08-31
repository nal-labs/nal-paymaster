package contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func GetHash(
	eth *ethclient.Client,
	op *userop.UserOperation,
	data *Data,
) ([32]byte, error) {
	pm, err := NewContract(data.Paymaster, eth)
	if err != nil {
		println("GetHash-NewContract", err.Error())
		return [32]byte{}, err
	}

	return pm.GetHash(
		&bind.CallOpts{},
		data.PaymasterMode,
		UserOperation(*op),
	)
}
