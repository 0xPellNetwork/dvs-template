package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/0xPellNetwork/dvs-template/docker/test/dvse2e/services/dvse2e"
)

// CheckBLSAggrSigCmd is the command for checking the BLS aggregation signature of the given group number.
var CheckBLSAggrSigCmd = &cobra.Command{
	Use:  "check-aggr-sigs",
	RunE: handleCheckBLSAggrSig,
}

func init() {
	CheckBLSAggrSigCmdFlagDVSNodeURL.AddToCmdFlag(CheckBLSAggrSigCmd)
	CheckBLSAggrSigCmdFlagDVSServiceManagerAddress.AddToCmdFlag(CheckBLSAggrSigCmd)
	CheckBLSAggrSigCmdFlagGroupNumber.AddToCmdFlag(CheckBLSAggrSigCmd)
	CheckBLSAggrSigCmdFlagThreshold.AddToCmdFlag(CheckBLSAggrSigCmd)

	// flags for trigger new block
	CheckBLSAggrSigCmdFlagETHRPCURL.AddToCmdFlag(CheckBLSAggrSigCmd)
	CheckBLSAggrSigCmdFlagSenderPrivateKey.AddToCmdFlag(CheckBLSAggrSigCmd)
	CheckBLSAggrSigCmdFlagReceiverAddress.AddToCmdFlag(CheckBLSAggrSigCmd)
	CheckBLSAggrSigCmdFlagTimesForTriggerNewBlock.AddToCmdFlag(CheckBLSAggrSigCmd)

	// flags required
	_ = CheckBLSAggrSigCmdFlagDVSNodeURL.MarkRequired(CheckBLSAggrSigCmd)
	_ = CheckBLSAggrSigCmdFlagDVSServiceManagerAddress.MarkRequired(CheckBLSAggrSigCmd)
	_ = CheckBLSAggrSigCmdFlagETHRPCURL.MarkRequired(CheckBLSAggrSigCmd)
}

// handleCheckBLSAggrSig is the handler for the check-bls-sigs command.
// It checks the BLS aggregation signature of the given group number.
func handleCheckBLSAggrSig(cmd *cobra.Command, args []string) error {
	// check flags
	if CheckBLSAggrSigCmdFlagTimesForTriggerNewBlock.Value == 0 {
		CheckBLSAggrSigCmdFlagTimesForTriggerNewBlock.Value = CheckBLSAggrSigCmdFlagTimesForTriggerNewBlock.Default
	}
	return execCheckBLSAggrSig(cmd)
}

// execCheckBLSAggrSig is the executor for the check-bls-sigs command.
// It will create a squaring request and send it to the DVS node.
// It will then query the request by using the hash.
// It will then trigger new blocks for checking the BLS aggregation signature on chain.
// It will then check the BLS aggregation signature of the given group number.
func execCheckBLSAggrSig(cmd *cobra.Command) error {
	ctx := cmd.Context()
	groupNumbers := []uint32{uint32(CheckBLSAggrSigCmdFlagGroupNumber.Value)}
	groupThresholdPercentage := []uint32{uint32(CheckBLSAggrSigCmdFlagThreshold.Value)}

	// create a new PellE2ERunner
	per, err := dvse2e.NewPellE2ERunner(ctx,
		CheckBLSAggrSigCmdFlagETHRPCURL.Value,
		CheckBLSAggrSigCmdFlagDVSNodeURL.Value,
		CheckBLSAggrSigCmdFlagDVSServiceManagerAddress.Value,
		logger,
	)
	if err != nil {
		return err
	}

	// create a squaring request
	req, err := per.PrepareSquaringRequest(ctx, groupNumbers, groupThresholdPercentage)
	if err != nil {
		return err
	}

	logger.Info("Requesting DVS request data",
		"req", req,
	)

	// send to DVS
	reqResp, err := per.RequestDVSAsync(ctx, req)
	if err != nil {
		logger.Error("Failed to request DVS", "error", err)
		return err
	}
	if reqResp == nil {
		logger.Error("Failed to request DVS, reqResp is nil")
		return fmt.Errorf("reqResp is nil")
	}
	if reqResp.Hash == nil {
		logger.Error("Failed to request DVS, reqResp.Hash is nil")
		return fmt.Errorf("reqResp.Hash is nil")
	}

	logger.Info("RequestDVSAsync result", "resp", reqResp)

	var secondsForRequestToBeProcessed = 5 * time.Second
	logger.Info("⌛ Waiting for the request to be processed", "seconds", secondsForRequestToBeProcessed)
	time.Sleep(secondsForRequestToBeProcessed)

	// query request by using hash
	logger.Info("Querying request by using hash", "hash", reqResp.Hash.String())
	taskResult, err := per.QueryRequest(ctx, reqResp.Hash.String())
	if err != nil {
		logger.Error("Failed to query request", "error", err)
		return err
	}

	if taskResult == nil {
		logger.Error("Failed to QueryRequest taskResult is nil")
		return fmt.Errorf("taskResult is nil")
	}

	logger.Info("Got taskResult",
		"hashHex", reqResp.Hash.String(),
		"hash", reqResp,
		"taskResult", taskResult,
	)

	// Test for search request result
	request, err := per.SearchRequest(ctx,
		"SecondEventType.SecondEventKey='SecondEventValue'", nil, nil,
	)
	if err != nil {
		logger.Error("SearchRequest failed",
			"query", "SecondEventType.SecondEventKey='SecondEventValue'",
			"error", err,
		)
		return err
	}

	if request == nil {
		logger.Error("SearchRequest returned no results",
			"query", "SecondEventType.SecondEventKey='SecondEventValue'",
		)
		return fmt.Errorf("search request returned no results")
	} else {
		logger.Info("SearchRequest successful",
			"query", "SecondEventType.SecondEventKey='SecondEventValue'",
			"results", request,
		)
	}

	// Trigger new blocks before check bls sigs, because the contract interface CheckSignatures
	// requires the current block height to be greater than the height when the task was created, but
	// in the test cast, no new blocks are generated, so we make it.
	logger.Info("Triggering new blocks")
	err = per.TriggerAnvilNewBlocks(
		CheckBLSAggrSigCmdFlagTimesForTriggerNewBlock.Value,
		CheckBLSAggrSigCmdFlagSenderPrivateKey.Value,
		CheckBLSAggrSigCmdFlagReceiverAddress.Value,
	)
	if err != nil {
		logger.Error("Failed to trigger new blocks", "error", err)
		return err
	}

	var secondsForNewBlocksToBeGenerated = 2 * time.Second
	logger.Info("⌛ Wainting for new blocks to be generated",
		"seconds", secondsForNewBlocksToBeGenerated,
	)
	time.Sleep(secondsForNewBlocksToBeGenerated)

	// Verify BLS signatures on chain
	logger.Info("Checking BLS signatures on chain after new blocks are generated")
	err = per.VerifyBLSSigsOnChain(req, taskResult)
	if err != nil {
		logger.Error("Failed to verify BLS signatures on chain", "error", err)
		return err
	}

	logger.Info("✅ BLS signatures verified successfully")

	return nil
}
