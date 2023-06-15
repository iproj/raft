// Code generated by protoc-gen-go.
// source: request_vote_request.proto
// DO NOT EDIT!

package protobuf

import proto "github.com/golang/protobuf/proto"
import math "math"

// discarding unused import gogoproto "github.com/golang/protobuf/gogoproto/gogo.pb"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type RequestVoteRequest struct {
	Term             *uint64 `protobuf:"varint,1,req" json:"Term,omitempty"`
	LastLogIndex     *uint64 `protobuf:"varint,2,req" json:"LastLogIndex,omitempty"`
	LastLogTerm      *uint64 `protobuf:"varint,3,req" json:"LastLogTerm,omitempty"`
	CandidateName    *string `protobuf:"bytes,4,req" json:"CandidateName,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *RequestVoteRequest) Reset()         { *m = RequestVoteRequest{} }
func (m *RequestVoteRequest) String() string { return proto.CompactTextString(m) }
func (*RequestVoteRequest) ProtoMessage()    {}

func (m *RequestVoteRequest) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *RequestVoteRequest) GetLastLogIndex() uint64 {
	if m != nil && m.LastLogIndex != nil {
		return *m.LastLogIndex
	}
	return 0
}

func (m *RequestVoteRequest) GetLastLogTerm() uint64 {
	if m != nil && m.LastLogTerm != nil {
		return *m.LastLogTerm
	}
	return 0
}

func (m *RequestVoteRequest) GetCandidateName() string {
	if m != nil && m.CandidateName != nil {
		return *m.CandidateName
	}
	return ""
}

func init() {
}
