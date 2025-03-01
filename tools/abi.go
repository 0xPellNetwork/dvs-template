package tools

import (
	csquaringManager "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"
)

func AbiEncodeTaskResponse(h *csquaringManager.IIncredibleSquaringServiceManagerTaskResponse) ([]byte, error) {

	taskResponseType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "referenceTaskIndex",
			Type: "uint32",
		},
		{
			Name: "numberSquared",
			Type: "uint256",
		},
	})
	if err != nil {
		return nil, err
	}

	arguments := abi.Arguments{
		{
			Type: taskResponseType,
		},
	}

	bytes, err := arguments.Pack(h)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func GetTaskResponseDigest(h *csquaringManager.IIncredibleSquaringServiceManagerTaskResponse) ([32]byte, error) {
	encodeTaskResponseByte, err := AbiEncodeTaskResponse(h)
	if err != nil {
		return [32]byte{}, err
	}

	var taskResponseDigest [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(encodeTaskResponseByte)
	copy(taskResponseDigest[:], hasher.Sum(nil)[:32])

	return taskResponseDigest, nil
}
