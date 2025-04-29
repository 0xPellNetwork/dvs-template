package dvse2e

import (
	"context"
	sqtypes "github.com/0xPellNetwork/dvs-template/dvs/squared/types"
	"github.com/0xPellNetwork/pelldvs/rpc/client/http"
	ctypes "github.com/0xPellNetwork/pelldvs/rpc/core/types"
)

// PrepareSquaringRequest prepares a squaring request, the `Squared` field is set to "2".
func (per *PellDVSE2ERunner) PrepareSquaringRequest(
	ctx context.Context,
	groupNumbers []uint32,
	groupThresholdPercentage []uint32,
) (*sqtypes.RequestNumberSquaredIn, error) {
	blockNumber, err := per.Client.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	chainID, err := per.Client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	// convert params groupNumbers >> groupNumbersByteList
	groupNumbersByteList := make([]byte, len(groupNumbers))
	for i, groupNumber := range groupNumbers {
		groupNumbersByteList[i] = byte(groupNumber)
	}

	numberTobeSquared := "2"
	req := &sqtypes.RequestNumberSquaredIn{Task: &sqtypes.TaskRequest{
		TaskIndex:                0,
		Height:                   uint32(blockNumber),
		ChainId:                  chainID.Uint64(),
		Squared:                  numberTobeSquared,
		GroupNumbers:             groupNumbersByteList,
		GroupThresholdPercentage: groupThresholdPercentage[0],
	}}

	return req, nil
}

// RequestDVSAsync sends a squaring request to the DVS node to sign a message.
func (per *PellDVSE2ERunner) RequestDVSAsync(ctx context.Context,
	req *sqtypes.RequestNumberSquaredIn) (
	*ctypes.ResultRequestDvsAsync, error,
) {
	httpClient, err := http.New(per.DVSNodeRPCURL, "")
	if err != nil {
		return nil, err
	}

	data, err := per.dvsMsgEncoder.EncodeMsgs(req)
	if err != nil {
		return nil, err
	}

	var groupNumbers []uint32
	for _, groupNumber := range req.Task.GroupNumbers {
		groupNumbers = append(groupNumbers, uint32(groupNumber))
	}

	result, err := httpClient.RequestDVSAsync(
		ctx,
		data,
		int64(req.Task.Height),
		int64(req.Task.ChainId),
		groupNumbers,
		[]uint32{req.Task.GroupThresholdPercentage},
	)
	return result, err
}

// QueryRequest queries a request by hash.
func (per *PellDVSE2ERunner) QueryRequest(ctx context.Context, hash string) (*ctypes.ResultDvsRequest, error) {
	httpClient, err := http.New(per.DVSNodeRPCURL, "")
	if err != nil {
		return nil, err
	}
	result, err := httpClient.QueryRequest(ctx, hash)
	return result, err
}

// SearchRequest searches for requests.
func (per *PellDVSE2ERunner) SearchRequest(ctx context.Context, query string, pagePtr, perPagePtr *int) (*ctypes.ResultDvsRequestSearch, error) {
	httpClient, err := http.New(per.DVSNodeRPCURL, "")
	if err != nil {
		return nil, err
	}
	result, err := httpClient.SearchRequest(ctx, query, pagePtr, perPagePtr)
	return result, err
}
