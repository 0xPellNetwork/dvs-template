package dvse2e

import (
	"math/big"

	sqcontract "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
	"github.com/0xPellNetwork/pelldvs/crypto/bls"
)

func convertToBN254G1Point(input *bls.G1Point) sqcontract.BN254G1Point {
	output := sqcontract.BN254G1Point{
		X: input.X.BigInt(big.NewInt(0)),
		Y: input.Y.BigInt(big.NewInt(0)),
	}
	return output
}
