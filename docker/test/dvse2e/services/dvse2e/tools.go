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

func convertToBN254G2Point(input *bls.G2Point) sqcontract.BN254G2Point {
	output := sqcontract.BN254G2Point{
		X: [2]*big.Int{input.X.A1.BigInt(big.NewInt(0)), input.X.A0.BigInt(big.NewInt(0))},
		Y: [2]*big.Int{input.Y.A1.BigInt(big.NewInt(0)), input.Y.A0.BigInt(big.NewInt(0))},
	}
	return output
}

func convertToBN254G2Point2(p []byte) sqcontract.BN254G2Point {
	if len(p) != 128 {
		panic("invalid G2 point length")
	}
	return sqcontract.BN254G2Point{
		X: [2]*big.Int{
			new(big.Int).SetBytes(p[32:64]), // X[1]
			new(big.Int).SetBytes(p[:32]),   // X[0]
		},
		Y: [2]*big.Int{
			new(big.Int).SetBytes(p[96:]),   // Y[1]
			new(big.Int).SetBytes(p[64:96]), // Y[0]
		},
	}
}
