package result

import (
	"math/big"

	csquaringmanager "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"

	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
)

// ResultHandler implements the result handler interface
// for processing squared number service results
type ResultHandler struct{}

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
	return []byte(r.Squared), nil
}

// GetDigest computes the digest of the result message for signing
func (p *ResultHandler) GetDigest(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.RequestNumberSquaredOut)
	if !ok {
		return nil, nil
	}

	squared, _ := new(big.Int).SetString(r.Squared, 10)

	// Construct the task response structure
	taskResponse := &csquaringmanager.IIncredibleSquaringServiceManagerTaskResponse{
		ReferenceTaskIndex: r.TaskIndex,
		NumberSquared:      squared,
	}

	// Compute the response digest
	return calcTaskResponseDigest(taskResponse)
}

func calcTaskResponseDigest(h *csquaringmanager.IIncredibleSquaringServiceManagerTaskResponse) ([]byte, error) {
	encodeTaskResponseByte, err := abiEncodeTaskResponse(h)
	if err != nil {
		return nil, err
	}

	var taskResponseDigest [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(encodeTaskResponseByte)
	copy(taskResponseDigest[:], hasher.Sum(nil)[:32])

	return taskResponseDigest[:], nil
}

func abiEncodeTaskResponse(h *csquaringmanager.IIncredibleSquaringServiceManagerTaskResponse) ([]byte, error) {
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

	return arguments.Pack(h)
}
