// Code generated by protoc-gen-go.
// source: snapshot_recovery_request.proto
// DO NOT EDIT!

package protobuf

import proto "github.com/golang/protobuf/proto"
import math "math"

// discarding unused import gogoproto "github.com/golang/protobuf/gogoproto/gogo.pb"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type SnapshotRecoveryRequest struct {
	LeaderName       *string                         `protobuf:"bytes,1,req" json:"LeaderName,omitempty"`
	LastIndex        *uint64                         `protobuf:"varint,2,req" json:"LastIndex,omitempty"`
	LastTerm         *uint64                         `protobuf:"varint,3,req" json:"LastTerm,omitempty"`
	Peers            []*SnapshotRecoveryRequest_Peer `protobuf:"bytes,4,rep" json:"Peers,omitempty"`
	State            []byte                          `protobuf:"bytes,5,req" json:"State,omitempty"`
	XXX_unrecognized []byte                          `json:"-"`
}

func (m *SnapshotRecoveryRequest) Reset()         { *m = SnapshotRecoveryRequest{} }
func (m *SnapshotRecoveryRequest) String() string { return proto.CompactTextString(m) }
func (*SnapshotRecoveryRequest) ProtoMessage()    {}

func (m *SnapshotRecoveryRequest) GetLeaderName() string {
	if m != nil && m.LeaderName != nil {
		return *m.LeaderName
	}
	return ""
}

func (m *SnapshotRecoveryRequest) GetLastIndex() uint64 {
	if m != nil && m.LastIndex != nil {
		return *m.LastIndex
	}
	return 0
}

func (m *SnapshotRecoveryRequest) GetLastTerm() uint64 {
	if m != nil && m.LastTerm != nil {
		return *m.LastTerm
	}
	return 0
}

func (m *SnapshotRecoveryRequest) GetPeers() []*SnapshotRecoveryRequest_Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}

func (m *SnapshotRecoveryRequest) GetState() []byte {
	if m != nil {
		return m.State
	}
	return nil
}

type SnapshotRecoveryRequest_Peer struct {
	Name             *string `protobuf:"bytes,1,req" json:"Name,omitempty"`
	ConnectionString *string `protobuf:"bytes,2,req" json:"ConnectionString,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SnapshotRecoveryRequest_Peer) Reset()         { *m = SnapshotRecoveryRequest_Peer{} }
func (m *SnapshotRecoveryRequest_Peer) String() string { return proto.CompactTextString(m) }
func (*SnapshotRecoveryRequest_Peer) ProtoMessage()    {}

func (m *SnapshotRecoveryRequest_Peer) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *SnapshotRecoveryRequest_Peer) GetConnectionString() string {
	if m != nil && m.ConnectionString != nil {
		return *m.ConnectionString
	}
	return ""
}

func init() {
}
