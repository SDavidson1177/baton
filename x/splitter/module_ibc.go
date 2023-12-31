package splitter

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"baton/x/splitter/keeper"
	"baton/x/splitter/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	middlewaretypes "github.com/cosmos/ibc-go/v7/modules/apps/30-middleware/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
)

type IBCModule struct {
	keeper       keeper.Keeper
	scopedkeeper ibcexported.ScopedKeeper
}

func NewIBCModule(k keeper.Keeper, s ibcexported.ScopedKeeper) IBCModule {
	return IBCModule{
		keeper:       k,
		scopedkeeper: s,
	}
}

// OnChanOpenInit implements the IBCModule interface
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {

	var metadata middlewaretypes.MiddlewareVersion

	if version != "" {
		// try to unmarshal JSON-encoded version string and pass
		// the app-specific version to app callback.
		// otherwise, pass version directly to app callback.

		err := json.Unmarshal([]byte(version), &metadata)
		if err != nil {
			// call the underlying application's onChanOpenInit callback
			return im.keeper.GetApp().OnChanOpenInit(
				ctx,
				order,
				connectionHops,
				portID,
				channelID,
				chanCap,
				counterparty,
				version,
			)
		}
	} else {
		// TODO: better way to get the app's default version
		metadata = middlewaretypes.MiddlewareVersion{
			MiddlewareVersion: "splitter",
			AppVersion:        "ics20-1",
		}
	}

	// CUSTOM LOGIC GOES HERE

	// call the underlying application's OnChanOpenInit callback.
	// if the version string is empty, OnChanOpenInit is expected to return
	// a default version string representing the version(s) it supports
	appVersion, err := im.keeper.GetApp().OnChanOpenInit(
		ctx,
		order,
		connectionHops,
		portID,
		channelID,
		chanCap,
		counterparty,
		metadata.AppVersion, // The version asso
	)
	if err != nil {
		return "", sdkerrors.Wrapf(types.ErrInvalidVersion, "could not complete app chan open init")
	}

	metadata.AppVersion = appVersion

	// Marshal the version
	version_bytes, err := json.Marshal(metadata)
	if err != nil {
		return "", sdkerrors.Wrapf(types.ErrInvalidVersion, "cannot marshal new version")
	}

	return string(version_bytes), nil
}

// OnChanOpenTry implements the IBCModule interface
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {

	// try to unmarshal JSON-encoded version string and pass
	// the app-specific version to app callback.
	// otherwise, pass version directly to app callback.
	var cpMetadata middlewaretypes.MiddlewareVersion
	err := json.Unmarshal([]byte(counterpartyVersion), &cpMetadata)
	if err != nil {
		// call the underlying application's OnChanOpenTry callback
		return im.keeper.GetApp().OnChanOpenTry(
			ctx,
			order,
			connectionHops,
			portID,
			channelID,
			chanCap,
			counterparty,
			counterpartyVersion,
		)
	}

	// select mutually compatible middleware version
	// TODO: Check this
	// if !isCompatible(cpMetadata.middlewareVersion) {
	// 	return "", error
	// }
	// middlewareVersion = selectMiddlewareVersion(cpMetadata.middlewareVersion)

	// call the underlying application's OnChanOpenTry callback
	appVersion, err := im.keeper.GetApp().OnChanOpenTry(
		ctx,
		order,
		connectionHops,
		portID,
		channelID,
		chanCap,
		counterparty,
		cpMetadata.AppVersion, // note we only pass counterparty app version here
	)
	if err != nil {
		return "", sdkerrors.Wrapf(types.ErrInvalidVersion, "could not complete app chan open init")
	}

	cpMetadata.AppVersion = appVersion

	// Marshal the version
	version_bytes, err := json.Marshal(cpMetadata)
	if err != nil {
		return "", sdkerrors.Wrapf(types.ErrInvalidVersion, "cannot marshal new version")
	}

	return string(version_bytes), nil
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartChannelID string,
	counterpartyVersion string,
) error {
	var cpMetadata middlewaretypes.MiddlewareVersion
	err := json.Unmarshal([]byte(counterpartyVersion), &cpMetadata)
	if err != nil {
		// call the underlying application's OnChanOpenAck callback
		return im.keeper.GetApp().OnChanOpenAck(
			ctx,
			portID,
			channelID,
			counterpartChannelID,
			counterpartyVersion,
		)
	}

	// TODO: Check this
	// if !isSupported(cpMetadata.middlewareVersion) {
	// 	return error
	// }
	// doCustomLogic()

	// call the underlying application's OnChanOpenAck callback
	return im.keeper.GetApp().OnChanOpenAck(
		ctx,
		portID,
		channelID,
		counterpartChannelID,
		cpMetadata.AppVersion,
	)
}

func (im IBCModule) UpdateEndpointChainID(ctx sdk.Context, port string, channel string) error {
	// Packet to update the chain id
	sp := types.SplitterPacket{
		ChainId: ctx.ChainID(),
	}

	capKey, _ := im.scopedkeeper.GetCapability(ctx, host.ChannelCapabilityPath(port, channel))

	data, err_mar := sp.Marshal()
	if err_mar != nil {
		return err_mar
	}

	now := uint64(time.Now().UnixNano())

	// TODO: Determine what the timeout should be
	timeout_t := now + 6*uint64(time.Hour)

	_, err := im.keeper.SendPacketBypass(ctx, capKey, port, channel, clienttypes.Height{RevisionNumber: 0, RevisionHeight: 0}, timeout_t, data)
	if err != nil {
		fmt.Printf("error for send %v\n", err.Error())
	}

	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	err := im.keeper.GetApp().OnChanOpenConfirm(ctx, portID, channelID)
	if err == nil {
		return im.UpdateEndpointChainID(ctx, portID, channelID)
	}
	return err
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Disallow user-initiated channel closing for channels
	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	var splitter_packet types.SplitterPacket
	err := splitter_packet.Unmarshal(modulePacket.Data)

	if err != nil {
		// An error has occurred
		fmt.Printf("cannot unmarshal splitter packet")
		return nil
	}

	// TODO: Unify packets to avoid this silly hack
	switch splitter_packet.Type {
	case types.TYPE_HANDSHAKE:
		// Store the port/channel and chain id pair to the KV store
		store := im.keeper.GetStore(ctx)

		// Get the current
		var cc_map types.ChannelChainMap
		map_bytes := store.Get([]byte(types.ChannelChainKey))

		err = cc_map.Unmarshal(map_bytes)
		if err != nil {
			// error
			fmt.Printf("cannot unmarshal channel chain map")
			return nil
		}

		// Check if the channel is already accounted for
		for _, v := range cc_map.Values {
			if v.Port == modulePacket.DestinationPort && v.Channel == modulePacket.DestinationChannel {
				return nil
			}
		}

		// Track endpoint info about the channel
		cc_map.Values = append(cc_map.Values, &types.ChannelChain{
			Port:    modulePacket.DestinationPort,
			Channel: modulePacket.DestinationChannel,
			Chain:   splitter_packet.ChainId,
		})

		map_bytes_set, err := cc_map.Marshal()
		if err != nil {
			// error
			fmt.Printf("cannot marshal channel chain map")
			return nil
		}

		store.Set([]byte(types.ChannelChainKey), map_bytes_set)

		// Send a response
		err = im.UpdateEndpointChainID(ctx, modulePacket.SourcePort, modulePacket.SourceChannel)
		if err != nil {
			fmt.Printf("error resending: %v\n", err.Error())
		}

		// Acknowledge the packet and return
		// return channeltypes.NewResultAcknowledgement([]byte{byte(1)})
		return nil
	case types.TYPE_WRAPPER:
		// TODO: Need to prevent the scenario where the same malicious channel sends the same
		// packet over and over again to increment amount.

		// Increment the number of packets received
		// Use packet hash has key.
		store := im.keeper.GetStore(ctx)
		hash := sha256.Sum256(modulePacket.GetData())

		packet_val := store.Get(hash[:])
		if packet_val == nil {
			// Create a new entry to track this packet
			d := types.SplitterPacketTracker{
				Port:    modulePacket.SourcePort,
				Channel: modulePacket.SourceChannel,
				Amount:  1,
			}

			d_bytes, err := d.Marshal()
			if err != nil {
				return nil
			}

			store.Set(hash[:], d_bytes)

			// Do not process packet yet
			return nil
		}

		// Update the packet tracker
		var packet_tracker types.SplitterPacketTracker
		err = packet_tracker.Unmarshal(packet_val)
		if err != nil {
			return nil
		}

		// TODO: Determine how many packets must be received
		// Counter can only be above a certain threshold N if
		// packet is correctly sent OR there are colluding
		// malicious chains on N channels.
		if packet_tracker.Amount+1 < 2 {
			// don't send packet yet
			packet_tracker.Amount++
			tracker_bytes, err := packet_tracker.Marshal()
			if err == nil {
				store.Set(hash[:], tracker_bytes)
			}
			return nil
		}

		// Send the packet. Update amount in case packet with the same
		// hash is sent later
		packet_tracker.Amount = 0
		tracker_bytes, err := packet_tracker.Marshal()
		if err == nil {
			store.Set(hash[:], tracker_bytes)
		}

		// Send the packet
		modulePacket.Data = splitter_packet.PacketData
		modulePacket.SourcePort = splitter_packet.SourcePort
		modulePacket.SourceChannel = splitter_packet.SourceChannel
		modulePacket.DestinationPort = splitter_packet.DstPort
		modulePacket.DestinationChannel = splitter_packet.DstChannel
		return im.keeper.GetApp().OnRecvPacket(ctx, modulePacket, relayer)
	}

	// this line is used by starport scaffolding # oracle/packet/module/recv
	return nil
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	// TODO: Fix acknowledgements

	return im.keeper.GetApp().OnAcknowledgementPacket(ctx, modulePacket, acknowledgement, relayer)
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return im.keeper.GetApp().OnTimeoutPacket(ctx, modulePacket, relayer)
}
