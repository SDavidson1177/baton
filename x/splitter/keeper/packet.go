package keeper

import (
	"baton/x/splitter/types"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

func (k Keeper) SendPacket(
	ctx sdk.Context,
	channelCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (uint64, error) {
	channel, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return 0, fmt.Errorf("cannot find channel with port and channel: %v/%v", sourcePort, sourceChannel)
	}

	// Wrap the data
	wrapped_data := types.SplitterPacket{
		Type:          types.TYPE_WRAPPER,
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		DstPort:       channel.Counterparty.PortId,
		DstChannel:    channel.Counterparty.ChannelId,
		PacketData:    data,
	}

	wrapped_data_bytes, err := wrapped_data.Marshal()
	if err != nil {
		return 0, err
	}

	store := k.GetStore(ctx)

	// Get the current
	var cc_map types.ChannelChainMap
	map_bytes := store.Get([]byte(types.ChannelChainKey))

	err = cc_map.Unmarshal(map_bytes)
	if err != nil {
		// error
		return 0, err
	}

	// Check if the channel is already accounted for
	found = false
	var cc string

	for _, v := range cc_map.Values {
		if v.Port == sourcePort && v.Channel == sourceChannel {
			found = true
			cc = v.Chain
			break
		}
	}

	// Send on all channels
	if found {
		for _, v := range cc_map.Values {
			if v.Chain == cc {
				cap, f := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(v.Port, v.Channel))
				if f {
					k.channelKeeper.SendPacket(ctx, cap, v.Port, v.Channel, timeoutHeight, timeoutTimestamp, wrapped_data_bytes)
				}
			}
		}
	}

	// default ack
	return 1, nil
}

func (k Keeper) SendPacketBypass(
	ctx sdk.Context,
	channelCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (uint64, error) {
	return k.channelKeeper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

func (k Keeper) WriteAcknowledgement(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet exported.PacketI,
	acknowledgement exported.Acknowledgement,
) error {
	return k.channelKeeper.WriteAcknowledgement(ctx, chanCap, packet, acknowledgement)
}
