package types

import (
	"github.com/cometbft/cometbft/libs/bytes"
	"github.com/cometbft/cometbft/types"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// EventAttribute is a single key-value pair, associated with an event.
type EventAttribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Index bool   `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"` // nondeterministic
}

// Event allows application developers to attach additional information to
// ResponseFinalizeBlock and ResponseCheckTx.
// Later, transactions may be queried using these events.
type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type_      string            `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Attributes []*EventAttribute `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes,omitempty"`
}

type ExecTxResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code      uint32   `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Data      []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Log       string   `protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`   // nondeterministic
	Info      string   `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"` // nondeterministic
	GasWanted string   `protobuf:"varint,5,opt,name=gas_wanted,proto3" json:"gas_wanted,omitempty"`
	GasUsed   string   `protobuf:"varint,6,opt,name=gas_used,proto3" json:"gas_used,omitempty"`
	Events    []*Event `protobuf:"bytes,7,rep,name=events,proto3" json:"events,omitempty"` // nondeterministic
	Codespace string   `protobuf:"bytes,8,opt,name=codespace,proto3" json:"codespace,omitempty"`
}

// Result of querying for a tx
type ResultTx struct {
	Hash     bytes.HexBytes `json:"hash"`
	Height   string         `json:"height"`
	Index    uint32         `json:"index"`
	TxResult ExecTxResult   `json:"tx_result"`
	Tx       types.Tx       `json:"tx"`
	Proof    types.TxProof  `json:"proof,omitempty"`
}

// Result of searching for txs
type ResultTxSearch struct {
	Txs        []*ResultTx `json:"txs"`
	TotalCount string      `json:"total_count"`
}

type FeesResponse struct {
	JsonRPC string         `json:"jsonrpc"`
	ID      int32          `json:"id"`
	Result  ResultTxSearch `json:"result"`
}
