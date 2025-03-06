package result

import (
	"github.com/cosmos/gogoproto/proto"
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

	return []byte(r.Squared.String()), nil
}

// GetDigest computes the digest of the result message for signing
func (p *ResultHandler) GetDigest(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.RequestNumberSquaredOut)
	if !ok {
		return nil, nil
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(r.Squared.BigInt().Bytes())

	// Compute the response digest
	return hasher.Sum(nil), nil
}
