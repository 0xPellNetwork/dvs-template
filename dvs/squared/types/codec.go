package types

import (
	dvstypes "github.com/0xPellNetwork/pellapp-sdk/pelldvs/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces registers the proto interfaces of the current module
// into the proto registry. This allows pellapp-sdk router to correctly
// deserialize the messages.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_DVSRequest_serviceDesc)
	msgservice.RegisterMsgServiceDesc(registry, &_DVSResponse_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &dvstypes.RequestPostRequestValidatedData{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &RequestNumberSquaredIn{})
}
