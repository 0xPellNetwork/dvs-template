package commands

import (
	"github.com/0xPellNetwork/pelldvs/cmd/pelldvs/commands/chains/chainflags"
)

var CheckBLSAggrSigCmdFlagDVSNodeURL = &chainflags.StringFlag{
	Name:    "node-url",
	Usage:   "dvs node url",
	Default: "http://127.0.0.1:26657",
}

var CheckBLSAggrSigCmdFlagGroupNumber = &chainflags.IntFlag{
	Name:    "group",
	Usage:   "group number",
	Default: 0,
}

var CheckBLSAggrSigCmdFlagThreshold = &chainflags.IntFlag{
	Name:    "threshold",
	Usage:   "Threshold",
	Default: 60,
}

var CheckBLSAggrSigCmdFlagDVSServiceManagerAddress = &chainflags.StringFlag{
	Name: "service-manager",
	Aliases: []string{
		"service-manager-address",
	},
}

var CheckBLSAggrSigCmdFlagETHRPCURL = &chainflags.StringFlag{
	Name:    "rpc-url",
	Default: "http://eth:8545",
	Aliases: []string{"eth-rpc-url"},
}

var CheckBLSAggrSigCmdFlagSenderPrivateKey = &chainflags.StringFlag{
	Name:    "sender-private-key",
	Default: "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
}

var CheckBLSAggrSigCmdFlagReceiverAddress = &chainflags.StringFlag{
	Name:    "receiver-address",
	Default: "4860f78301d7ef2dd42a1a4a0a230cc8c38d1996",
}

var CheckBLSAggrSigCmdFlagTimesForTriggerNewBlock = &chainflags.IntFlag{
	Name:    "trigger-times",
	Default: 2,
}
