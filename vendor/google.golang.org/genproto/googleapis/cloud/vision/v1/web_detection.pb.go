// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/vision/v1/web_detection.proto

package vision

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Relevant information for the image from the Internet.
type WebDetection struct {
	// Deduced entities from similar images on the Internet.
	WebEntities []*WebDetection_WebEntity `protobuf:"bytes,1,rep,name=web_entities,json=webEntities,proto3" json:"web_entities,omitempty"`
	// Fully matching images from the Internet.
	// Can include resized copies of the query image.
	FullMatchingImages []*WebDetection_WebImage `protobuf:"bytes,2,rep,name=full_matching_images,json=fullMatchingImages,proto3" json:"full_matching_images,omitempty"`
	// Partial matching images from the Internet.
	// Those images are similar enough to share some key-point features. For
	// example an original image will likely have partial matching for its crops.
	PartialMatchingImages []*WebDetection_WebImage `protobuf:"bytes,3,rep,name=partial_matching_images,json=partialMatchingImages,proto3" json:"partial_matching_images,omitempty"`
	// Web pages containing the matching images from the Internet.
	PagesWithMatchingImages []*WebDetection_WebPage `protobuf:"bytes,4,rep,name=pages_with_matching_images,json=pagesWithMatchingImages,proto3" json:"pages_with_matching_images,omitempty"`
	// The visually similar image results.
	VisuallySimilarImages []*WebDetection_WebImage `protobuf:"bytes,6,rep,name=visually_similar_images,json=visuallySimilarImages,proto3" json:"visually_similar_images,omitempty"`
	// The service's best guess as to the topic of the request image.
	// Inferred from similar images on the open web.
	BestGuessLabels      []*WebDetection_WebLabel `protobuf:"bytes,8,rep,name=best_guess_labels,json=bestGuessLabels,proto3" json:"best_guess_labels,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *WebDetection) Reset()         { *m = WebDetection{} }
func (m *WebDetection) String() string { return proto.CompactTextString(m) }
func (*WebDetection) ProtoMessage()    {}
func (*WebDetection) Descriptor() ([]byte, []int) {
	return fileDescriptor_894df371610d13f4, []int{0}
}

func (m *WebDetection) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WebDetection.Unmarshal(m, b)
}
func (m *WebDetection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WebDetection.Marshal(b, m, deterministic)
}
func (m *WebDetection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WebDetection.Merge(m, src)
}
func (m *WebDetection) XXX_Size() int {
	return xxx_messageInfo_WebDetection.Size(m)
}
func (m *WebDetection) XXX_DiscardUnknown() {
	xxx_messageInfo_WebDetection.DiscardUnknown(m)
}

var xxx_messageInfo_WebDetection proto.InternalMessageInfo

func (m *WebDetection) GetWebEntities() []*WebDetection_WebEntity {
	if m != nil {
		return m.WebEntities
	}
	return nil
}

func (m *WebDetection) GetFullMatchingImages() []*WebDetection_WebImage {
	if m != nil {
		return m.FullMatchingImages
	}
	return nil
}

func (m *WebDetection) GetPartialMatchingImages() []*WebDetection_WebImage {
	if m != nil {
		return m.PartialMatchingImages
	}
	return nil
}

func (m *WebDetection) GetPagesWithMatchingImages() []*WebDetection_WebPage {
	if m != nil {
		return m.PagesWithMatchingImages
	}
	return nil
}

func (m *WebDetection) GetVisuallySimilarImages() []*WebDetection_WebImage {
	if m != nil {
		return m.VisuallySimilarImages
	}
	return nil
}

func (m *WebDetection) GetBestGuessLabels() []*WebDetection_WebLabel {
	if m != nil {
		return m.BestGuessLabels
	}
	return nil
}

// Entity deduced from similar images on the Internet.
type WebDetection_WebEntity struct {
	// Opaque entity ID.
	EntityId string `protobuf:"bytes,1,opt,name=entity_id,json=entityId,proto3" json:"entity_id,omitempty"`
	// Overall relevancy score for the entity.
	// Not normalized and not comparable across different image queries.
	Score float32 `protobuf:"fixed32,2,opt,name=score,proto3" json:"score,omitempty"`
	// Canonical description of the entity, in English.
	Description          string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WebDetection_WebEntity) Reset()         { *m = WebDetection_WebEntity{} }
func (m *WebDetection_WebEntity) String() string { return proto.CompactTextString(m) }
func (*WebDetection_WebEntity) ProtoMessage()    {}
func (*WebDetection_WebEntity) Descriptor() ([]byte, []int) {
	return fileDescriptor_894df371610d13f4, []int{0, 0}
}

func (m *WebDetection_WebEntity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WebDetection_WebEntity.Unmarshal(m, b)
}
func (m *WebDetection_WebEntity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WebDetection_WebEntity.Marshal(b, m, deterministic)
}
func (m *WebDetection_WebEntity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WebDetection_WebEntity.Merge(m, src)
}
func (m *WebDetection_WebEntity) XXX_Size() int {
	return xxx_messageInfo_WebDetection_WebEntity.Size(m)
}
func (m *WebDetection_WebEntity) XXX_DiscardUnknown() {
	xxx_messageInfo_WebDetection_WebEntity.DiscardUnknown(m)
}

var xxx_messageInfo_WebDetection_WebEntity proto.InternalMessageInfo

func (m *WebDetection_WebEntity) GetEntityId() string {
	if m != nil {
		return m.EntityId
	}
	return ""
}

func (m *WebDetection_WebEntity) GetScore() float32 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *WebDetection_WebEntity) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

// Metadata for online images.
type WebDetection_WebImage struct {
	// The result image URL.
	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	// (Deprecated) Overall relevancy score for the image.
	Score                float32  `protobuf:"fixed32,2,opt,name=score,proto3" json:"score,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WebDetection_WebImage) Reset()         { *m = WebDetection_WebImage{} }
func (m *WebDetection_WebImage) String() string { return proto.CompactTextString(m) }
func (*WebDetection_WebImage) ProtoMessage()    {}
func (*WebDetection_WebImage) Descriptor() ([]byte, []int) {
	return fileDescriptor_894df371610d13f4, []int{0, 1}
}

func (m *WebDetection_WebImage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WebDetection_WebImage.Unmarshal(m, b)
}
func (m *WebDetection_WebImage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WebDetection_WebImage.Marshal(b, m, deterministic)
}
func (m *WebDetection_WebImage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WebDetection_WebImage.Merge(m, src)
}
func (m *WebDetection_WebImage) XXX_Size() int {
	return xxx_messageInfo_WebDetection_WebImage.Size(m)
}
func (m *WebDetection_WebImage) XXX_DiscardUnknown() {
	xxx_messageInfo_WebDetection_WebImage.DiscardUnknown(m)
}

var xxx_messageInfo_WebDetection_WebImage proto.InternalMessageInfo

func (m *WebDetection_WebImage) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *WebDetection_WebImage) GetScore() float32 {
	if m != nil {
		return m.Score
	}
	return 0
}

// Metadata for web pages.
type WebDetection_WebPage struct {
	// The result web page URL.
	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	// (Deprecated) Overall relevancy score for the web page.
	Score float32 `protobuf:"fixed32,2,opt,name=score,proto3" json:"score,omitempty"`
	// Title for the web page, may contain HTML markups.
	PageTitle string `protobuf:"bytes,3,opt,name=page_title,json=pageTitle,proto3" json:"page_title,omitempty"`
	// Fully matching images on the page.
	// Can include resized copies of the query image.
	FullMatchingImages []*WebDetection_WebImage `protobuf:"bytes,4,rep,name=full_matching_images,json=fullMatchingImages,proto3" json:"full_matching_images,omitempty"`
	// Partial matching images on the page.
	// Those images are similar enough to share some key-point features. For
	// example an original image will likely have partial matching for its
	// crops.
	PartialMatchingImages []*WebDetection_WebImage `protobuf:"bytes,5,rep,name=partial_matching_images,json=partialMatchingImages,proto3" json:"partial_matching_images,omitempty"`
	XXX_NoUnkeyedLiteral  struct{}                 `json:"-"`
	XXX_unrecognized      []byte                   `json:"-"`
	XXX_sizecache         int32                    `json:"-"`
}

func (m *WebDetection_WebPage) Reset()         { *m = WebDetection_WebPage{} }
func (m *WebDetection_WebPage) String() string { return proto.CompactTextString(m) }
func (*WebDetection_WebPage) ProtoMessage()    {}
func (*WebDetection_WebPage) Descriptor() ([]byte, []int) {
	return fileDescriptor_894df371610d13f4, []int{0, 2}
}

func (m *WebDetection_WebPage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WebDetection_WebPage.Unmarshal(m, b)
}
func (m *WebDetection_WebPage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WebDetection_WebPage.Marshal(b, m, deterministic)
}
func (m *WebDetection_WebPage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WebDetection_WebPage.Merge(m, src)
}
func (m *WebDetection_WebPage) XXX_Size() int {
	return xxx_messageInfo_WebDetection_WebPage.Size(m)
}
func (m *WebDetection_WebPage) XXX_DiscardUnknown() {
	xxx_messageInfo_WebDetection_WebPage.DiscardUnknown(m)
}

var xxx_messageInfo_WebDetection_WebPage proto.InternalMessageInfo

func (m *WebDetection_WebPage) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *WebDetection_WebPage) GetScore() float32 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *WebDetection_WebPage) GetPageTitle() string {
	if m != nil {
		return m.PageTitle
	}
	return ""
}

func (m *WebDetection_WebPage) GetFullMatchingImages() []*WebDetection_WebImage {
	if m != nil {
		return m.FullMatchingImages
	}
	return nil
}

func (m *WebDetection_WebPage) GetPartialMatchingImages() []*WebDetection_WebImage {
	if m != nil {
		return m.PartialMatchingImages
	}
	return nil
}

// Label to provide extra metadata for the web detection.
type WebDetection_WebLabel struct {
	// Label for extra metadata.
	Label string `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	// The BCP-47 language code for `label`, such as "en-US" or "sr-Latn".
	// For more information, see
	// http://www.unicode.org/reports/tr35/#Unicode_locale_identifier.
	LanguageCode         string   `protobuf:"bytes,2,opt,name=language_code,json=languageCode,proto3" json:"language_code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WebDetection_WebLabel) Reset()         { *m = WebDetection_WebLabel{} }
func (m *WebDetection_WebLabel) String() string { return proto.CompactTextString(m) }
func (*WebDetection_WebLabel) ProtoMessage()    {}
func (*WebDetection_WebLabel) Descriptor() ([]byte, []int) {
	return fileDescriptor_894df371610d13f4, []int{0, 3}
}

func (m *WebDetection_WebLabel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WebDetection_WebLabel.Unmarshal(m, b)
}
func (m *WebDetection_WebLabel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WebDetection_WebLabel.Marshal(b, m, deterministic)
}
func (m *WebDetection_WebLabel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WebDetection_WebLabel.Merge(m, src)
}
func (m *WebDetection_WebLabel) XXX_Size() int {
	return xxx_messageInfo_WebDetection_WebLabel.Size(m)
}
func (m *WebDetection_WebLabel) XXX_DiscardUnknown() {
	xxx_messageInfo_WebDetection_WebLabel.DiscardUnknown(m)
}

var xxx_messageInfo_WebDetection_WebLabel proto.InternalMessageInfo

func (m *WebDetection_WebLabel) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *WebDetection_WebLabel) GetLanguageCode() string {
	if m != nil {
		return m.LanguageCode
	}
	return ""
}

func init() {
	proto.RegisterType((*WebDetection)(nil), "google.cloud.vision.v1.WebDetection")
	proto.RegisterType((*WebDetection_WebEntity)(nil), "google.cloud.vision.v1.WebDetection.WebEntity")
	proto.RegisterType((*WebDetection_WebImage)(nil), "google.cloud.vision.v1.WebDetection.WebImage")
	proto.RegisterType((*WebDetection_WebPage)(nil), "google.cloud.vision.v1.WebDetection.WebPage")
	proto.RegisterType((*WebDetection_WebLabel)(nil), "google.cloud.vision.v1.WebDetection.WebLabel")
}

func init() {
	proto.RegisterFile("google/cloud/vision/v1/web_detection.proto", fileDescriptor_894df371610d13f4)
}

var fileDescriptor_894df371610d13f4 = []byte{
	// 512 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x94, 0x51, 0x6f, 0xd3, 0x30,
	0x10, 0xc7, 0x95, 0xb6, 0x1b, 0xad, 0x5b, 0x04, 0xb3, 0x06, 0x8b, 0x02, 0x48, 0x15, 0xbc, 0x54,
	0x08, 0x12, 0x6d, 0x3c, 0xc2, 0xd3, 0xc6, 0x34, 0x4d, 0x02, 0x54, 0x02, 0x62, 0x82, 0x17, 0xe3,
	0x26, 0x26, 0x3d, 0xc9, 0x8d, 0xa3, 0xd8, 0x69, 0xd5, 0x6f, 0xc2, 0x33, 0x1f, 0x88, 0xcf, 0xc3,
	0x23, 0x3a, 0xdb, 0x45, 0xd5, 0xba, 0x49, 0x65, 0x42, 0xbc, 0xdd, 0x5d, 0xee, 0xff, 0xfb, 0xdb,
	0x97, 0x93, 0xc9, 0xd3, 0x42, 0xa9, 0x42, 0x8a, 0x24, 0x93, 0xaa, 0xc9, 0x93, 0x39, 0x68, 0x50,
	0x65, 0x32, 0x3f, 0x4c, 0x16, 0x62, 0xc2, 0x72, 0x61, 0x44, 0x66, 0x40, 0x95, 0x71, 0x55, 0x2b,
	0xa3, 0xe8, 0x7d, 0xd7, 0x1b, 0xdb, 0xde, 0xd8, 0xf5, 0xc6, 0xf3, 0xc3, 0xe8, 0xa1, 0x67, 0xf0,
	0x0a, 0x12, 0x5e, 0x96, 0xca, 0x70, 0x14, 0x69, 0xa7, 0x7a, 0xfc, 0xb3, 0x4b, 0x06, 0x17, 0x62,
	0xf2, 0x7a, 0x05, 0xa3, 0xef, 0xc9, 0x00, 0xe9, 0xa2, 0x34, 0x60, 0x40, 0xe8, 0x30, 0x18, 0xb6,
	0x47, 0xfd, 0xa3, 0x38, 0xbe, 0x9a, 0x1e, 0xaf, 0x6b, 0x31, 0x39, 0x45, 0xdd, 0x32, 0xed, 0x2f,
	0x7c, 0x08, 0x42, 0x53, 0x46, 0xf6, 0xbf, 0x35, 0x52, 0xb2, 0x19, 0x37, 0xd9, 0x14, 0xca, 0x82,
	0xc1, 0x8c, 0x17, 0x42, 0x87, 0x2d, 0x8b, 0x7e, 0xbe, 0x2d, 0xfa, 0x1c, 0x55, 0x29, 0x45, 0xd4,
	0x5b, 0x4f, 0xb2, 0x25, 0x4d, 0x05, 0x39, 0xa8, 0x78, 0x6d, 0x80, 0x6f, 0x7a, 0xb4, 0x6f, 0xe2,
	0x71, 0xcf, 0xd3, 0x2e, 0xd9, 0x00, 0x89, 0x2a, 0x0c, 0xd8, 0x02, 0xcc, 0x74, 0xc3, 0xa9, 0x63,
	0x9d, 0x9e, 0x6d, 0xeb, 0x34, 0x46, 0xa3, 0x03, 0xcb, 0xbb, 0x00, 0x33, 0xdd, 0xbc, 0xd1, 0x1c,
	0x74, 0xc3, 0xa5, 0x5c, 0x32, 0x0d, 0x33, 0x90, 0xbc, 0x5e, 0xf9, 0xec, 0xde, 0xe8, 0x46, 0x2b,
	0xda, 0x07, 0x07, 0xf3, 0x36, 0x9f, 0xc9, 0xde, 0x44, 0x68, 0xc3, 0x8a, 0x46, 0x68, 0xcd, 0x24,
	0x9f, 0x08, 0xa9, 0xc3, 0xee, 0xdf, 0x19, 0xbc, 0x41, 0x55, 0x7a, 0x07, 0x39, 0x67, 0x88, 0xb1,
	0xb9, 0x8e, 0xbe, 0x92, 0xde, 0x9f, 0x75, 0xa0, 0x0f, 0x48, 0xcf, 0x2e, 0xd4, 0x92, 0x41, 0x1e,
	0x06, 0xc3, 0x60, 0xd4, 0x4b, 0xbb, 0xae, 0x70, 0x9e, 0xd3, 0x7d, 0xb2, 0xa3, 0x33, 0x55, 0x8b,
	0xb0, 0x35, 0x0c, 0x46, 0xad, 0xd4, 0x25, 0x74, 0x48, 0xfa, 0xb9, 0xd0, 0x59, 0x0d, 0x15, 0x1a,
	0x85, 0x6d, 0x2b, 0x5a, 0x2f, 0x45, 0x47, 0xa4, 0xbb, 0xba, 0x1f, 0xbd, 0x4b, 0xda, 0x4d, 0x2d,
	0x3d, 0x1a, 0xc3, 0xab, 0xa9, 0xd1, 0xf7, 0x16, 0xb9, 0xe5, 0x87, 0xbf, 0xad, 0x86, 0x3e, 0x22,
	0x04, 0x7f, 0x13, 0x33, 0x60, 0xa4, 0xf0, 0x07, 0xe9, 0x61, 0xe5, 0x23, 0x16, 0xae, 0xdd, 0xee,
	0xce, 0x7f, 0xd8, 0xee, 0x9d, 0x7f, 0xb7, 0xdd, 0xd1, 0xa9, 0x1d, 0xa7, 0xfd, 0x7b, 0x38, 0x08,
	0xbb, 0x0c, 0x7e, 0x38, 0x2e, 0xa1, 0x4f, 0xc8, 0x6d, 0xc9, 0xcb, 0xa2, 0xc1, 0x61, 0x64, 0x2a,
	0x77, 0x63, 0xea, 0xa5, 0x83, 0x55, 0xf1, 0x44, 0xe5, 0xe2, 0x78, 0x49, 0xa2, 0x4c, 0xcd, 0xae,
	0x39, 0xd1, 0xf1, 0xde, 0xfa, 0x91, 0xc6, 0xf8, 0x02, 0x8d, 0x83, 0x2f, 0xaf, 0x7c, 0x73, 0xa1,
	0x90, 0x14, 0xab, 0xba, 0x48, 0x0a, 0x51, 0xda, 0xf7, 0x29, 0x71, 0x9f, 0x78, 0x05, 0xfa, 0xf2,
	0x23, 0xf8, 0xd2, 0x45, 0xbf, 0x82, 0xe0, 0x47, 0xab, 0x73, 0x76, 0xf2, 0xe9, 0xdd, 0x64, 0xd7,
	0x4a, 0x5e, 0xfc, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x99, 0x1b, 0x06, 0x03, 0x36, 0x05, 0x00, 0x00,
}
