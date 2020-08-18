// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: opentelemetry/proto/trace/v1/trace_config.proto

package v1

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

// How spans should be sampled:
// - Always off
// - Always on
// - Always follow the parent Span's decision (off if no parent).
type ConstantSampler_ConstantDecision int32

const (
	ConstantSampler_ALWAYS_OFF    ConstantSampler_ConstantDecision = 0
	ConstantSampler_ALWAYS_ON     ConstantSampler_ConstantDecision = 1
	ConstantSampler_ALWAYS_PARENT ConstantSampler_ConstantDecision = 2
)

var ConstantSampler_ConstantDecision_name = map[int32]string{
	0: "ALWAYS_OFF",
	1: "ALWAYS_ON",
	2: "ALWAYS_PARENT",
}

var ConstantSampler_ConstantDecision_value = map[string]int32{
	"ALWAYS_OFF":    0,
	"ALWAYS_ON":     1,
	"ALWAYS_PARENT": 2,
}

func (x ConstantSampler_ConstantDecision) String() string {
	return proto.EnumName(ConstantSampler_ConstantDecision_name, int32(x))
}

func (ConstantSampler_ConstantDecision) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5936aa8fa6443e6f, []int{1, 0}
}

// Global configuration of the trace service. All fields must be specified, or
// the default (zero) values will be used for each type.
type TraceConfig struct {
	// The global default sampler used to make decisions on span sampling.
	//
	// Types that are valid to be assigned to Sampler:
	//	*TraceConfig_ConstantSampler
	//	*TraceConfig_ProbabilitySampler
	//	*TraceConfig_RateLimitingSampler
	Sampler isTraceConfig_Sampler `protobuf_oneof:"sampler"`
	// The global default max number of attributes per span.
	MaxNumberOfAttributes int64 `protobuf:"varint,4,opt,name=max_number_of_attributes,json=maxNumberOfAttributes,proto3" json:"max_number_of_attributes,omitempty"`
	// The global default max number of annotation events per span.
	MaxNumberOfTimedEvents int64 `protobuf:"varint,5,opt,name=max_number_of_timed_events,json=maxNumberOfTimedEvents,proto3" json:"max_number_of_timed_events,omitempty"`
	// The global default max number of attributes per timed event.
	MaxNumberOfAttributesPerTimedEvent int64 `protobuf:"varint,6,opt,name=max_number_of_attributes_per_timed_event,json=maxNumberOfAttributesPerTimedEvent,proto3" json:"max_number_of_attributes_per_timed_event,omitempty"`
	// The global default max number of link entries per span.
	MaxNumberOfLinks int64 `protobuf:"varint,7,opt,name=max_number_of_links,json=maxNumberOfLinks,proto3" json:"max_number_of_links,omitempty"`
	// The global default max number of attributes per span.
	MaxNumberOfAttributesPerLink int64    `protobuf:"varint,8,opt,name=max_number_of_attributes_per_link,json=maxNumberOfAttributesPerLink,proto3" json:"max_number_of_attributes_per_link,omitempty"`
	XXX_NoUnkeyedLiteral         struct{} `json:"-"`
	XXX_unrecognized             []byte   `json:"-"`
	XXX_sizecache                int32    `json:"-"`
}

func (m *TraceConfig) Reset()         { *m = TraceConfig{} }
func (m *TraceConfig) String() string { return proto.CompactTextString(m) }
func (*TraceConfig) ProtoMessage()    {}
func (*TraceConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_5936aa8fa6443e6f, []int{0}
}
func (m *TraceConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TraceConfig.Unmarshal(m, b)
}
func (m *TraceConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TraceConfig.Marshal(b, m, deterministic)
}
func (m *TraceConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceConfig.Merge(m, src)
}
func (m *TraceConfig) XXX_Size() int {
	return xxx_messageInfo_TraceConfig.Size(m)
}
func (m *TraceConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceConfig.DiscardUnknown(m)
}

var xxx_messageInfo_TraceConfig proto.InternalMessageInfo

type isTraceConfig_Sampler interface {
	isTraceConfig_Sampler()
}

type TraceConfig_ConstantSampler struct {
	ConstantSampler *ConstantSampler `protobuf:"bytes,1,opt,name=constant_sampler,json=constantSampler,proto3,oneof" json:"constant_sampler,omitempty"`
}
type TraceConfig_ProbabilitySampler struct {
	ProbabilitySampler *ProbabilitySampler `protobuf:"bytes,2,opt,name=probability_sampler,json=probabilitySampler,proto3,oneof" json:"probability_sampler,omitempty"`
}
type TraceConfig_RateLimitingSampler struct {
	RateLimitingSampler *RateLimitingSampler `protobuf:"bytes,3,opt,name=rate_limiting_sampler,json=rateLimitingSampler,proto3,oneof" json:"rate_limiting_sampler,omitempty"`
}

func (*TraceConfig_ConstantSampler) isTraceConfig_Sampler()     {}
func (*TraceConfig_ProbabilitySampler) isTraceConfig_Sampler()  {}
func (*TraceConfig_RateLimitingSampler) isTraceConfig_Sampler() {}

func (m *TraceConfig) GetSampler() isTraceConfig_Sampler {
	if m != nil {
		return m.Sampler
	}
	return nil
}

func (m *TraceConfig) GetConstantSampler() *ConstantSampler {
	if x, ok := m.GetSampler().(*TraceConfig_ConstantSampler); ok {
		return x.ConstantSampler
	}
	return nil
}

func (m *TraceConfig) GetProbabilitySampler() *ProbabilitySampler {
	if x, ok := m.GetSampler().(*TraceConfig_ProbabilitySampler); ok {
		return x.ProbabilitySampler
	}
	return nil
}

func (m *TraceConfig) GetRateLimitingSampler() *RateLimitingSampler {
	if x, ok := m.GetSampler().(*TraceConfig_RateLimitingSampler); ok {
		return x.RateLimitingSampler
	}
	return nil
}

func (m *TraceConfig) GetMaxNumberOfAttributes() int64 {
	if m != nil {
		return m.MaxNumberOfAttributes
	}
	return 0
}

func (m *TraceConfig) GetMaxNumberOfTimedEvents() int64 {
	if m != nil {
		return m.MaxNumberOfTimedEvents
	}
	return 0
}

func (m *TraceConfig) GetMaxNumberOfAttributesPerTimedEvent() int64 {
	if m != nil {
		return m.MaxNumberOfAttributesPerTimedEvent
	}
	return 0
}

func (m *TraceConfig) GetMaxNumberOfLinks() int64 {
	if m != nil {
		return m.MaxNumberOfLinks
	}
	return 0
}

func (m *TraceConfig) GetMaxNumberOfAttributesPerLink() int64 {
	if m != nil {
		return m.MaxNumberOfAttributesPerLink
	}
	return 0
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*TraceConfig) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*TraceConfig_ConstantSampler)(nil),
		(*TraceConfig_ProbabilitySampler)(nil),
		(*TraceConfig_RateLimitingSampler)(nil),
	}
}

// Sampler that always makes a constant decision on span sampling.
type ConstantSampler struct {
	Decision             ConstantSampler_ConstantDecision `protobuf:"varint,1,opt,name=decision,proto3,enum=opentelemetry.proto.trace.v1.ConstantSampler_ConstantDecision" json:"decision,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *ConstantSampler) Reset()         { *m = ConstantSampler{} }
func (m *ConstantSampler) String() string { return proto.CompactTextString(m) }
func (*ConstantSampler) ProtoMessage()    {}
func (*ConstantSampler) Descriptor() ([]byte, []int) {
	return fileDescriptor_5936aa8fa6443e6f, []int{1}
}
func (m *ConstantSampler) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConstantSampler.Unmarshal(m, b)
}
func (m *ConstantSampler) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConstantSampler.Marshal(b, m, deterministic)
}
func (m *ConstantSampler) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConstantSampler.Merge(m, src)
}
func (m *ConstantSampler) XXX_Size() int {
	return xxx_messageInfo_ConstantSampler.Size(m)
}
func (m *ConstantSampler) XXX_DiscardUnknown() {
	xxx_messageInfo_ConstantSampler.DiscardUnknown(m)
}

var xxx_messageInfo_ConstantSampler proto.InternalMessageInfo

func (m *ConstantSampler) GetDecision() ConstantSampler_ConstantDecision {
	if m != nil {
		return m.Decision
	}
	return ConstantSampler_ALWAYS_OFF
}

// Sampler that tries to uniformly sample traces with a given probability.
// The probability of sampling a trace is equal to that of the specified probability.
type ProbabilitySampler struct {
	// The desired probability of sampling. Must be within [0.0, 1.0].
	SamplingProbability  float64  `protobuf:"fixed64,1,opt,name=samplingProbability,proto3" json:"samplingProbability,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProbabilitySampler) Reset()         { *m = ProbabilitySampler{} }
func (m *ProbabilitySampler) String() string { return proto.CompactTextString(m) }
func (*ProbabilitySampler) ProtoMessage()    {}
func (*ProbabilitySampler) Descriptor() ([]byte, []int) {
	return fileDescriptor_5936aa8fa6443e6f, []int{2}
}
func (m *ProbabilitySampler) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProbabilitySampler.Unmarshal(m, b)
}
func (m *ProbabilitySampler) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProbabilitySampler.Marshal(b, m, deterministic)
}
func (m *ProbabilitySampler) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProbabilitySampler.Merge(m, src)
}
func (m *ProbabilitySampler) XXX_Size() int {
	return xxx_messageInfo_ProbabilitySampler.Size(m)
}
func (m *ProbabilitySampler) XXX_DiscardUnknown() {
	xxx_messageInfo_ProbabilitySampler.DiscardUnknown(m)
}

var xxx_messageInfo_ProbabilitySampler proto.InternalMessageInfo

func (m *ProbabilitySampler) GetSamplingProbability() float64 {
	if m != nil {
		return m.SamplingProbability
	}
	return 0
}

// Sampler that tries to sample with a rate per time window.
type RateLimitingSampler struct {
	// Rate per second.
	Qps                  int64    `protobuf:"varint,1,opt,name=qps,proto3" json:"qps,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RateLimitingSampler) Reset()         { *m = RateLimitingSampler{} }
func (m *RateLimitingSampler) String() string { return proto.CompactTextString(m) }
func (*RateLimitingSampler) ProtoMessage()    {}
func (*RateLimitingSampler) Descriptor() ([]byte, []int) {
	return fileDescriptor_5936aa8fa6443e6f, []int{3}
}
func (m *RateLimitingSampler) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RateLimitingSampler.Unmarshal(m, b)
}
func (m *RateLimitingSampler) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RateLimitingSampler.Marshal(b, m, deterministic)
}
func (m *RateLimitingSampler) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RateLimitingSampler.Merge(m, src)
}
func (m *RateLimitingSampler) XXX_Size() int {
	return xxx_messageInfo_RateLimitingSampler.Size(m)
}
func (m *RateLimitingSampler) XXX_DiscardUnknown() {
	xxx_messageInfo_RateLimitingSampler.DiscardUnknown(m)
}

var xxx_messageInfo_RateLimitingSampler proto.InternalMessageInfo

func (m *RateLimitingSampler) GetQps() int64 {
	if m != nil {
		return m.Qps
	}
	return 0
}

func init() {
	proto.RegisterEnum("opentelemetry.proto.trace.v1.ConstantSampler_ConstantDecision", ConstantSampler_ConstantDecision_name, ConstantSampler_ConstantDecision_value)
	proto.RegisterType((*TraceConfig)(nil), "opentelemetry.proto.trace.v1.TraceConfig")
	proto.RegisterType((*ConstantSampler)(nil), "opentelemetry.proto.trace.v1.ConstantSampler")
	proto.RegisterType((*ProbabilitySampler)(nil), "opentelemetry.proto.trace.v1.ProbabilitySampler")
	proto.RegisterType((*RateLimitingSampler)(nil), "opentelemetry.proto.trace.v1.RateLimitingSampler")
}

func init() {
	proto.RegisterFile("opentelemetry/proto/trace/v1/trace_config.proto", fileDescriptor_5936aa8fa6443e6f)
}

var fileDescriptor_5936aa8fa6443e6f = []byte{
	// 509 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0x4d, 0x6f, 0x9b, 0x4c,
	0x14, 0x85, 0x43, 0xfc, 0xe6, 0xeb, 0x46, 0x49, 0x78, 0x07, 0xa5, 0x42, 0x55, 0xa4, 0xa6, 0x6c,
	0xea, 0x8d, 0x21, 0x4e, 0x17, 0x95, 0xba, 0xa8, 0x64, 0x27, 0x71, 0xba, 0xb0, 0x1c, 0x44, 0x2c,
	0x55, 0xf5, 0x06, 0xc1, 0x64, 0x4c, 0x47, 0x85, 0x19, 0x3a, 0x8c, 0xad, 0x64, 0xd3, 0x5d, 0x7f,
	0x51, 0xff, 0x60, 0xc5, 0x98, 0xf2, 0x11, 0x3b, 0x48, 0xdd, 0x71, 0xef, 0xe1, 0x3c, 0x67, 0xc6,
	0xbe, 0x5c, 0x70, 0x78, 0x4a, 0x98, 0x24, 0x31, 0x49, 0x88, 0x14, 0x4f, 0x4e, 0x2a, 0xb8, 0xe4,
	0x8e, 0x14, 0x01, 0x26, 0xce, 0xb2, 0xbf, 0x7a, 0xf0, 0x31, 0x67, 0x73, 0x1a, 0xd9, 0x4a, 0x43,
	0x67, 0x0d, 0xc3, 0xaa, 0x69, 0xab, 0xf7, 0xec, 0x65, 0xdf, 0xfa, 0xb5, 0x03, 0x87, 0xd3, 0xbc,
	0xb8, 0x52, 0x1e, 0x34, 0x03, 0x1d, 0x73, 0x96, 0xc9, 0x80, 0x49, 0x3f, 0x0b, 0x92, 0x34, 0x26,
	0xc2, 0xd4, 0xce, 0xb5, 0xee, 0xe1, 0x65, 0xcf, 0x6e, 0x03, 0xd9, 0x57, 0x85, 0xeb, 0x7e, 0x65,
	0xfa, 0xbc, 0xe5, 0x9d, 0xe0, 0x66, 0x0b, 0x61, 0x30, 0x52, 0xc1, 0xc3, 0x20, 0xa4, 0x31, 0x95,
	0x4f, 0x25, 0x7e, 0x5b, 0xe1, 0x2f, 0xda, 0xf1, 0x6e, 0x65, 0xac, 0x12, 0x50, 0xba, 0xd6, 0x45,
	0x11, 0x9c, 0x8a, 0x40, 0x12, 0x3f, 0xa6, 0x09, 0x95, 0x94, 0x45, 0x65, 0x4c, 0x47, 0xc5, 0xf4,
	0xdb, 0x63, 0xbc, 0x40, 0x92, 0x71, 0xe1, 0xac, 0x72, 0x0c, 0xb1, 0xde, 0x46, 0x1f, 0xc0, 0x4c,
	0x82, 0x47, 0x9f, 0x2d, 0x92, 0x90, 0x08, 0x9f, 0xcf, 0xfd, 0x40, 0x4a, 0x41, 0xc3, 0x85, 0x24,
	0x99, 0xf9, 0xdf, 0xb9, 0xd6, 0xed, 0x78, 0xa7, 0x49, 0xf0, 0x38, 0x51, 0xf2, 0xdd, 0x7c, 0x50,
	0x8a, 0xe8, 0x23, 0xbc, 0x6e, 0x1a, 0x25, 0x4d, 0xc8, 0x83, 0x4f, 0x96, 0x84, 0xc9, 0xcc, 0xdc,
	0x51, 0xd6, 0x57, 0x35, 0xeb, 0x34, 0x97, 0x6f, 0x94, 0x8a, 0xa6, 0xd0, 0x7d, 0x29, 0xd4, 0x4f,
	0x89, 0xa8, 0xa3, 0xcc, 0x5d, 0x45, 0xb2, 0x36, 0x1e, 0xc2, 0x25, 0xa2, 0xc2, 0xa2, 0x1e, 0x18,
	0x4d, 0x6a, 0x4c, 0xd9, 0xf7, 0xcc, 0xdc, 0x53, 0x00, 0xbd, 0x06, 0x18, 0xe7, 0x7d, 0x74, 0x0b,
	0x6f, 0x5b, 0x0f, 0x91, 0xbb, 0xcd, 0x7d, 0x65, 0x3e, 0x7b, 0x29, 0x3d, 0x27, 0x0d, 0x0f, 0x60,
	0xaf, 0xf8, 0x77, 0xac, 0xdf, 0x1a, 0x9c, 0x3c, 0x1b, 0x21, 0x34, 0x83, 0xfd, 0x07, 0x82, 0x69,
	0x46, 0x39, 0x53, 0x33, 0x78, 0x7c, 0xf9, 0xe9, 0x9f, 0x66, 0xb0, 0xac, 0xaf, 0x0b, 0x8a, 0x57,
	0xf2, 0xac, 0x6b, 0xd0, 0x9f, 0xab, 0xe8, 0x18, 0x60, 0x30, 0xfe, 0x32, 0xf8, 0x7a, 0xef, 0xdf,
	0x8d, 0x46, 0xfa, 0x16, 0x3a, 0x82, 0x83, 0xbf, 0xf5, 0x44, 0xd7, 0xd0, 0xff, 0x70, 0x54, 0x94,
	0xee, 0xc0, 0xbb, 0x99, 0x4c, 0xf5, 0x6d, 0x6b, 0x04, 0x68, 0x7d, 0x30, 0xd1, 0x05, 0x18, 0xea,
	0x5a, 0x94, 0x45, 0x35, 0x55, 0x5d, 0x41, 0xf3, 0x36, 0x49, 0xd6, 0x3b, 0x30, 0x36, 0x4c, 0x1e,
	0xd2, 0xa1, 0xf3, 0x23, 0xcd, 0x94, 0xb1, 0xe3, 0xe5, 0x8f, 0xc3, 0x9f, 0xf0, 0x86, 0xf2, 0xd6,
	0x1f, 0x61, 0xa8, 0xd7, 0x3e, 0x67, 0x37, 0x97, 0x5c, 0x6d, 0x76, 0x1b, 0x51, 0xf9, 0x6d, 0x11,
	0xda, 0x98, 0x27, 0x6a, 0x7f, 0xf4, 0xaa, 0x05, 0xd2, 0x60, 0xf5, 0x56, 0xeb, 0x24, 0x22, 0xcc,
	0x89, 0xb8, 0x83, 0x79, 0x1c, 0x13, 0x2c, 0xb9, 0x28, 0xf7, 0x4b, 0xb8, 0xab, 0x5e, 0x78, 0xff,
	0x27, 0x00, 0x00, 0xff, 0xff, 0xe4, 0x80, 0x54, 0x8b, 0x86, 0x04, 0x00, 0x00,
}
