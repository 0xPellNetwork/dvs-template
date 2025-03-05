package server

import (
	"context"
	"math/big"

	"cosmossdk.io/math"
	csquaringManager "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
	"github.com/0xPellNetwork/pellapp-sdk/pelldvs"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
	"github.com/0xPellNetwork/pelldvs/crypto/bls"

	chainConnector "github.com/0xPellNetwork/dvs-template/chain_connector"
	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

var ChainConnector *chainConnector.Client

func (s Server) ResponseNumberSquared(ctx context.Context, in *types.RequestNumberSquaredIn) (*types.ResponseNumberSquaredOut, error) {
	pkgCtx := sdktypes.UnwrapContext(ctx)

	validatedData, err := pelldvs.GetDvsRequestValidatedData(pkgCtx)
	if err != nil {
		return nil, err
	}

	// Convert []uint32 to bytes
	groupNumbersBytes := make([]byte, len(in.Task.GroupNumbers))
	for i, num := range in.Task.GroupNumbers {
		groupNumbersBytes[i] = byte(num)
	}

	squared, _ := math.NewIntFromString(string(validatedData.Data))
	// Construct task parameters
	task := csquaringManager.IIncredibleSquaringServiceManagerTask{
		NumberToBeSquared:        in.Task.Squared.BigInt(),
		TaskCreatedBlock:         in.Task.Height,
		GroupNumbers:             groupNumbersBytes,
		GroupThresholdPercentage: in.Task.GroupThresholdPercentage,
	}

	// Construct TaskResponse parameters
	taskResponse := csquaringManager.IIncredibleSquaringServiceManagerTaskResponse{
		ReferenceTaskIndex: in.Task.TaskIndex,
		NumberSquared:      squared.BigInt(),
	}

	// Construct NonSignerStakesAndSignature parameters
	nonSignerPubkeysG1 := make([]csquaringManager.BN254G1Point, len(validatedData.NonSignersPubkeysG1))
	for i, pubkey := range validatedData.NonSignersPubkeysG1 {
		nonSignerPubkeysG1[i] = csquaringManager.BN254G1Point{
			X: new(big.Int).SetBytes(pubkey[:32]),
			Y: new(big.Int).SetBytes(pubkey[32:]),
		}
	}

	quorumApksG1 := []csquaringManager.BN254G1Point{}
	for _, apk := range validatedData.QuorumApksG1 {
		tapk := bls.NewZeroG1Point()
		_ = tapk.Unmarshal(apk)
		quorumApksG1 = append(quorumApksG1, csquaringManager.BN254G1Point{
			X: tapk.X.BigInt(big.NewInt(0)),
			Y: tapk.Y.BigInt(big.NewInt(0)),
		})
	}

	signersAggSigG1 := csquaringManager.BN254G1Point{
		X: new(big.Int).SetBytes(validatedData.SignersAggSigG1[:32]),
		Y: new(big.Int).SetBytes(validatedData.SignersAggSigG1[32:]),
	}

	nonSignerStakeIndices := make([][]uint32, len(validatedData.NonSignerStakeIndices))
	for i, indices := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices[i] = indices.NonSignerStakeIndice
	}

	signersApkG2 := csquaringManager.BN254G2Point{
		X: [2]*big.Int{
			new(big.Int).SetBytes(validatedData.SignersApkG2[:32]),
			new(big.Int).SetBytes(validatedData.SignersApkG2[32:64]),
		},
		Y: [2]*big.Int{
			new(big.Int).SetBytes(validatedData.SignersApkG2[64:96]),
			new(big.Int).SetBytes(validatedData.SignersApkG2[96:]),
		},
	}

	nonSignerStakesAndSignature := csquaringManager.IBLSSignatureVerifierNonSignerStakesAndSignature{
		NonSignerPubkeys:            nonSignerPubkeysG1,
		GroupApks:                   quorumApksG1,
		ApkG2:                       signersApkG2,
		Sigma:                       signersAggSigG1,
		NonSignerGroupBitmapIndices: validatedData.NonSignerQuorumBitmapIndices,
		GroupApkIndices:             validatedData.QuorumApkIndices,
		TotalStakeIndices:           validatedData.TotalStakeIndices,
		NonSignerStakeIndices:       nonSignerStakeIndices,
	}

	err = ChainConnector.RespondToTask(uint64(pkgCtx.ChainID()), task, taskResponse, nonSignerStakesAndSignature)
	if err != nil {
		s.logger.Error("Failed to respond to task", "error", err)
		return nil, err
	}

	return &types.ResponseNumberSquaredOut{}, nil
}
