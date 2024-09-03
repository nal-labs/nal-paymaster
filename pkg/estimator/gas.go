package estimator

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
)

func updateOpVerificationGasLimit(op *userop.UserOperation, val *big.Int) (*userop.UserOperation, error) {
	opData, err := op.ToMap()
	if err != nil {
		return nil, err
	}
	opData["verificationGasLimit"] = common.BigToHash(val).String()
	return userop.New(opData)
}

func updateOpCallGasLimit(op *userop.UserOperation, val *big.Int) (*userop.UserOperation, error) {
	opData, err := op.ToMap()
	if err != nil {
		return nil, err
	}
	opData["callGasLimit"] = common.BigToHash(val).String()
	return userop.New(opData)
}

func updateOpPreVerificationGas(
	op *userop.UserOperation,
	ov *gas.Overhead,
) (*userop.UserOperation, error) {
	opData, err := op.ToMap()
	if err != nil {
		return nil, err
	}

	pvg, err := ov.CalcPreVerificationGasWithBuffer(op)
	if err != nil {
		return nil, err
	}

	newPVG := new(big.Int)
	newPVG.Add(pvg, big.NewInt(100))
	opData["preVerificationGas"] = common.BigToHash(newPVG).String()
	return userop.New(opData)
}

func updateOpPaymasterAndData(
	op *userop.UserOperation,
	paymasterAndDataHex string,
) (*userop.UserOperation, error) {
	opData, err := op.ToMap()
	if err != nil {
		return nil, err
	}
	opData["paymasterAndData"] = paymasterAndDataHex
	return userop.New(opData)
}

type GasEstimator struct {
	signer  *signer.EOA
	rpc     *rpc.Client
	eth     *ethclient.Client
	chainID *big.Int
	ov      *gas.Overhead
}

func New(
	signer *signer.EOA,
	rpc *rpc.Client,
	eth *ethclient.Client,
	chain *big.Int,
	ov *gas.Overhead,
) *GasEstimator {
	return &GasEstimator{
		signer:  signer,
		rpc:     rpc,
		eth:     eth,
		chainID: chain,
		ov:      ov,
	}
}

func (g *GasEstimator) OverrideOpGasLimitsForPND(
	op *userop.UserOperation,
	ep common.Address,
	data *contract.Data,
) (*userop.UserOperation, error) {
	// Generate a PND for EstimateGas.
	pndWithoutSig, err := contract.EncodePaymasterAndData(data, []byte{})
	if err != nil {
		return nil, err
	}
	fmt.Printf("%x\n", pndWithoutSig)
	op.PaymasterAndData = pndWithoutSig
	hash, err := contract.GetHash(g.eth, op, data)
	if err != nil {
		println("OverrideOpGasLimitsForPND-GetHash", err.Error())
		return nil, err
	}
	sig, err := contract.Sign(hash[:], g.signer)
	if err != nil {
		println("OverrideOpGasLimitsForPND-Sign", err.Error())
		return nil, err
	}
	pnd, err := contract.EncodePaymasterAndData(data, sig)
	fmt.Printf("%x\n", pnd)
	if err != nil {
		println("OverrideOpGasLimitsForPND-EncodePaymasterAndData", err.Error())
		return nil, err
	}
	println("OverrideOpGasLimitsForPND-EncodePaymasterAndData", pnd)
	pmOp, err := updateOpPaymasterAndData(op, hexutil.Encode(pnd))
	if err != nil {
		println("OverrideOpGasLimitsForPND-updateOpPaymasterAndData", err.Error())
		return nil, err
	}

	// Run EstimateGas.
	vgl, cgl, err := gas.EstimateGas(&gas.EstimateInput{
		Rpc:         g.rpc,
		EntryPoint:  ep,
		Op:          pmOp,
		Ov:          g.ov,
		ChainID:     g.chainID,
		MaxGasLimit: maxGasLimit,
		Tracer:      "",
	})
	if err != nil {
		println("OverrideOpGasLimitsForPND-EstimateGas", err.Error())
		return nil, err
	}

	// Update gas fields.
	// pmOp, err = updateOpPaymasterAndData(pmOp, DummyPaymasterAndDataHex)
	if err != nil {
		println("OverrideOpGasLimitsForPND-updateOpPaymasterAndData", err.Error())
		return nil, err
	}
	pmOp, err = updateOpVerificationGasLimit(pmOp, big.NewInt(int64(vgl)))
	if err != nil {
		println("OverrideOpGasLimitsForPND-updateOpVerificationGasLimit", err.Error())
		return nil, err
	}
	pmOp, err = updateOpCallGasLimit(pmOp, big.NewInt(int64(cgl)))
	if err != nil {
		println("OverrideOpGasLimitsForPND-updateOpCallGasLimit", err.Error())
		return nil, err
	}
	pmOp, err = updateOpPreVerificationGas(pmOp, g.ov)
	if err != nil {
		println("OverrideOpGasLimitsForPND-updateOpPreVerificationGas", err.Error())
		return nil, err
	}

	return pmOp, nil
}
