// Code generated by protoc-gen-go from "crashme/gtfs.proto"
// DO NOT EDIT!

package gtfs

import proto "code.google.com/p/goprotobuf/proto"
import "math"

// Reference proto, math & os imports to suppress error if they are not otherwise used.
var _ = proto.GetString
var _ = math.Inf
var _ error

type FeedHeader_Incrementality int32

const (
	FeedHeader_FULL_DATASET FeedHeader_Incrementality = 0
	FeedHeader_DIFFERENTIAL FeedHeader_Incrementality = 1
)

var FeedHeader_Incrementality_name = map[int32]string{
	0: "FULL_DATASET",
	1: "DIFFERENTIAL",
}
var FeedHeader_Incrementality_value = map[string]int32{
	"FULL_DATASET": 0,
	"DIFFERENTIAL": 1,
}

func NewFeedHeader_Incrementality(x FeedHeader_Incrementality) *FeedHeader_Incrementality {
	e := FeedHeader_Incrementality(x)
	return &e
}
func (x FeedHeader_Incrementality) String() string {
	return proto.EnumName(FeedHeader_Incrementality_name, int32(x))
}

type TripUpdate_StopTimeUpdate_ScheduleRelationship int32

const (
	TripUpdate_StopTimeUpdate_SCHEDULED TripUpdate_StopTimeUpdate_ScheduleRelationship = 0
	TripUpdate_StopTimeUpdate_SKIPPED   TripUpdate_StopTimeUpdate_ScheduleRelationship = 1
	TripUpdate_StopTimeUpdate_NO_DATA   TripUpdate_StopTimeUpdate_ScheduleRelationship = 2
)

var TripUpdate_StopTimeUpdate_ScheduleRelationship_name = map[int32]string{
	0: "SCHEDULED",
	1: "SKIPPED",
	2: "NO_DATA",
}
var TripUpdate_StopTimeUpdate_ScheduleRelationship_value = map[string]int32{
	"SCHEDULED": 0,
	"SKIPPED":   1,
	"NO_DATA":   2,
}

func NewTripUpdate_StopTimeUpdate_ScheduleRelationship(x TripUpdate_StopTimeUpdate_ScheduleRelationship) *TripUpdate_StopTimeUpdate_ScheduleRelationship {
	e := TripUpdate_StopTimeUpdate_ScheduleRelationship(x)
	return &e
}
func (x TripUpdate_StopTimeUpdate_ScheduleRelationship) String() string {
	return proto.EnumName(TripUpdate_StopTimeUpdate_ScheduleRelationship_name, int32(x))
}

type VehiclePosition_VehicleStopStatus int32

const (
	VehiclePosition_INCOMING_AT   VehiclePosition_VehicleStopStatus = 0
	VehiclePosition_STOPPED_AT    VehiclePosition_VehicleStopStatus = 1
	VehiclePosition_IN_TRANSIT_TO VehiclePosition_VehicleStopStatus = 2
)

var VehiclePosition_VehicleStopStatus_name = map[int32]string{
	0: "INCOMING_AT",
	1: "STOPPED_AT",
	2: "IN_TRANSIT_TO",
}
var VehiclePosition_VehicleStopStatus_value = map[string]int32{
	"INCOMING_AT":   0,
	"STOPPED_AT":    1,
	"IN_TRANSIT_TO": 2,
}

func NewVehiclePosition_VehicleStopStatus(x VehiclePosition_VehicleStopStatus) *VehiclePosition_VehicleStopStatus {
	e := VehiclePosition_VehicleStopStatus(x)
	return &e
}
func (x VehiclePosition_VehicleStopStatus) String() string {
	return proto.EnumName(VehiclePosition_VehicleStopStatus_name, int32(x))
}

type VehiclePosition_CongestionLevel int32

const (
	VehiclePosition_UNKNOWN_CONGESTION_LEVEL VehiclePosition_CongestionLevel = 0
	VehiclePosition_RUNNING_SMOOTHLY         VehiclePosition_CongestionLevel = 1
	VehiclePosition_STOP_AND_GO              VehiclePosition_CongestionLevel = 2
	VehiclePosition_CONGESTION               VehiclePosition_CongestionLevel = 3
	VehiclePosition_SEVERE_CONGESTION        VehiclePosition_CongestionLevel = 4
)

var VehiclePosition_CongestionLevel_name = map[int32]string{
	0: "UNKNOWN_CONGESTION_LEVEL",
	1: "RUNNING_SMOOTHLY",
	2: "STOP_AND_GO",
	3: "CONGESTION",
	4: "SEVERE_CONGESTION",
}
var VehiclePosition_CongestionLevel_value = map[string]int32{
	"UNKNOWN_CONGESTION_LEVEL": 0,
	"RUNNING_SMOOTHLY":         1,
	"STOP_AND_GO":              2,
	"CONGESTION":               3,
	"SEVERE_CONGESTION":        4,
}

func NewVehiclePosition_CongestionLevel(x VehiclePosition_CongestionLevel) *VehiclePosition_CongestionLevel {
	e := VehiclePosition_CongestionLevel(x)
	return &e
}
func (x VehiclePosition_CongestionLevel) String() string {
	return proto.EnumName(VehiclePosition_CongestionLevel_name, int32(x))
}

type Alert_Cause int32

const (
	Alert_UNKNOWN_CAUSE     Alert_Cause = 1
	Alert_OTHER_CAUSE       Alert_Cause = 2
	Alert_TECHNICAL_PROBLEM Alert_Cause = 3
	Alert_STRIKE            Alert_Cause = 4
	Alert_DEMONSTRATION     Alert_Cause = 5
	Alert_ACCIDENT          Alert_Cause = 6
	Alert_HOLIDAY           Alert_Cause = 7
	Alert_WEATHER           Alert_Cause = 8
	Alert_MAINTENANCE       Alert_Cause = 9
	Alert_CONSTRUCTION      Alert_Cause = 10
	Alert_POLICE_ACTIVITY   Alert_Cause = 11
	Alert_MEDICAL_EMERGENCY Alert_Cause = 12
)

var Alert_Cause_name = map[int32]string{
	1:  "UNKNOWN_CAUSE",
	2:  "OTHER_CAUSE",
	3:  "TECHNICAL_PROBLEM",
	4:  "STRIKE",
	5:  "DEMONSTRATION",
	6:  "ACCIDENT",
	7:  "HOLIDAY",
	8:  "WEATHER",
	9:  "MAINTENANCE",
	10: "CONSTRUCTION",
	11: "POLICE_ACTIVITY",
	12: "MEDICAL_EMERGENCY",
}
var Alert_Cause_value = map[string]int32{
	"UNKNOWN_CAUSE":     1,
	"OTHER_CAUSE":       2,
	"TECHNICAL_PROBLEM": 3,
	"STRIKE":            4,
	"DEMONSTRATION":     5,
	"ACCIDENT":          6,
	"HOLIDAY":           7,
	"WEATHER":           8,
	"MAINTENANCE":       9,
	"CONSTRUCTION":      10,
	"POLICE_ACTIVITY":   11,
	"MEDICAL_EMERGENCY": 12,
}

func NewAlert_Cause(x Alert_Cause) *Alert_Cause {
	e := Alert_Cause(x)
	return &e
}
func (x Alert_Cause) String() string {
	return proto.EnumName(Alert_Cause_name, int32(x))
}

type Alert_Effect int32

const (
	Alert_NO_SERVICE         Alert_Effect = 1
	Alert_REDUCED_SERVICE    Alert_Effect = 2
	Alert_SIGNIFICANT_DELAYS Alert_Effect = 3
	Alert_DETOUR             Alert_Effect = 4
	Alert_ADDITIONAL_SERVICE Alert_Effect = 5
	Alert_MODIFIED_SERVICE   Alert_Effect = 6
	Alert_OTHER_EFFECT       Alert_Effect = 7
	Alert_UNKNOWN_EFFECT     Alert_Effect = 8
	Alert_STOP_MOVED         Alert_Effect = 9
)

var Alert_Effect_name = map[int32]string{
	1: "NO_SERVICE",
	2: "REDUCED_SERVICE",
	3: "SIGNIFICANT_DELAYS",
	4: "DETOUR",
	5: "ADDITIONAL_SERVICE",
	6: "MODIFIED_SERVICE",
	7: "OTHER_EFFECT",
	8: "UNKNOWN_EFFECT",
	9: "STOP_MOVED",
}
var Alert_Effect_value = map[string]int32{
	"NO_SERVICE":         1,
	"REDUCED_SERVICE":    2,
	"SIGNIFICANT_DELAYS": 3,
	"DETOUR":             4,
	"ADDITIONAL_SERVICE": 5,
	"MODIFIED_SERVICE":   6,
	"OTHER_EFFECT":       7,
	"UNKNOWN_EFFECT":     8,
	"STOP_MOVED":         9,
}

func NewAlert_Effect(x Alert_Effect) *Alert_Effect {
	e := Alert_Effect(x)
	return &e
}
func (x Alert_Effect) String() string {
	return proto.EnumName(Alert_Effect_name, int32(x))
}

type TripDescriptor_ScheduleRelationship int32

const (
	TripDescriptor_SCHEDULED   TripDescriptor_ScheduleRelationship = 0
	TripDescriptor_ADDED       TripDescriptor_ScheduleRelationship = 1
	TripDescriptor_UNSCHEDULED TripDescriptor_ScheduleRelationship = 2
	TripDescriptor_CANCELED    TripDescriptor_ScheduleRelationship = 3
	TripDescriptor_REPLACEMENT TripDescriptor_ScheduleRelationship = 5
)

var TripDescriptor_ScheduleRelationship_name = map[int32]string{
	0: "SCHEDULED",
	1: "ADDED",
	2: "UNSCHEDULED",
	3: "CANCELED",
	5: "REPLACEMENT",
}
var TripDescriptor_ScheduleRelationship_value = map[string]int32{
	"SCHEDULED":   0,
	"ADDED":       1,
	"UNSCHEDULED": 2,
	"CANCELED":    3,
	"REPLACEMENT": 5,
}

func NewTripDescriptor_ScheduleRelationship(x TripDescriptor_ScheduleRelationship) *TripDescriptor_ScheduleRelationship {
	e := TripDescriptor_ScheduleRelationship(x)
	return &e
}
func (x TripDescriptor_ScheduleRelationship) String() string {
	return proto.EnumName(TripDescriptor_ScheduleRelationship_name, int32(x))
}

type FeedMessage struct {
	Header           *FeedHeader   `protobuf:"bytes,1,req,name=header" json:"header"`
	Entity           []*FeedEntity `protobuf:"bytes,2,rep,name=entity" json:"entity"`
	XXX_unrecognized []byte
}

func (this *FeedMessage) Reset()         { *this = FeedMessage{} }
func (this *FeedMessage) String() string { return proto.CompactTextString(this) }

type FeedHeader struct {
	GtfsRealtimeVersion *string                    `protobuf:"bytes,1,req,name=gtfs_realtime_version" json:"gtfs_realtime_version"`
	Incrementality      *FeedHeader_Incrementality `protobuf:"varint,2,opt,name=incrementality,enum=transit_realtime.FeedHeader_Incrementality,def=0" json:"incrementality"`
	Timestamp           *uint64                    `protobuf:"varint,3,opt,name=timestamp" json:"timestamp"`
	XXX_unrecognized    []byte
}

func (this *FeedHeader) Reset()         { *this = FeedHeader{} }
func (this *FeedHeader) String() string { return proto.CompactTextString(this) }

const Default_FeedHeader_Incrementality FeedHeader_Incrementality = FeedHeader_FULL_DATASET

type FeedEntity struct {
	Id               *string          `protobuf:"bytes,1,req,name=id" json:"id"`
	IsDeleted        *bool            `protobuf:"varint,2,opt,name=is_deleted,def=0" json:"is_deleted"`
	TripUpdate       *TripUpdate      `protobuf:"bytes,3,opt,name=trip_update" json:"trip_update"`
	Vehicle          *VehiclePosition `protobuf:"bytes,4,opt,name=vehicle" json:"vehicle"`
	Alert            *Alert           `protobuf:"bytes,5,opt,name=alert" json:"alert"`
	XXX_unrecognized []byte
}

func (this *FeedEntity) Reset()         { *this = FeedEntity{} }
func (this *FeedEntity) String() string { return proto.CompactTextString(this) }

const Default_FeedEntity_IsDeleted bool = false

type TripUpdate struct {
	Trip             *TripDescriptor              `protobuf:"bytes,1,req,name=trip" json:"trip"`
	Vehicle          *VehicleDescriptor           `protobuf:"bytes,3,opt,name=vehicle" json:"vehicle"`
	StopTimeUpdate   []*TripUpdate_StopTimeUpdate `protobuf:"bytes,2,rep,name=stop_time_update" json:"stop_time_update"`
	XXX_unrecognized []byte
}

func (this *TripUpdate) Reset()         { *this = TripUpdate{} }
func (this *TripUpdate) String() string { return proto.CompactTextString(this) }

type TripUpdate_StopTimeEvent struct {
	Delay            *int32 `protobuf:"varint,1,opt,name=delay" json:"delay"`
	Time             *int64 `protobuf:"varint,2,opt,name=time" json:"time"`
	Uncertainty      *int32 `protobuf:"varint,3,opt,name=uncertainty" json:"uncertainty"`
	XXX_unrecognized []byte
}

func (this *TripUpdate_StopTimeEvent) Reset()         { *this = TripUpdate_StopTimeEvent{} }
func (this *TripUpdate_StopTimeEvent) String() string { return proto.CompactTextString(this) }

type TripUpdate_StopTimeUpdate struct {
	StopSequence         *uint32                                         `protobuf:"varint,1,opt,name=stop_sequence" json:"stop_sequence"`
	StopId               *string                                         `protobuf:"bytes,4,opt,name=stop_id" json:"stop_id"`
	Arrival              *TripUpdate_StopTimeEvent                       `protobuf:"bytes,2,opt,name=arrival" json:"arrival"`
	Departure            *TripUpdate_StopTimeEvent                       `protobuf:"bytes,3,opt,name=departure" json:"departure"`
	ScheduleRelationship *TripUpdate_StopTimeUpdate_ScheduleRelationship `protobuf:"varint,5,opt,name=schedule_relationship,enum=transit_realtime.TripUpdate_StopTimeUpdate_ScheduleRelationship,def=0" json:"schedule_relationship"`
	XXX_unrecognized     []byte
}

func (this *TripUpdate_StopTimeUpdate) Reset()         { *this = TripUpdate_StopTimeUpdate{} }
func (this *TripUpdate_StopTimeUpdate) String() string { return proto.CompactTextString(this) }

const Default_TripUpdate_StopTimeUpdate_ScheduleRelationship TripUpdate_StopTimeUpdate_ScheduleRelationship = TripUpdate_StopTimeUpdate_SCHEDULED

type VehiclePosition struct {
	Trip                *TripDescriptor                    `protobuf:"bytes,1,opt,name=trip" json:"trip"`
	Vehicle             *VehicleDescriptor                 `protobuf:"bytes,8,opt,name=vehicle" json:"vehicle"`
	Position            *Position                          `protobuf:"bytes,2,opt,name=position" json:"position"`
	CurrentStopSequence *uint32                            `protobuf:"varint,3,opt,name=current_stop_sequence" json:"current_stop_sequence"`
	StopId              *string                            `protobuf:"bytes,7,opt,name=stop_id" json:"stop_id"`
	CurrentStatus       *VehiclePosition_VehicleStopStatus `protobuf:"varint,4,opt,name=current_status,enum=transit_realtime.VehiclePosition_VehicleStopStatus,def=2" json:"current_status"`
	Timestamp           *uint64                            `protobuf:"varint,5,opt,name=timestamp" json:"timestamp"`
	CongestionLevel     *VehiclePosition_CongestionLevel   `protobuf:"varint,6,opt,name=congestion_level,enum=transit_realtime.VehiclePosition_CongestionLevel" json:"congestion_level"`
	XXX_unrecognized    []byte
}

func (this *VehiclePosition) Reset()         { *this = VehiclePosition{} }
func (this *VehiclePosition) String() string { return proto.CompactTextString(this) }

const Default_VehiclePosition_CurrentStatus VehiclePosition_VehicleStopStatus = VehiclePosition_IN_TRANSIT_TO

type Alert struct {
	ActivePeriod     []*TimeRange      `protobuf:"bytes,1,rep,name=active_period" json:"active_period"`
	InformedEntity   []*EntitySelector `protobuf:"bytes,5,rep,name=informed_entity" json:"informed_entity"`
	Cause            *Alert_Cause      `protobuf:"varint,6,opt,name=cause,enum=transit_realtime.Alert_Cause,def=1" json:"cause"`
	Effect           *Alert_Effect     `protobuf:"varint,7,opt,name=effect,enum=transit_realtime.Alert_Effect,def=8" json:"effect"`
	Url              *TranslatedString `protobuf:"bytes,8,opt,name=url" json:"url"`
	HeaderText       *TranslatedString `protobuf:"bytes,10,opt,name=header_text" json:"header_text"`
	DescriptionText  *TranslatedString `protobuf:"bytes,11,opt,name=description_text" json:"description_text"`
	XXX_unrecognized []byte
}

func (this *Alert) Reset()         { *this = Alert{} }
func (this *Alert) String() string { return proto.CompactTextString(this) }

const Default_Alert_Cause Alert_Cause = Alert_UNKNOWN_CAUSE
const Default_Alert_Effect Alert_Effect = Alert_UNKNOWN_EFFECT

type TimeRange struct {
	Start            *uint64 `protobuf:"varint,1,opt,name=start" json:"start"`
	End              *uint64 `protobuf:"varint,2,opt,name=end" json:"end"`
	XXX_unrecognized []byte
}

func (this *TimeRange) Reset()         { *this = TimeRange{} }
func (this *TimeRange) String() string { return proto.CompactTextString(this) }

type Position struct {
	Latitude         *float32 `protobuf:"fixed32,1,req,name=latitude" json:"latitude"`
	Longitude        *float32 `protobuf:"fixed32,2,req,name=longitude" json:"longitude"`
	Bearing          *float32 `protobuf:"fixed32,3,opt,name=bearing" json:"bearing"`
	Odometer         *float64 `protobuf:"fixed64,4,opt,name=odometer" json:"odometer"`
	Speed            *float32 `protobuf:"fixed32,5,opt,name=speed" json:"speed"`
	XXX_unrecognized []byte
}

func (this *Position) Reset()         { *this = Position{} }
func (this *Position) String() string { return proto.CompactTextString(this) }

type TripDescriptor struct {
	TripId               *string                              `protobuf:"bytes,1,opt,name=trip_id" json:"trip_id"`
	RouteId              *string                              `protobuf:"bytes,5,opt,name=route_id" json:"route_id"`
	StartTime            *string                              `protobuf:"bytes,2,opt,name=start_time" json:"start_time"`
	StartDate            *string                              `protobuf:"bytes,3,opt,name=start_date" json:"start_date"`
	ScheduleRelationship *TripDescriptor_ScheduleRelationship `protobuf:"varint,4,opt,name=schedule_relationship,enum=transit_realtime.TripDescriptor_ScheduleRelationship" json:"schedule_relationship"`
	XXX_unrecognized     []byte
}

func (this *TripDescriptor) Reset()         { *this = TripDescriptor{} }
func (this *TripDescriptor) String() string { return proto.CompactTextString(this) }

type VehicleDescriptor struct {
	Id               *string `protobuf:"bytes,1,opt,name=id" json:"id"`
	Label            *string `protobuf:"bytes,2,opt,name=label" json:"label"`
	LicensePlate     *string `protobuf:"bytes,3,opt,name=license_plate" json:"license_plate"`
	XXX_unrecognized []byte
}

func (this *VehicleDescriptor) Reset()         { *this = VehicleDescriptor{} }
func (this *VehicleDescriptor) String() string { return proto.CompactTextString(this) }

type EntitySelector struct {
	AgencyId         *string         `protobuf:"bytes,1,opt,name=agency_id" json:"agency_id"`
	RouteId          *string         `protobuf:"bytes,2,opt,name=route_id" json:"route_id"`
	RouteType        *int32          `protobuf:"varint,3,opt,name=route_type" json:"route_type"`
	Trip             *TripDescriptor `protobuf:"bytes,4,opt,name=trip" json:"trip"`
	StopId           *string         `protobuf:"bytes,5,opt,name=stop_id" json:"stop_id"`
	XXX_unrecognized []byte
}

func (this *EntitySelector) Reset()         { *this = EntitySelector{} }
func (this *EntitySelector) String() string { return proto.CompactTextString(this) }

type TranslatedString struct {
	Translation      []*TranslatedString_Translation `protobuf:"bytes,1,rep,name=translation" json:"translation"`
	XXX_unrecognized []byte
}

func (this *TranslatedString) Reset()         { *this = TranslatedString{} }
func (this *TranslatedString) String() string { return proto.CompactTextString(this) }

type TranslatedString_Translation struct {
	Text             *string `protobuf:"bytes,1,req,name=text" json:"text"`
	Language         *string `protobuf:"bytes,2,opt,name=language" json:"language"`
	XXX_unrecognized []byte
}

func (this *TranslatedString_Translation) Reset()         { *this = TranslatedString_Translation{} }
func (this *TranslatedString_Translation) String() string { return proto.CompactTextString(this) }

func init() {
	proto.RegisterEnum("transit_realtime.FeedHeader_Incrementality", FeedHeader_Incrementality_name, FeedHeader_Incrementality_value)
	proto.RegisterEnum("transit_realtime.TripUpdate_StopTimeUpdate_ScheduleRelationship", TripUpdate_StopTimeUpdate_ScheduleRelationship_name, TripUpdate_StopTimeUpdate_ScheduleRelationship_value)
	proto.RegisterEnum("transit_realtime.VehiclePosition_VehicleStopStatus", VehiclePosition_VehicleStopStatus_name, VehiclePosition_VehicleStopStatus_value)
	proto.RegisterEnum("transit_realtime.VehiclePosition_CongestionLevel", VehiclePosition_CongestionLevel_name, VehiclePosition_CongestionLevel_value)
	proto.RegisterEnum("transit_realtime.Alert_Cause", Alert_Cause_name, Alert_Cause_value)
	proto.RegisterEnum("transit_realtime.Alert_Effect", Alert_Effect_name, Alert_Effect_value)
	proto.RegisterEnum("transit_realtime.TripDescriptor_ScheduleRelationship", TripDescriptor_ScheduleRelationship_name, TripDescriptor_ScheduleRelationship_value)
}
