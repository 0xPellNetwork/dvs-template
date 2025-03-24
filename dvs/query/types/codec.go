package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces registers the proto interfaces of the current module
// into the proto registry. This allows pellapp-sdk router to correctly
// deserialize the messages.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_QueryService_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &GetDataRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &ListDataRequest{})

	// If you have custom types to register, you can do it here
	// For example: registry.RegisterImplementations((*sdk.Msg)(nil), &MsgQuery{})

	// Maybe you don't need to register any messages for pure query services
	// This function can be left empty, but it's recommended to keep it for consistency

}
