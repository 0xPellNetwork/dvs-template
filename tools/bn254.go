package tools

import (
	"math/big"

	csquaringManager "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
	"github.com/0xPellNetwork/pelldvs/crypto/bls"
)

func ConvertToBN254G1Point(input *bls.G1Point) csquaringManager.BN254G1Point {
	output := csquaringManager.BN254G1Point{
		X: input.X.BigInt(big.NewInt(0)),
		Y: input.Y.BigInt(big.NewInt(0)),
	}
	return output
}

func ConvertToBN254G2Point(input *bls.G2Point) csquaringManager.BN254G2Point {
	output := csquaringManager.BN254G2Point{
		X: [2]*big.Int{input.X.A1.BigInt(big.NewInt(0)), input.X.A0.BigInt(big.NewInt(0))},
		Y: [2]*big.Int{input.Y.A1.BigInt(big.NewInt(0)), input.Y.A0.BigInt(big.NewInt(0))},
	}
	return output
}
