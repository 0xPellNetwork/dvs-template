package dvse2e

import (
	"context"
	"fmt"
	resulthandlers "github.com/0xPellNetwork/dvs-template/dvs/squared/result"
	sqtypes "github.com/0xPellNetwork/dvs-template/dvs/squared/types"
	"github.com/0xPellNetwork/pellapp-sdk/service/tx"

	sqcontract "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
	"github.com/ethereum/go-ethereum/common"

	"github.com/0xPellNetwork/pelldvs-interactor/chainlibs/eth"
	"github.com/0xPellNetwork/pelldvs/libs/log"
	ctypes "github.com/0xPellNetwork/pelldvs/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// PellDVSE2ERunner is responsible for running end-to-end operations in the DVS system with Ethereum-based integration.
// It handles client communications, manages DVS services, encodes messages, and processes service results.
type PellDVSE2ERunner struct {
	DVSNodeRPCURL string
	logger        log.Logger

	ethRPCURL string

	Client         eth.Client
	serviceManager *sqcontract.ContractIncredibleSquaringServiceManager

	dvsMsgEncoder    tx.MsgEncoder
	dvsResultHandler *resulthandlers.ResultHandler
}

// NewPellE2ERunner initializes and returns a new instance of PellDVSE2ERunner with required dependencies and configurations.
func NewPellE2ERunner(
	ctx context.Context,
	ethRPCURL string,
	DVSNodeRPCURL string,
	serviceManagerAddress string,
	logger log.Logger,
) (*PellDVSE2ERunner, error) {
	per := &PellDVSE2ERunner{
		logger:        logger.With("module", "PellDVSE2ERunner"),
		ethRPCURL:     ethRPCURL,
		DVSNodeRPCURL: DVSNodeRPCURL,
	}

	per.logger.Info("NewPellE2ERunner",
		"DVSNodeRPCURL", DVSNodeRPCURL,
		"ethRPCURL", ethRPCURL,
		"serviceManagerAddress", serviceManagerAddress,
	)

	client, err := eth.NewClient(ethRPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create eth client for RPCURL `%s`: %v", ethRPCURL, err)
	}
	per.Client = client

	manager, err := sqcontract.NewContractIncredibleSquaringServiceManager(
		common.HexToAddress(serviceManagerAddress), client,
	)
	if err != nil {
		return nil, err
	}
	per.serviceManager = manager

	per.dvsMsgEncoder = tx.NewDefaultDecoder(codec.NewProtoCodec(codectypes.NewInterfaceRegistry()))
	per.dvsResultHandler = resulthandlers.NewResultHandler()

	return per, nil
}

// VerifyBLSSigsOnChain verifies BLS aggregated signatures on-chain for a given request and response result.
func (per *PellDVSE2ERunner) VerifyBLSSigsOnChain(
	req *sqtypes.RequestNumberSquaredIn, requestResult *ctypes.ResultDvsRequest,
) error {
	if requestResult == nil {
		return fmt.Errorf("verifyBLSSigsOnChain: nil requestResult")
	}
	if requestResult.DvsResponse == nil {
		return fmt.Errorf("verifyBLSSigsOnChain: nil requestResult.DvsResponse")
	}

	if len(requestResult.DvsResponse.SignersAggSigG1) == 0 {
		return fmt.Errorf("no responses to verify")
	}

	err := per.callVerify(req, requestResult)
	return err
}
