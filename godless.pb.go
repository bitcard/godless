// Code generated by protoc-gen-go. DO NOT EDIT.
// source: godless.proto

/*
Package godless is a generated protocol buffer package.

It is generated from these files:
	godless.proto

It has these top-level messages:
	NamespaceMessage
	NamespaceEntryMessage
	IndexMessage
	IndexEntryMessage
	APIResponseMessage
	APIQueryResponseMessage
	APIReflectResponseMessage
	QueryMessage
	QueryJoinMessage
	QueryRowJoinMessage
	QueryRowJoinEntryMessage
	QuerySelectMessage
	QueryWhereMessage
	QueryPredicateMessage
*/
package godless

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type NamespaceMessage struct {
	Entries []*NamespaceEntryMessage `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
}

func (m *NamespaceMessage) Reset()                    { *m = NamespaceMessage{} }
func (m *NamespaceMessage) String() string            { return proto.CompactTextString(m) }
func (*NamespaceMessage) ProtoMessage()               {}
func (*NamespaceMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *NamespaceMessage) GetEntries() []*NamespaceEntryMessage {
	if m != nil {
		return m.Entries
	}
	return nil
}

type NamespaceEntryMessage struct {
	Table  string   `protobuf:"bytes,1,opt,name=table" json:"table,omitempty"`
	Row    string   `protobuf:"bytes,2,opt,name=row" json:"row,omitempty"`
	Entry  string   `protobuf:"bytes,3,opt,name=entry" json:"entry,omitempty"`
	Points []string `protobuf:"bytes,4,rep,name=points" json:"points,omitempty"`
}

func (m *NamespaceEntryMessage) Reset()                    { *m = NamespaceEntryMessage{} }
func (m *NamespaceEntryMessage) String() string            { return proto.CompactTextString(m) }
func (*NamespaceEntryMessage) ProtoMessage()               {}
func (*NamespaceEntryMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *NamespaceEntryMessage) GetTable() string {
	if m != nil {
		return m.Table
	}
	return ""
}

func (m *NamespaceEntryMessage) GetRow() string {
	if m != nil {
		return m.Row
	}
	return ""
}

func (m *NamespaceEntryMessage) GetEntry() string {
	if m != nil {
		return m.Entry
	}
	return ""
}

func (m *NamespaceEntryMessage) GetPoints() []string {
	if m != nil {
		return m.Points
	}
	return nil
}

type IndexMessage struct {
	Entries []*IndexEntryMessage `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
}

func (m *IndexMessage) Reset()                    { *m = IndexMessage{} }
func (m *IndexMessage) String() string            { return proto.CompactTextString(m) }
func (*IndexMessage) ProtoMessage()               {}
func (*IndexMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *IndexMessage) GetEntries() []*IndexEntryMessage {
	if m != nil {
		return m.Entries
	}
	return nil
}

type IndexEntryMessage struct {
	Table string   `protobuf:"bytes,1,opt,name=table" json:"table,omitempty"`
	Links []string `protobuf:"bytes,2,rep,name=links" json:"links,omitempty"`
}

func (m *IndexEntryMessage) Reset()                    { *m = IndexEntryMessage{} }
func (m *IndexEntryMessage) String() string            { return proto.CompactTextString(m) }
func (*IndexEntryMessage) ProtoMessage()               {}
func (*IndexEntryMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *IndexEntryMessage) GetTable() string {
	if m != nil {
		return m.Table
	}
	return ""
}

func (m *IndexEntryMessage) GetLinks() []string {
	if m != nil {
		return m.Links
	}
	return nil
}

type APIResponseMessage struct {
	Message          string                     `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
	Error            string                     `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
	Type             uint32                     `protobuf:"varint,3,opt,name=type" json:"type,omitempty"`
	QueryResponse    *APIQueryResponseMessage   `protobuf:"bytes,4,opt,name=queryResponse" json:"queryResponse,omitempty"`
	RefelectResponse *APIReflectResponseMessage `protobuf:"bytes,5,opt,name=refelectResponse" json:"refelectResponse,omitempty"`
}

func (m *APIResponseMessage) Reset()                    { *m = APIResponseMessage{} }
func (m *APIResponseMessage) String() string            { return proto.CompactTextString(m) }
func (*APIResponseMessage) ProtoMessage()               {}
func (*APIResponseMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *APIResponseMessage) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *APIResponseMessage) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *APIResponseMessage) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *APIResponseMessage) GetQueryResponse() *APIQueryResponseMessage {
	if m != nil {
		return m.QueryResponse
	}
	return nil
}

func (m *APIResponseMessage) GetRefelectResponse() *APIReflectResponseMessage {
	if m != nil {
		return m.RefelectResponse
	}
	return nil
}

type APIQueryResponseMessage struct {
	Namespace *NamespaceMessage `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
}

func (m *APIQueryResponseMessage) Reset()                    { *m = APIQueryResponseMessage{} }
func (m *APIQueryResponseMessage) String() string            { return proto.CompactTextString(m) }
func (*APIQueryResponseMessage) ProtoMessage()               {}
func (*APIQueryResponseMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *APIQueryResponseMessage) GetNamespace() *NamespaceMessage {
	if m != nil {
		return m.Namespace
	}
	return nil
}

type APIReflectResponseMessage struct {
	Type      uint32            `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	Path      string            `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	Index     *IndexMessage     `protobuf:"bytes,3,opt,name=index" json:"index,omitempty"`
	Namespace *NamespaceMessage `protobuf:"bytes,4,opt,name=namespace" json:"namespace,omitempty"`
}

func (m *APIReflectResponseMessage) Reset()                    { *m = APIReflectResponseMessage{} }
func (m *APIReflectResponseMessage) String() string            { return proto.CompactTextString(m) }
func (*APIReflectResponseMessage) ProtoMessage()               {}
func (*APIReflectResponseMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *APIReflectResponseMessage) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *APIReflectResponseMessage) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *APIReflectResponseMessage) GetIndex() *IndexMessage {
	if m != nil {
		return m.Index
	}
	return nil
}

func (m *APIReflectResponseMessage) GetNamespace() *NamespaceMessage {
	if m != nil {
		return m.Namespace
	}
	return nil
}

type QueryMessage struct {
	OpCode uint32              `protobuf:"varint,1,opt,name=opCode" json:"opCode,omitempty"`
	Table  string              `protobuf:"bytes,2,opt,name=table" json:"table,omitempty"`
	Join   *QueryJoinMessage   `protobuf:"bytes,3,opt,name=join" json:"join,omitempty"`
	Select *QuerySelectMessage `protobuf:"bytes,4,opt,name=select" json:"select,omitempty"`
}

func (m *QueryMessage) Reset()                    { *m = QueryMessage{} }
func (m *QueryMessage) String() string            { return proto.CompactTextString(m) }
func (*QueryMessage) ProtoMessage()               {}
func (*QueryMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *QueryMessage) GetOpCode() uint32 {
	if m != nil {
		return m.OpCode
	}
	return 0
}

func (m *QueryMessage) GetTable() string {
	if m != nil {
		return m.Table
	}
	return ""
}

func (m *QueryMessage) GetJoin() *QueryJoinMessage {
	if m != nil {
		return m.Join
	}
	return nil
}

func (m *QueryMessage) GetSelect() *QuerySelectMessage {
	if m != nil {
		return m.Select
	}
	return nil
}

type QueryJoinMessage struct {
	Rows []*QueryRowJoinMessage `protobuf:"bytes,1,rep,name=rows" json:"rows,omitempty"`
}

func (m *QueryJoinMessage) Reset()                    { *m = QueryJoinMessage{} }
func (m *QueryJoinMessage) String() string            { return proto.CompactTextString(m) }
func (*QueryJoinMessage) ProtoMessage()               {}
func (*QueryJoinMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *QueryJoinMessage) GetRows() []*QueryRowJoinMessage {
	if m != nil {
		return m.Rows
	}
	return nil
}

type QueryRowJoinMessage struct {
	Row     string                      `protobuf:"bytes,1,opt,name=row" json:"row,omitempty"`
	Entries []*QueryRowJoinEntryMessage `protobuf:"bytes,2,rep,name=entries" json:"entries,omitempty"`
}

func (m *QueryRowJoinMessage) Reset()                    { *m = QueryRowJoinMessage{} }
func (m *QueryRowJoinMessage) String() string            { return proto.CompactTextString(m) }
func (*QueryRowJoinMessage) ProtoMessage()               {}
func (*QueryRowJoinMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *QueryRowJoinMessage) GetRow() string {
	if m != nil {
		return m.Row
	}
	return ""
}

func (m *QueryRowJoinMessage) GetEntries() []*QueryRowJoinEntryMessage {
	if m != nil {
		return m.Entries
	}
	return nil
}

type QueryRowJoinEntryMessage struct {
	Entry string `protobuf:"bytes,1,opt,name=entry" json:"entry,omitempty"`
	Point string `protobuf:"bytes,2,opt,name=point" json:"point,omitempty"`
}

func (m *QueryRowJoinEntryMessage) Reset()                    { *m = QueryRowJoinEntryMessage{} }
func (m *QueryRowJoinEntryMessage) String() string            { return proto.CompactTextString(m) }
func (*QueryRowJoinEntryMessage) ProtoMessage()               {}
func (*QueryRowJoinEntryMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *QueryRowJoinEntryMessage) GetEntry() string {
	if m != nil {
		return m.Entry
	}
	return ""
}

func (m *QueryRowJoinEntryMessage) GetPoint() string {
	if m != nil {
		return m.Point
	}
	return ""
}

type QuerySelectMessage struct {
	Limit uint32             `protobuf:"varint,1,opt,name=limit" json:"limit,omitempty"`
	Where *QueryWhereMessage `protobuf:"bytes,2,opt,name=where" json:"where,omitempty"`
}

func (m *QuerySelectMessage) Reset()                    { *m = QuerySelectMessage{} }
func (m *QuerySelectMessage) String() string            { return proto.CompactTextString(m) }
func (*QuerySelectMessage) ProtoMessage()               {}
func (*QuerySelectMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *QuerySelectMessage) GetLimit() uint32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *QuerySelectMessage) GetWhere() *QueryWhereMessage {
	if m != nil {
		return m.Where
	}
	return nil
}

type QueryWhereMessage struct {
	OpCode    uint32                 `protobuf:"varint,1,opt,name=opCode" json:"opCode,omitempty"`
	Predicate *QueryPredicateMessage `protobuf:"bytes,2,opt,name=predicate" json:"predicate,omitempty"`
	Clauses   []*QueryWhereMessage   `protobuf:"bytes,3,rep,name=clauses" json:"clauses,omitempty"`
}

func (m *QueryWhereMessage) Reset()                    { *m = QueryWhereMessage{} }
func (m *QueryWhereMessage) String() string            { return proto.CompactTextString(m) }
func (*QueryWhereMessage) ProtoMessage()               {}
func (*QueryWhereMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *QueryWhereMessage) GetOpCode() uint32 {
	if m != nil {
		return m.OpCode
	}
	return 0
}

func (m *QueryWhereMessage) GetPredicate() *QueryPredicateMessage {
	if m != nil {
		return m.Predicate
	}
	return nil
}

func (m *QueryWhereMessage) GetClauses() []*QueryWhereMessage {
	if m != nil {
		return m.Clauses
	}
	return nil
}

type QueryPredicateMessage struct {
	OpCode   uint32   `protobuf:"varint,1,opt,name=opCode" json:"opCode,omitempty"`
	Keys     []string `protobuf:"bytes,2,rep,name=keys" json:"keys,omitempty"`
	Literals []string `protobuf:"bytes,3,rep,name=literals" json:"literals,omitempty"`
	Userow   bool     `protobuf:"varint,4,opt,name=userow" json:"userow,omitempty"`
}

func (m *QueryPredicateMessage) Reset()                    { *m = QueryPredicateMessage{} }
func (m *QueryPredicateMessage) String() string            { return proto.CompactTextString(m) }
func (*QueryPredicateMessage) ProtoMessage()               {}
func (*QueryPredicateMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *QueryPredicateMessage) GetOpCode() uint32 {
	if m != nil {
		return m.OpCode
	}
	return 0
}

func (m *QueryPredicateMessage) GetKeys() []string {
	if m != nil {
		return m.Keys
	}
	return nil
}

func (m *QueryPredicateMessage) GetLiterals() []string {
	if m != nil {
		return m.Literals
	}
	return nil
}

func (m *QueryPredicateMessage) GetUserow() bool {
	if m != nil {
		return m.Userow
	}
	return false
}

func init() {
	proto.RegisterType((*NamespaceMessage)(nil), "godless.NamespaceMessage")
	proto.RegisterType((*NamespaceEntryMessage)(nil), "godless.NamespaceEntryMessage")
	proto.RegisterType((*IndexMessage)(nil), "godless.IndexMessage")
	proto.RegisterType((*IndexEntryMessage)(nil), "godless.IndexEntryMessage")
	proto.RegisterType((*APIResponseMessage)(nil), "godless.APIResponseMessage")
	proto.RegisterType((*APIQueryResponseMessage)(nil), "godless.APIQueryResponseMessage")
	proto.RegisterType((*APIReflectResponseMessage)(nil), "godless.APIReflectResponseMessage")
	proto.RegisterType((*QueryMessage)(nil), "godless.QueryMessage")
	proto.RegisterType((*QueryJoinMessage)(nil), "godless.QueryJoinMessage")
	proto.RegisterType((*QueryRowJoinMessage)(nil), "godless.QueryRowJoinMessage")
	proto.RegisterType((*QueryRowJoinEntryMessage)(nil), "godless.QueryRowJoinEntryMessage")
	proto.RegisterType((*QuerySelectMessage)(nil), "godless.QuerySelectMessage")
	proto.RegisterType((*QueryWhereMessage)(nil), "godless.QueryWhereMessage")
	proto.RegisterType((*QueryPredicateMessage)(nil), "godless.QueryPredicateMessage")
}

func init() { proto.RegisterFile("godless.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 626 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0xcf, 0x6e, 0xd3, 0x4e,
	0x10, 0x96, 0x13, 0x3b, 0x6d, 0x26, 0xad, 0x94, 0xee, 0x2f, 0xe9, 0xcf, 0x2d, 0xa8, 0x0a, 0x7b,
	0x8a, 0x84, 0xa8, 0xaa, 0x14, 0x09, 0x24, 0x90, 0x50, 0x45, 0xa9, 0x14, 0x04, 0x55, 0x30, 0x07,
	0x2e, 0x5c, 0xdc, 0x64, 0xdb, 0x9a, 0x3a, 0x5e, 0xb3, 0xbb, 0x51, 0xc8, 0xd3, 0xc0, 0x95, 0xb7,
	0xe3, 0x11, 0x90, 0xc7, 0xbb, 0xfe, 0x13, 0xdb, 0x07, 0x6e, 0x3b, 0xde, 0x6f, 0xbe, 0xf9, 0x66,
	0x76, 0x66, 0x0c, 0xfb, 0x77, 0x7c, 0x11, 0x32, 0x29, 0x4f, 0x63, 0xc1, 0x15, 0x27, 0x3b, 0xda,
	0xa4, 0x1f, 0xa0, 0x7f, 0xed, 0x2f, 0x99, 0x8c, 0xfd, 0x39, 0xfb, 0xc8, 0xa4, 0xf4, 0xef, 0x18,
	0x79, 0x09, 0x3b, 0x2c, 0x52, 0x22, 0x60, 0xd2, 0xb5, 0x46, 0xed, 0x71, 0x6f, 0x72, 0x72, 0x6a,
	0xbc, 0x33, 0xec, 0xbb, 0x48, 0x89, 0x8d, 0x76, 0xf0, 0x0c, 0x9c, 0x2e, 0x61, 0x58, 0x8b, 0x20,
	0x03, 0x70, 0x94, 0x7f, 0x13, 0x32, 0xd7, 0x1a, 0x59, 0xe3, 0xae, 0x97, 0x1a, 0xa4, 0x0f, 0x6d,
	0xc1, 0xd7, 0x6e, 0x0b, 0xbf, 0x25, 0xc7, 0x04, 0x97, 0x70, 0x6d, 0xdc, 0x76, 0x8a, 0x43, 0x83,
	0x1c, 0x42, 0x27, 0xe6, 0x41, 0xa4, 0xa4, 0x6b, 0x8f, 0xda, 0xe3, 0xae, 0xa7, 0x2d, 0x7a, 0x09,
	0x7b, 0xd3, 0x68, 0xc1, 0x7e, 0x98, 0x28, 0xcf, 0xb7, 0x85, 0x1f, 0x67, 0xc2, 0x11, 0x57, 0x2f,
	0xfa, 0x0d, 0x1c, 0x54, 0x6e, 0x1b, 0x04, 0x0f, 0xc0, 0x09, 0x83, 0xe8, 0x41, 0xba, 0x2d, 0xd4,
	0x91, 0x1a, 0xf4, 0x8f, 0x05, 0xe4, 0x62, 0x36, 0xf5, 0x98, 0x8c, 0x79, 0x24, 0xb3, 0x32, 0xba,
	0xb0, 0xb3, 0x4c, 0x8f, 0x9a, 0xc4, 0x98, 0x98, 0xa5, 0x10, 0x5c, 0xe8, 0xcc, 0x53, 0x83, 0x10,
	0xb0, 0xd5, 0x26, 0x66, 0x98, 0xfa, 0xbe, 0x87, 0x67, 0x72, 0x05, 0xfb, 0xdf, 0x57, 0x4c, 0x6c,
	0x0c, 0xb7, 0x6b, 0x8f, 0xac, 0x71, 0x6f, 0x32, 0xca, 0xf2, 0xba, 0x98, 0x4d, 0x3f, 0x15, 0x01,
	0x26, 0xbb, 0xb2, 0x1b, 0xb9, 0x86, 0xbe, 0x60, 0xb7, 0x2c, 0x64, 0x73, 0x95, 0x51, 0x39, 0x48,
	0x45, 0x8b, 0x54, 0x1e, 0xbb, 0x2d, 0x42, 0x0c, 0x59, 0xc5, 0x97, 0x7a, 0xf0, 0x7f, 0x43, 0x64,
	0xf2, 0x02, 0xba, 0x91, 0xe9, 0x01, 0x4c, 0xbc, 0x37, 0x39, 0xaa, 0xf6, 0x8f, 0xa1, 0xce, 0xb1,
	0xf4, 0xb7, 0x05, 0x47, 0x8d, 0x1a, 0xb2, 0xea, 0x58, 0x85, 0xea, 0x10, 0xb0, 0x63, 0x5f, 0xdd,
	0xeb, 0x32, 0xe2, 0x99, 0x3c, 0x05, 0x27, 0x48, 0x5e, 0x13, 0xcb, 0xd8, 0x9b, 0x0c, 0xcb, 0x1d,
	0x60, 0xc2, 0xa6, 0x98, 0xb2, 0x56, 0xfb, 0x1f, 0xb4, 0xfe, 0xb2, 0x60, 0x0f, 0xb3, 0x37, 0xf2,
	0x0e, 0xa1, 0xc3, 0xe3, 0xb7, 0x7c, 0x61, 0x04, 0x6a, 0x2b, 0xef, 0xa3, 0x56, 0xb1, 0x8f, 0x9e,
	0x81, 0xfd, 0x8d, 0x07, 0x91, 0xd6, 0x98, 0x87, 0x44, 0xca, 0xf7, 0x3c, 0x88, 0x4c, 0x48, 0x84,
	0x91, 0x73, 0xe8, 0x48, 0xac, 0xbf, 0xd6, 0xf8, 0xa8, 0xec, 0xf0, 0x19, 0xef, 0x8c, 0x8b, 0x86,
	0xd2, 0x4b, 0xe8, 0x6f, 0xd3, 0x91, 0x33, 0xb0, 0x05, 0x5f, 0x9b, 0xe9, 0x78, 0x5c, 0xa6, 0xf1,
	0xf8, 0xba, 0x14, 0x3a, 0x41, 0xd2, 0x05, 0xfc, 0x57, 0x73, 0x69, 0x26, 0xd7, 0xca, 0x27, 0xf7,
	0x55, 0x3e, 0x7b, 0x2d, 0x64, 0x7f, 0x52, 0xcb, 0x5e, 0x3f, 0x82, 0x57, 0xe0, 0x36, 0x81, 0xf2,
	0x95, 0x60, 0x15, 0x57, 0xc2, 0x00, 0x1c, 0x5c, 0x02, 0xa6, 0xae, 0x68, 0xd0, 0xaf, 0x40, 0xaa,
	0x15, 0x49, 0xa7, 0x76, 0x19, 0x28, 0xfd, 0x34, 0xa9, 0x41, 0xce, 0xc0, 0x59, 0xdf, 0x33, 0x91,
	0xbe, 0x4c, 0x71, 0x55, 0x20, 0xc3, 0x97, 0xe4, 0x2a, 0xeb, 0x16, 0x04, 0xd2, 0x9f, 0x16, 0x1c,
	0x54, 0x2e, 0x1b, 0x5f, 0xfe, 0x35, 0x74, 0x63, 0xc1, 0x16, 0xc1, 0xdc, 0x57, 0x26, 0xc6, 0x49,
	0x39, 0xc6, 0xcc, 0x5c, 0x67, 0x0d, 0x96, 0x39, 0x24, 0xab, 0x6c, 0x1e, 0xfa, 0x2b, 0xc9, 0xa4,
	0xdb, 0xde, 0x5a, 0x65, 0x55, 0x7d, 0x06, 0x4a, 0xd7, 0x30, 0xac, 0x65, 0x6e, 0x14, 0x49, 0xc0,
	0x7e, 0x60, 0x1b, 0xb3, 0xcf, 0xf0, 0x4c, 0x8e, 0x61, 0x37, 0x0c, 0x14, 0x13, 0x7e, 0x98, 0xc6,
	0xee, 0x7a, 0x99, 0x9d, 0xf0, 0xac, 0x24, 0x4b, 0x9e, 0x3e, 0xe9, 0xc4, 0x5d, 0x4f, 0x5b, 0x37,
	0x1d, 0xfc, 0xad, 0x9c, 0xff, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x60, 0xb9, 0xb4, 0x2b, 0x67, 0x06,
	0x00, 0x00,
}
