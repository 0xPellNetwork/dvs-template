package dvse2e

import (
	"context"
	"fmt"
	"math/big"

	sqcontract "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
	"github.com/0xPellNetwork/pelldvs/crypto/bls"
	coretypes "github.com/0xPellNetwork/pelldvs/rpc/core/types"
	"github.com/pkg/errors"

	sqtypes "github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

// callVerify calls the Verify method on the service manager contract.
func (per *PellDVSE2ERunner) callVerify(req *sqtypes.RequestNumberSquaredIn,
	requestResult *coretypes.ResultDvsRequest,
) error {
	per.logger.Info("callVerify", "requestResult", requestResult)

	// Construct NonSignerStakesAndSignature parameters
	nonSignerPubkeysG1 := make([]sqcontract.BN254G1Point, len(requestResult.DvsResponse.NonSignersPubkeysG1))
	for i, pubkey := range requestResult.DvsResponse.NonSignersPubkeysG1 {
		nonSignerPubkeysG1[i] = sqcontract.BN254G1Point{
			X: new(big.Int).SetBytes(pubkey[:32]),
			Y: new(big.Int).SetBytes(pubkey[32:]),
		}
	}

	var groupApksG1 []sqcontract.BN254G1Point
	for _, apk := range requestResult.DvsResponse.GroupApksG1 {
		tapk := bls.NewZeroG1Point()
		_ = tapk.Unmarshal(apk)
		groupApksG1 = append(groupApksG1, convertToBN254G1Point(tapk))
	}

	signersApkG2 := sqcontract.BN254G2Point{
		X: [2]*big.Int{
			new(big.Int).SetBytes(requestResult.DvsResponse.SignersApkG2[:32]),
			new(big.Int).SetBytes(requestResult.DvsResponse.SignersApkG2[32:64]),
		},
		Y: [2]*big.Int{
			new(big.Int).SetBytes(requestResult.DvsResponse.SignersApkG2[64:96]),
			new(big.Int).SetBytes(requestResult.DvsResponse.SignersApkG2[96:]),
		},
	}

	signersAggSigG1 := sqcontract.BN254G1Point{
		X: new(big.Int).SetBytes(requestResult.DvsResponse.SignersAggSigG1[:32]),
		Y: new(big.Int).SetBytes(requestResult.DvsResponse.SignersAggSigG1[32:]),
	}

	nonSignerStakeIndices := make([][]uint32, len(requestResult.DvsResponse.NonSignerStakeIndices))
	for i, indices := range requestResult.DvsResponse.NonSignerStakeIndices {
		nonSignerStakeIndices[i] = indices.NonSignerStakeIndice
	}

	nonSignerStakesAndSignature := sqcontract.IBLSSignatureVerifierNonSignerStakesAndSignature{
		NonSignerPubkeys:            nonSignerPubkeysG1,
		GroupApks:                   groupApksG1,
		ApkG2:                       signersApkG2,
		Sigma:                       signersAggSigG1,
		NonSignerGroupBitmapIndices: requestResult.DvsResponse.NonSignerGroupBitmapIndices,
		GroupApkIndices:             requestResult.DvsResponse.GroupApkIndices,
		TotalStakeIndices:           requestResult.DvsResponse.TotalStakeIndices,
		NonSignerStakeIndices:       nonSignerStakeIndices,
	}

	blockNumber := uint32(requestResult.DvsRequest.Height)
	currentBlockNumber, _ := per.Client.BlockNumber(context.TODO())

	var groupNumbers []byte
	for _, groupNumber := range requestResult.DvsRequest.GroupNumbers {
		groupNumbers = append(groupNumbers, byte(groupNumber))
	}

	// We need to manually construct sqtypes.RequestNumberSquaredOut to simulate the behavior of the DVS node
	// Here we assume the DVS node has already calculated the square result
	squared, _ := new(big.Int).SetString(req.Task.Squared, 10)
	squared = squared.Mul(squared, big.NewInt(int64(2)))

	var resp = &sqtypes.RequestNumberSquaredOut{
		TaskIndex: req.Task.TaskIndex,
		Squared:   squared.String(),
	}
	per.logger.Info("callVerify", "resp", resp)
	msgHash, err := per.dvsResultHandler.GetDigest(resp)
	if err != nil {
		return err
	}

	if len(msgHash) == 0 {
		return errors.New("no response digest")
	}

	var msgHashBytes [32]byte
	copy(msgHashBytes[:], msgHash)

	fmt.Println()
	fmt.Println()

	per.logger.Info("params for signature verification",
		"taskBlockNumber", blockNumber,
		"currentBlockNumber", currentBlockNumber,
		"msgHash", msgHash,
		"groupNumbers", groupNumbers,
		"nonSignerStakesAndSignature", nonSignerStakesAndSignature,
	)

	_, _, err = per.serviceManager.CheckSignatures(nil, msgHashBytes,
		groupNumbers, blockNumber, nonSignerStakesAndSignature,
	)
	if err != nil {
		per.logger.Error("Signature verification failed",
			"error", err,
			"blockNumber", blockNumber,
			"currentBlock", currentBlockNumber,
		)
		return err
	}

	per.logger.Info("✅ successfully verified aggregate signature")

	return nil
}
