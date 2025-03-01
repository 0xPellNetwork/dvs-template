package result

import (
	"fmt"

	csquaringManager "github.com/0xPellNetwork/dvs-contracts-template/bindings/IncredibleSquaringServiceManager"
	"github.com/cosmos/gogoproto/proto"

	"github.com/0xPellNetwork/dvs-template/dvs/squared/types"
	"github.com/0xPellNetwork/dvs-template/tools"
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
	responseDigest, err := tools.GetTaskResponseDigest(taskResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get response digest: %v", err)
	}
	return responseDigest[:], nil
}
