// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: webpush.proto

package msgs

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// A DiscoveryRequest requests a set of versioned resources of the same type for
// a given Envoy node on some API.
type DiscoveryRequest struct {
	// The version_info provided in the request messages will be the version_info
	// received with the most recent successfully processed response or empty on
	// the first request. It is expected that no new request is sent after a
	// response is received until the Envoy instance is ready to ACK/NACK the new
	// configuration. ACK/NACK takes place by returning the new API config version
	// as applied or the previous API config version respectively. Each type_url
	// (see below) has an independent version associated with it.
	VersionInfo string `protobuf:"bytes,1,opt,name=version_info,json=versionInfo,proto3" json:"version_info,omitempty"`
	// List of resources to subscribe to, e.g. list of cluster names or a route
	// configuration name. If this is empty, all resources for the API are
	// returned. LDS/CDS expect empty resource_names, since this is global
	// discovery for the Envoy instance. The LDS and CDS responses will then imply
	// a number of resources that need to be fetched via EDS/RDS, which will be
	// explicitly enumerated in resource_names.
	ResourceNames []string `protobuf:"bytes,3,rep,name=resource_names,json=resourceNames,proto3" json:"resource_names,omitempty"`
	// Type of the resource that is being requested, e.g.
	// "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment". This is implicit
	// in requests made via singleton xDS APIs such as CDS, LDS, etc. but is
	// required for ADS.
	TypeUrl string `protobuf:"bytes,4,opt,name=type_url,json=typeUrl,proto3" json:"type_url,omitempty"`
	// nonce corresponding to DiscoveryResponse being ACK/NACKed. See above
	// discussion on version_info and the DiscoveryResponse nonce comment. This
	// may be empty if no nonce is available, e.g. at startup or for non-stream
	// xDS implementations.
	ResponseNonce string `protobuf:"bytes,5,opt,name=response_nonce,json=responseNonce,proto3" json:"response_nonce,omitempty"`
	// The response resources. These resources are typed and depend on the API being called.
	// google.protobuf.Any
	Resources            []*Any   `protobuf:"bytes,7,rep,name=resources,proto3" json:"resources,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DiscoveryRequest) Reset()         { *m = DiscoveryRequest{} }
func (m *DiscoveryRequest) String() string { return proto.CompactTextString(m) }
func (*DiscoveryRequest) ProtoMessage()    {}
func (*DiscoveryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{0}
}
func (m *DiscoveryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiscoveryRequest.Unmarshal(m, b)
}
func (m *DiscoveryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiscoveryRequest.Marshal(b, m, deterministic)
}
func (m *DiscoveryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiscoveryRequest.Merge(m, src)
}
func (m *DiscoveryRequest) XXX_Size() int {
	return xxx_messageInfo_DiscoveryRequest.Size(m)
}
func (m *DiscoveryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DiscoveryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DiscoveryRequest proto.InternalMessageInfo

func (m *DiscoveryRequest) GetVersionInfo() string {
	if m != nil {
		return m.VersionInfo
	}
	return ""
}

func (m *DiscoveryRequest) GetResourceNames() []string {
	if m != nil {
		return m.ResourceNames
	}
	return nil
}

func (m *DiscoveryRequest) GetTypeUrl() string {
	if m != nil {
		return m.TypeUrl
	}
	return ""
}

func (m *DiscoveryRequest) GetResponseNonce() string {
	if m != nil {
		return m.ResponseNonce
	}
	return ""
}

func (m *DiscoveryRequest) GetResources() []*Any {
	if m != nil {
		return m.Resources
	}
	return nil
}

type DiscoveryResponse struct {
	// The version of the response data.
	VersionInfo string `protobuf:"bytes,1,opt,name=version_info,json=versionInfo,proto3" json:"version_info,omitempty"`
	// The response resources. These resources are typed and depend on the API being called.
	// google.protobuf.Any
	Resources []*Any `protobuf:"bytes,2,rep,name=resources,proto3" json:"resources,omitempty"`
	// Type URL for resources. This must be consistent with the type_url in the
	// Any messages for resources if resources is non-empty. This effectively
	// identifies the xDS API when muxing over ADS.
	TypeUrl string `protobuf:"bytes,4,opt,name=type_url,json=typeUrl,proto3" json:"type_url,omitempty"`
	// For gRPC based subscriptions, the nonce provides a way to explicitly ack a
	// specific DiscoveryResponse in a following DiscoveryRequest. Additional
	// messages may have been sent by Envoy to the management server for the
	// previous version on the stream prior to this DiscoveryResponse, that were
	// unprocessed at response send time. The nonce allows the management server
	// to ignore any further DiscoveryRequests for the previous version until a
	// DiscoveryRequest bearing the nonce. The nonce is optional and is not
	// required for non-stream based xDS implementations.
	Nonce                string   `protobuf:"bytes,5,opt,name=nonce,proto3" json:"nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DiscoveryResponse) Reset()         { *m = DiscoveryResponse{} }
func (m *DiscoveryResponse) String() string { return proto.CompactTextString(m) }
func (*DiscoveryResponse) ProtoMessage()    {}
func (*DiscoveryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{1}
}
func (m *DiscoveryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiscoveryResponse.Unmarshal(m, b)
}
func (m *DiscoveryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiscoveryResponse.Marshal(b, m, deterministic)
}
func (m *DiscoveryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiscoveryResponse.Merge(m, src)
}
func (m *DiscoveryResponse) XXX_Size() int {
	return xxx_messageInfo_DiscoveryResponse.Size(m)
}
func (m *DiscoveryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DiscoveryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DiscoveryResponse proto.InternalMessageInfo

func (m *DiscoveryResponse) GetVersionInfo() string {
	if m != nil {
		return m.VersionInfo
	}
	return ""
}

func (m *DiscoveryResponse) GetResources() []*Any {
	if m != nil {
		return m.Resources
	}
	return nil
}

func (m *DiscoveryResponse) GetTypeUrl() string {
	if m != nil {
		return m.TypeUrl
	}
	return ""
}

func (m *DiscoveryResponse) GetNonce() string {
	if m != nil {
		return m.Nonce
	}
	return ""
}

// Copied to avoid dependency. type_url format: prefix/message.type
type Any struct {
	TypeUrl              string   `protobuf:"bytes,1,opt,name=type_url,json=typeUrl,proto3" json:"type_url,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Any) Reset()         { *m = Any{} }
func (m *Any) String() string { return proto.CompactTextString(m) }
func (*Any) ProtoMessage()    {}
func (*Any) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{2}
}
func (m *Any) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Any.Unmarshal(m, b)
}
func (m *Any) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Any.Marshal(b, m, deterministic)
}
func (m *Any) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Any.Merge(m, src)
}
func (m *Any) XXX_Size() int {
	return xxx_messageInfo_Any.Size(m)
}
func (m *Any) XXX_DiscardUnknown() {
	xxx_messageInfo_Any.DiscardUnknown(m)
}

var xxx_messageInfo_Any proto.InternalMessageInfo

func (m *Any) GetTypeUrl() string {
	if m != nil {
		return m.TypeUrl
	}
	return ""
}

func (m *Any) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

// Message is returned as PUSH PROMISE frames in the spec. The alternative protocol wraps it in
// Any field or other framing.
type MessageEnvelope struct {
	MessageId string `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	// Maps to the SubscribeResponse push parameter, returned as Link rel="urn:ietf:params:push"
	// in the push promise.
	Push string `protobuf:"bytes,2,opt,name=push,proto3" json:"push,omitempty"`
	// If 'dh' and 'salt' are set, will contain encrypted data.
	// Otherwise it is a plaintext message.
	Data []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	// Identifies the sender.
	Sender               *Vapid   `protobuf:"bytes,4,opt,name=sender,proto3" json:"sender,omitempty"`
	ContentEncoding      string   `protobuf:"bytes,7,opt,name=content_encoding,json=contentEncoding,proto3" json:"content_encoding,omitempty"`
	Salt                 []byte   `protobuf:"bytes,8,opt,name=salt,proto3" json:"salt,omitempty"`
	Dh                   []byte   `protobuf:"bytes,9,opt,name=dh,proto3" json:"dh,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MessageEnvelope) Reset()         { *m = MessageEnvelope{} }
func (m *MessageEnvelope) String() string { return proto.CompactTextString(m) }
func (*MessageEnvelope) ProtoMessage()    {}
func (*MessageEnvelope) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{3}
}
func (m *MessageEnvelope) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageEnvelope.Unmarshal(m, b)
}
func (m *MessageEnvelope) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageEnvelope.Marshal(b, m, deterministic)
}
func (m *MessageEnvelope) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageEnvelope.Merge(m, src)
}
func (m *MessageEnvelope) XXX_Size() int {
	return xxx_messageInfo_MessageEnvelope.Size(m)
}
func (m *MessageEnvelope) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageEnvelope.DiscardUnknown(m)
}

var xxx_messageInfo_MessageEnvelope proto.InternalMessageInfo

func (m *MessageEnvelope) GetMessageId() string {
	if m != nil {
		return m.MessageId
	}
	return ""
}

func (m *MessageEnvelope) GetPush() string {
	if m != nil {
		return m.Push
	}
	return ""
}

func (m *MessageEnvelope) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *MessageEnvelope) GetSender() *Vapid {
	if m != nil {
		return m.Sender
	}
	return nil
}

func (m *MessageEnvelope) GetContentEncoding() string {
	if m != nil {
		return m.ContentEncoding
	}
	return ""
}

func (m *MessageEnvelope) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *MessageEnvelope) GetDh() []byte {
	if m != nil {
		return m.Dh
	}
	return nil
}

// Vapid is the proto variant of a Webpush JWT.
//
// For HTTP, included in Authorization header:
// Authorization: vapid t=B64url k=B64url
//
// Decoded t is of form: { "typ": "JWT", "alg": "ES256" }.JWT.SIG
//
// { "crv":"P-256",
//   "kty":"EC",
//   "x":"DUfHPKLVFQzVvnCPGyfucbECzPDa7rWbXriLcysAjEc",
//   "y":"F6YK5h4SDYic-dRuU_RCPCfA5aq9ojSwk5Y2EmClBPs" }
type Vapid struct {
	// json payload of VAPID ( without base64 encoding)
	// Can also be a proto message when used over other transports.
	// Verification requires converting back to base64 !
	// We decode to reduce the binary size
	Data []byte `protobuf:"bytes,7,opt,name=data,proto3" json:"data,omitempty"`
	// Public key of the signer, 64 bytes, EC256.
	// Included in 'k' parameter for HTTP.
	K []byte `protobuf:"bytes,4,opt,name=k,proto3" json:"k,omitempty"`
	// If empty, it is assumed to be the constant value {typ=JWT,alg=ES256}
	TType                []byte   `protobuf:"bytes,32,opt,name=t_type,json=tType,proto3" json:"t_type,omitempty"`
	TSignature           []byte   `protobuf:"bytes,33,opt,name=t_signature,json=tSignature,proto3" json:"t_signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Vapid) Reset()         { *m = Vapid{} }
func (m *Vapid) String() string { return proto.CompactTextString(m) }
func (*Vapid) ProtoMessage()    {}
func (*Vapid) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{4}
}
func (m *Vapid) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Vapid.Unmarshal(m, b)
}
func (m *Vapid) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Vapid.Marshal(b, m, deterministic)
}
func (m *Vapid) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vapid.Merge(m, src)
}
func (m *Vapid) XXX_Size() int {
	return xxx_messageInfo_Vapid.Size(m)
}
func (m *Vapid) XXX_DiscardUnknown() {
	xxx_messageInfo_Vapid.DiscardUnknown(m)
}

var xxx_messageInfo_Vapid proto.InternalMessageInfo

func (m *Vapid) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Vapid) GetK() []byte {
	if m != nil {
		return m.K
	}
	return nil
}

func (m *Vapid) GetTType() []byte {
	if m != nil {
		return m.TType
	}
	return nil
}

func (m *Vapid) GetTSignature() []byte {
	if m != nil {
		return m.TSignature
	}
	return nil
}

//
type SubscribeRequest struct {
	// A UA should group subscriptions in a set. First request from a
	// UA will not include a set - it is typically a subscription associated with
	// the UA itself.
	PushSet string `protobuf:"bytes,1,opt,name=push_set,json=pushSet,proto3" json:"push_set,omitempty"`
	// Included as Crypto-Key: p256ecdsa parameter.
	// Corresponds to the applicationServerKey parameter in the PushSubscriptionOptions in
	// the W3C API
	SenderVapid          string   `protobuf:"bytes,2,opt,name=sender_vapid,json=senderVapid,proto3" json:"sender_vapid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubscribeRequest) Reset()         { *m = SubscribeRequest{} }
func (m *SubscribeRequest) String() string { return proto.CompactTextString(m) }
func (*SubscribeRequest) ProtoMessage()    {}
func (*SubscribeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{5}
}
func (m *SubscribeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscribeRequest.Unmarshal(m, b)
}
func (m *SubscribeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscribeRequest.Marshal(b, m, deterministic)
}
func (m *SubscribeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscribeRequest.Merge(m, src)
}
func (m *SubscribeRequest) XXX_Size() int {
	return xxx_messageInfo_SubscribeRequest.Size(m)
}
func (m *SubscribeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscribeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubscribeRequest proto.InternalMessageInfo

func (m *SubscribeRequest) GetPushSet() string {
	if m != nil {
		return m.PushSet
	}
	return ""
}

func (m *SubscribeRequest) GetSenderVapid() string {
	if m != nil {
		return m.SenderVapid
	}
	return ""
}

// Subscribe response includes the elements in the spec.
type SubscribeResponse struct {
	// Returned as Link: rel="urn:ietf:params:push"
	// Spec examples use a full path ( /push/xxxx1 )
	// TODO: clarify if it can be a full URL
	Push string `protobuf:"bytes,1,opt,name=push,proto3" json:"push,omitempty"`
	// Optional response: it
	// returned as Link: rel=urn:ietf:params:push:set
	// Spec examples use a full path ( /subscription-set/xxxx2 ).
	// TODO: clarify it can be a full URL, like subscription
	PushSet string `protobuf:"bytes,2,opt,name=push_set,json=pushSet,proto3" json:"push_set,omitempty"`
	// Push subscription resource. This is the full URL where the UA will use to
	// receive the messages, using the PUSH promise http2 frame.
	//
	//
	// Returned as Location header in the spec
	Location             string   `protobuf:"bytes,3,opt,name=location,proto3" json:"location,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubscribeResponse) Reset()         { *m = SubscribeResponse{} }
func (m *SubscribeResponse) String() string { return proto.CompactTextString(m) }
func (*SubscribeResponse) ProtoMessage()    {}
func (*SubscribeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{6}
}
func (m *SubscribeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscribeResponse.Unmarshal(m, b)
}
func (m *SubscribeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscribeResponse.Marshal(b, m, deterministic)
}
func (m *SubscribeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscribeResponse.Merge(m, src)
}
func (m *SubscribeResponse) XXX_Size() int {
	return xxx_messageInfo_SubscribeResponse.Size(m)
}
func (m *SubscribeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscribeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SubscribeResponse proto.InternalMessageInfo

func (m *SubscribeResponse) GetPush() string {
	if m != nil {
		return m.Push
	}
	return ""
}

func (m *SubscribeResponse) GetPushSet() string {
	if m != nil {
		return m.PushSet
	}
	return ""
}

func (m *SubscribeResponse) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

type PushRequest struct {
	// The value returned in the SubscribeResponse push, without the hostname.
	Push    string `protobuf:"bytes,1,opt,name=push,proto3" json:"push,omitempty"`
	Ttl     int32  `protobuf:"varint,2,opt,name=ttl,proto3" json:"ttl,omitempty"`
	Data    []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Urgency string `protobuf:"bytes,4,opt,name=urgency,proto3" json:"urgency,omitempty"`
	// Prefer header indicating delivery receipt request.
	RespondAsync         bool     `protobuf:"varint,5,opt,name=respond_async,json=respondAsync,proto3" json:"respond_async,omitempty"`
	Topic                string   `protobuf:"bytes,6,opt,name=topic,proto3" json:"topic,omitempty"`
	ContentEncoding      string   `protobuf:"bytes,7,opt,name=content_encoding,json=contentEncoding,proto3" json:"content_encoding,omitempty"`
	Salt                 string   `protobuf:"bytes,8,opt,name=salt,proto3" json:"salt,omitempty"`
	Dh                   string   `protobuf:"bytes,9,opt,name=dh,proto3" json:"dh,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushRequest) Reset()         { *m = PushRequest{} }
func (m *PushRequest) String() string { return proto.CompactTextString(m) }
func (*PushRequest) ProtoMessage()    {}
func (*PushRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{7}
}
func (m *PushRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushRequest.Unmarshal(m, b)
}
func (m *PushRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushRequest.Marshal(b, m, deterministic)
}
func (m *PushRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushRequest.Merge(m, src)
}
func (m *PushRequest) XXX_Size() int {
	return xxx_messageInfo_PushRequest.Size(m)
}
func (m *PushRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PushRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PushRequest proto.InternalMessageInfo

func (m *PushRequest) GetPush() string {
	if m != nil {
		return m.Push
	}
	return ""
}

func (m *PushRequest) GetTtl() int32 {
	if m != nil {
		return m.Ttl
	}
	return 0
}

func (m *PushRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *PushRequest) GetUrgency() string {
	if m != nil {
		return m.Urgency
	}
	return ""
}

func (m *PushRequest) GetRespondAsync() bool {
	if m != nil {
		return m.RespondAsync
	}
	return false
}

func (m *PushRequest) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *PushRequest) GetContentEncoding() string {
	if m != nil {
		return m.ContentEncoding
	}
	return ""
}

func (m *PushRequest) GetSalt() string {
	if m != nil {
		return m.Salt
	}
	return ""
}

func (m *PushRequest) GetDh() string {
	if m != nil {
		return m.Dh
	}
	return ""
}

type PushResponse struct {
	MessageId string `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	// If request includes the respond_async parameter.
	//
	PushReceipt          string   `protobuf:"bytes,2,opt,name=push_receipt,json=pushReceipt,proto3" json:"push_receipt,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushResponse) Reset()         { *m = PushResponse{} }
func (m *PushResponse) String() string { return proto.CompactTextString(m) }
func (*PushResponse) ProtoMessage()    {}
func (*PushResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{8}
}
func (m *PushResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushResponse.Unmarshal(m, b)
}
func (m *PushResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushResponse.Marshal(b, m, deterministic)
}
func (m *PushResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushResponse.Merge(m, src)
}
func (m *PushResponse) XXX_Size() int {
	return xxx_messageInfo_PushResponse.Size(m)
}
func (m *PushResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PushResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PushResponse proto.InternalMessageInfo

func (m *PushResponse) GetMessageId() string {
	if m != nil {
		return m.MessageId
	}
	return ""
}

func (m *PushResponse) GetPushReceipt() string {
	if m != nil {
		return m.PushReceipt
	}
	return ""
}

type MonitorRequest struct {
	// This is the push or push_set in the subscribe response.
	PushSet string `protobuf:"bytes,1,opt,name=push_set,json=pushSet,proto3" json:"push_set,omitempty"`
	// JWT token, signed with key
	Authorization string `protobuf:"bytes,2,opt,name=authorization,proto3" json:"authorization,omitempty"`
	// Public key used for signing, identifies sender/receiver
	Key                  string   `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MonitorRequest) Reset()         { *m = MonitorRequest{} }
func (m *MonitorRequest) String() string { return proto.CompactTextString(m) }
func (*MonitorRequest) ProtoMessage()    {}
func (*MonitorRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{9}
}
func (m *MonitorRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorRequest.Unmarshal(m, b)
}
func (m *MonitorRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorRequest.Marshal(b, m, deterministic)
}
func (m *MonitorRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorRequest.Merge(m, src)
}
func (m *MonitorRequest) XXX_Size() int {
	return xxx_messageInfo_MonitorRequest.Size(m)
}
func (m *MonitorRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorRequest proto.InternalMessageInfo

func (m *MonitorRequest) GetPushSet() string {
	if m != nil {
		return m.PushSet
	}
	return ""
}

func (m *MonitorRequest) GetAuthorization() string {
	if m != nil {
		return m.Authorization
	}
	return ""
}

func (m *MonitorRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type AckRequest struct {
	MessageId            string   `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AckRequest) Reset()         { *m = AckRequest{} }
func (m *AckRequest) String() string { return proto.CompactTextString(m) }
func (*AckRequest) ProtoMessage()    {}
func (*AckRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{10}
}
func (m *AckRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AckRequest.Unmarshal(m, b)
}
func (m *AckRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AckRequest.Marshal(b, m, deterministic)
}
func (m *AckRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AckRequest.Merge(m, src)
}
func (m *AckRequest) XXX_Size() int {
	return xxx_messageInfo_AckRequest.Size(m)
}
func (m *AckRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AckRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AckRequest proto.InternalMessageInfo

func (m *AckRequest) GetMessageId() string {
	if m != nil {
		return m.MessageId
	}
	return ""
}

type AckResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AckResponse) Reset()         { *m = AckResponse{} }
func (m *AckResponse) String() string { return proto.CompactTextString(m) }
func (*AckResponse) ProtoMessage()    {}
func (*AckResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{11}
}
func (m *AckResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AckResponse.Unmarshal(m, b)
}
func (m *AckResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AckResponse.Marshal(b, m, deterministic)
}
func (m *AckResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AckResponse.Merge(m, src)
}
func (m *AckResponse) XXX_Size() int {
	return xxx_messageInfo_AckResponse.Size(m)
}
func (m *AckResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AckResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AckResponse proto.InternalMessageInfo

type ReceiptRequest struct {
	ReceiptSubscription  string   `protobuf:"bytes,1,opt,name=receipt_subscription,json=receiptSubscription,proto3" json:"receipt_subscription,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReceiptRequest) Reset()         { *m = ReceiptRequest{} }
func (m *ReceiptRequest) String() string { return proto.CompactTextString(m) }
func (*ReceiptRequest) ProtoMessage()    {}
func (*ReceiptRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{12}
}
func (m *ReceiptRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReceiptRequest.Unmarshal(m, b)
}
func (m *ReceiptRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReceiptRequest.Marshal(b, m, deterministic)
}
func (m *ReceiptRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReceiptRequest.Merge(m, src)
}
func (m *ReceiptRequest) XXX_Size() int {
	return xxx_messageInfo_ReceiptRequest.Size(m)
}
func (m *ReceiptRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReceiptRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReceiptRequest proto.InternalMessageInfo

func (m *ReceiptRequest) GetReceiptSubscription() string {
	if m != nil {
		return m.ReceiptSubscription
	}
	return ""
}

type Receipt struct {
	MessageId            string   `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Receipt) Reset()         { *m = Receipt{} }
func (m *Receipt) String() string { return proto.CompactTextString(m) }
func (*Receipt) ProtoMessage()    {}
func (*Receipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_74c9ab69d6636d6b, []int{13}
}
func (m *Receipt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Receipt.Unmarshal(m, b)
}
func (m *Receipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Receipt.Marshal(b, m, deterministic)
}
func (m *Receipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Receipt.Merge(m, src)
}
func (m *Receipt) XXX_Size() int {
	return xxx_messageInfo_Receipt.Size(m)
}
func (m *Receipt) XXX_DiscardUnknown() {
	xxx_messageInfo_Receipt.DiscardUnknown(m)
}

var xxx_messageInfo_Receipt proto.InternalMessageInfo

func (m *Receipt) GetMessageId() string {
	if m != nil {
		return m.MessageId
	}
	return ""
}

func init() {
	proto.RegisterType((*DiscoveryRequest)(nil), "webpush.DiscoveryRequest")
	proto.RegisterType((*DiscoveryResponse)(nil), "webpush.DiscoveryResponse")
	proto.RegisterType((*Any)(nil), "webpush.Any")
	proto.RegisterType((*MessageEnvelope)(nil), "webpush.MessageEnvelope")
	proto.RegisterType((*Vapid)(nil), "webpush.Vapid")
	proto.RegisterType((*SubscribeRequest)(nil), "webpush.SubscribeRequest")
	proto.RegisterType((*SubscribeResponse)(nil), "webpush.SubscribeResponse")
	proto.RegisterType((*PushRequest)(nil), "webpush.PushRequest")
	proto.RegisterType((*PushResponse)(nil), "webpush.PushResponse")
	proto.RegisterType((*MonitorRequest)(nil), "webpush.MonitorRequest")
	proto.RegisterType((*AckRequest)(nil), "webpush.AckRequest")
	proto.RegisterType((*AckResponse)(nil), "webpush.AckResponse")
	proto.RegisterType((*ReceiptRequest)(nil), "webpush.ReceiptRequest")
	proto.RegisterType((*Receipt)(nil), "webpush.Receipt")
}

func init() { proto.RegisterFile("webpush.proto", fileDescriptor_74c9ab69d6636d6b) }

var fileDescriptor_74c9ab69d6636d6b = []byte{
	// 706 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0x4d, 0x6f, 0xd3, 0x4a,
	0x14, 0x95, 0x93, 0xe6, 0xc3, 0xd7, 0x4e, 0x9a, 0xfa, 0xf5, 0x49, 0xee, 0x93, 0x9e, 0x5e, 0xea,
	0x57, 0x50, 0x00, 0xa9, 0x11, 0x45, 0x62, 0x1f, 0xa0, 0x8b, 0x2e, 0x5a, 0x55, 0x0e, 0xb0, 0x60,
	0x81, 0xe5, 0xd8, 0xb7, 0xce, 0x28, 0xc9, 0x8c, 0xf1, 0x8c, 0x53, 0x99, 0xbf, 0xc1, 0xbf, 0x62,
	0xc3, 0x0f, 0xe1, 0x4f, 0xa0, 0xf9, 0x70, 0x49, 0xd4, 0xaa, 0xad, 0xd8, 0xcd, 0x3d, 0xe3, 0x39,
	0xf7, 0x9e, 0x73, 0x4f, 0x14, 0xe8, 0x5d, 0xe3, 0x2c, 0x2f, 0xf9, 0xfc, 0x38, 0x2f, 0x98, 0x60,
	0x5e, 0xc7, 0x94, 0xc1, 0x77, 0x0b, 0x06, 0xef, 0x08, 0x4f, 0xd8, 0x1a, 0x8b, 0x2a, 0xc4, 0x2f,
	0x25, 0x72, 0xe1, 0x1d, 0x82, 0xbb, 0xc6, 0x82, 0x13, 0x46, 0x23, 0x42, 0xaf, 0x98, 0x6f, 0x0d,
	0xad, 0x91, 0x1d, 0x3a, 0x06, 0x3b, 0xa3, 0x57, 0xcc, 0x7b, 0x02, 0xfd, 0x02, 0x39, 0x2b, 0x8b,
	0x04, 0x23, 0x1a, 0xaf, 0x90, 0xfb, 0xcd, 0x61, 0x73, 0x64, 0x87, 0xbd, 0x1a, 0xbd, 0x90, 0xa0,
	0x77, 0x00, 0x5d, 0x51, 0xe5, 0x18, 0x95, 0xc5, 0xd2, 0xdf, 0x51, 0x2c, 0x1d, 0x59, 0x7f, 0x28,
	0x96, 0x86, 0x21, 0x67, 0x94, 0x63, 0x44, 0x19, 0x4d, 0xd0, 0x6f, 0xa9, 0x0f, 0x7a, 0x35, 0x7a,
	0x21, 0x41, 0xef, 0x39, 0xd8, 0x35, 0x25, 0xf7, 0x3b, 0xc3, 0xe6, 0xc8, 0x39, 0x71, 0x8f, 0x6b,
	0x31, 0x13, 0x5a, 0x85, 0xbf, 0xaf, 0x83, 0x6f, 0x16, 0xec, 0x6d, 0x88, 0xd1, 0x34, 0x8f, 0x51,
	0xb3, 0xd5, 0xa4, 0x71, 0x6f, 0x93, 0xfb, 0x24, 0xed, 0x43, 0x6b, 0x53, 0x89, 0x2e, 0x82, 0xd7,
	0xd0, 0x9c, 0xd0, 0x6a, 0xeb, 0x9d, 0x75, 0xeb, 0xdd, 0x3a, 0x5e, 0x96, 0xe8, 0x37, 0x86, 0xd6,
	0xc8, 0x0d, 0x75, 0x11, 0xfc, 0xb0, 0x60, 0xf7, 0x1c, 0x39, 0x8f, 0x33, 0x3c, 0xa5, 0x6b, 0x5c,
	0xb2, 0x1c, 0xbd, 0x7f, 0x01, 0x56, 0x1a, 0x8a, 0x48, 0x6a, 0x68, 0x6c, 0x83, 0x9c, 0xa5, 0x9e,
	0x07, 0x3b, 0x72, 0x64, 0xc5, 0x63, 0x87, 0xea, 0x2c, 0xb1, 0x34, 0x16, 0xb1, 0xdf, 0x54, 0xdc,
	0xea, 0xec, 0x3d, 0x85, 0x36, 0x47, 0x9a, 0x62, 0xa1, 0x14, 0x38, 0x27, 0xfd, 0x1b, 0xb1, 0x1f,
	0xe3, 0x9c, 0xa4, 0xa1, 0xb9, 0xf5, 0x9e, 0xc1, 0x20, 0x61, 0x54, 0x20, 0x15, 0x11, 0xd2, 0x84,
	0xa5, 0x84, 0x66, 0x7e, 0x47, 0x71, 0xef, 0x1a, 0xfc, 0xd4, 0xc0, 0xb2, 0x0d, 0x8f, 0x97, 0xc2,
	0xef, 0xea, 0x36, 0xf2, 0xec, 0xf5, 0xa1, 0x91, 0xce, 0x7d, 0x5b, 0x21, 0x8d, 0x74, 0x1e, 0xc4,
	0xd0, 0x52, 0xfc, 0x37, 0x33, 0x75, 0x36, 0x66, 0x72, 0xc1, 0x5a, 0xa8, 0x71, 0xdc, 0xd0, 0x5a,
	0x78, 0x7f, 0x43, 0x5b, 0x44, 0xd2, 0x1f, 0x7f, 0xa8, 0x3d, 0x11, 0xef, 0xab, 0x1c, 0xbd, 0xff,
	0xc0, 0x11, 0x11, 0x27, 0x19, 0x8d, 0x45, 0x59, 0xa0, 0x7f, 0xa8, 0xee, 0x40, 0x4c, 0x6b, 0x24,
	0xb8, 0x84, 0xc1, 0xb4, 0x9c, 0xf1, 0xa4, 0x20, 0x33, 0xac, 0xe3, 0x7c, 0x00, 0x5d, 0xa9, 0x2d,
	0xe2, 0x28, 0x6a, 0xe7, 0x65, 0x3d, 0x45, 0x95, 0x74, 0x2d, 0x35, 0x5a, 0xcb, 0xc1, 0x8c, 0x71,
	0x8e, 0xc6, 0xd4, 0xac, 0xc1, 0x67, 0xd8, 0xdb, 0x60, 0x34, 0x99, 0xaa, 0x8d, 0xb6, 0x36, 0x8c,
	0xde, 0x6c, 0xd3, 0xd8, 0x6e, 0xf3, 0x0f, 0x74, 0x97, 0x2c, 0x89, 0x05, 0x61, 0x54, 0xed, 0xc1,
	0x0e, 0x6f, 0xea, 0xe0, 0xa7, 0x05, 0xce, 0x65, 0xc9, 0xe7, 0xf5, 0xb4, 0x77, 0x51, 0x0f, 0xa0,
	0x29, 0xc4, 0x52, 0xb1, 0xb6, 0x42, 0x79, 0xbc, 0x73, 0xab, 0x3e, 0x74, 0xca, 0x22, 0x43, 0x9a,
	0x54, 0x75, 0x30, 0x4d, 0xe9, 0xfd, 0x0f, 0xe6, 0x57, 0x95, 0x46, 0x31, 0xaf, 0x68, 0xa2, 0x02,
	0xda, 0x0d, 0x5d, 0x03, 0x4e, 0x24, 0x26, 0x53, 0x28, 0x58, 0x4e, 0x12, 0xbf, 0xad, 0xd3, 0xab,
	0x8a, 0x3f, 0x8d, 0x80, 0x7d, 0x2b, 0x02, 0xb6, 0x8a, 0xc0, 0x25, 0xb8, 0x5a, 0xac, 0x31, 0xf2,
	0x81, 0x40, 0x1f, 0x82, 0xab, 0x3c, 0x2d, 0x30, 0x41, 0x92, 0xd7, 0xbe, 0x3a, 0xb9, 0xa2, 0x50,
	0x50, 0x90, 0x40, 0xff, 0x9c, 0x51, 0x22, 0x58, 0xf1, 0x88, 0x7d, 0x1f, 0x41, 0x2f, 0x2e, 0xc5,
	0x9c, 0x15, 0xe4, 0xab, 0xde, 0x86, 0x26, 0xdc, 0x06, 0xa5, 0xdd, 0x0b, 0xac, 0xcc, 0xa6, 0xe4,
	0x31, 0x78, 0x01, 0x30, 0x49, 0x16, 0x75, 0x83, 0xfb, 0x87, 0x0e, 0x7a, 0xe0, 0xa8, 0x8f, 0xb5,
	0xc4, 0xe0, 0x2d, 0xf4, 0xcd, 0xac, 0xf5, 0xfb, 0x97, 0xb0, 0x6f, 0x04, 0x45, 0x5c, 0x47, 0x2b,
	0x57, 0xc3, 0x68, 0xa6, 0xbf, 0xcc, 0xdd, 0x74, 0xe3, 0x2a, 0x18, 0x41, 0xc7, 0x90, 0x3c, 0xd0,
	0xfd, 0xcd, 0xd1, 0xa7, 0x20, 0x23, 0x62, 0x5e, 0xce, 0x8e, 0x13, 0xb6, 0x1a, 0x27, 0x8c, 0x0b,
	0x42, 0x57, 0xe3, 0xeb, 0x3c, 0x8b, 0x05, 0x8e, 0xf3, 0x45, 0x36, 0x5e, 0xf1, 0x8c, 0xcf, 0xda,
	0xea, 0x7f, 0xe0, 0xd5, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x87, 0x03, 0x4a, 0x42, 0x18, 0x06,
	0x00, 0x00,
}
