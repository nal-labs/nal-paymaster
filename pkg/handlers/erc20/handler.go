package erc20

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"github.com/stackup-wallet/stackup-paymaster/pkg/contract"
	"github.com/stackup-wallet/stackup-paymaster/pkg/estimator"
	"github.com/stackup-wallet/stackup-paymaster/pkg/handlers"
)

type Handler struct {
	signer       *signer.EOA
	rpc          *rpc.Client
	eth          *ethclient.Client
	gasEstimator *estimator.GasEstimator
}

func New(
	signer *signer.EOA,
	rpc *rpc.Client,
	eth *ethclient.Client,
	chain *big.Int,
	ov *gas.Overhead,
) *Handler {
	return &Handler{
		rpc:          rpc,
		eth:          eth,
		signer:       signer,
		gasEstimator: estimator.New(signer, rpc, eth, chain, ov),
	}
}

func (h *Handler) Run(
	op *userop.UserOperation,
	ep common.Address,
	pm common.Address,
	tokenAddr common.Address,
) (*handlers.SponsorUserOperationResponse, error) {
	// Get paymaster data.
	data := contract.NewData(pm, tokenAddr, big.NewInt(0), 1, big.NewInt(30000))

	// Estimate gas values to account for paymasterAndData.
	pmOp, err := h.gasEstimator.OverrideOpGasLimitsForPND(op, ep, data)
	if err != nil {
		println("OverrideOpGasLimitsForPND", err.Error())
		return nil, err
	}

	// Fetch hash.
	hash, err := contract.GetHash(h.eth, pmOp, data)
	fmt.Println("GetHash", *pmOp)
	for _, b := range hash {
		println(b) // 分别输出: 65, 66, 67, 68, 69
	}
	if err != nil {
		println("GetHash", err.Error())
		return nil, err
	}

	// Sign hash.
	sig, err := contract.Sign(hash[:], h.signer)
	if err != nil {
		println("Sign", err.Error())
		return nil, err
	}

	// Encode final paymasterAndData.
	pnd, err := contract.EncodePaymasterAndData(data, sig)
	if err != nil {
		println("EncodePaymasterAndData", err.Error())
		return nil, err
	}

	return &handlers.SponsorUserOperationResponse{
		PaymasterAndData:     hexutil.Encode(pnd),
		PreVerificationGas:   hexutil.EncodeBig(pmOp.PreVerificationGas),
		VerificationGasLimit: hexutil.EncodeBig(pmOp.VerificationGasLimit),
		CallGasLimit:         hexutil.EncodeBig(pmOp.CallGasLimit),
	}, nil
}
