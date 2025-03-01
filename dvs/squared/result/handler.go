package result

import (
	"fmt"

	csquaringManager "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"

	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

// ResultHandler implements the result handler interface
// for processing squared number service results
type ResultHandler struct {
}

// NewResultHandler creates a new instance of the result handler
func NewResultHandler() *ResultHandler {
	return &ResultHandler{}
}

// GetData serializes the result message into a byte array
func (p *ResultHandler) GetData(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.RequestNumberSquaredOut)
	if !ok {
		return nil, nil
	}
	return []byte(r.Squared.String()), nil
}

// GetDigest computes the digest of the result message for signing
func (p *ResultHandler) GetDigest(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.RequestNumberSquaredOut)
	if !ok {
		return nil, nil
	}

	// Construct the task response structure
	taskResponse := &csquaringManager.IIncredibleSquaringServiceManagerTaskResponse{
		ReferenceTaskIndex: r.TaskIndex,
		NumberSquared:      r.Squared.BigInt(),
	}

	// Compute the response digest
	responseDigest, err := GetTaskResponseDigest(taskResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get response digest: %v", err)
	}
	return responseDigest[:], nil
}

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
