package ali_mns

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type NotifyStrategyType string

const (
	BACKOFF_RETRY           NotifyStrategyType = "BACKOFF_RETRY"
	EXPONENTIAL_DECAY_RETRY NotifyStrategyType = "EXPONENTIAL_DECAY_RETRY"
)

type NotifyContentFormatType string

const (
	XML        NotifyContentFormatType = "XML"
	SIMPLIFIED NotifyContentFormatType = "SIMPLIFIED"
	JSON       NotifyContentFormatType = "JSON"
	STREAM     NotifyContentFormatType = "STREAM"
)

type requestIDSetter interface {
	SetRequestID(reqID string)
}

type MessageResponse struct {
	RequestID string
	XMLName   xml.Name `xml:"Message" json:"-"`
	Code      string   `xml:"Code,omitempty" json:"code,omitempty"`
	Message   string   `xml:"Message,omitempty" json:"message,omitempty"`
	HostId    string   `xml:"HostId,omitempty" json:"host_id,omitempty"`
}

type ErrorResponse struct {
	XMLName   xml.Name `xml:"Error" json:"-"`
	Code      string   `xml:"Code,omitempty" json:"code,omitempty"`
	Message   string   `xml:"Message,omitempty" json:"message,omitempty"`
	RequestId string   `xml:"RequestId,omitempty" json:"request_id,omitempty"`
	HostId    string   `xml:"HostId,omitempty" json:"host_id,omitempty"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("code: %s, message: %s, requestId: %s, hostId %s", e.Code, e.Message, e.RequestId, e.HostId)
}

type MessageSendRequest struct {
	XMLName      xml.Name `xml:"Message" json:"-"`
	MessageBody  string   `xml:"MessageBody" json:"message_body"`
	DelaySeconds int64    `xml:"DelaySeconds" json:"delay_seconds"`
	Priority     int64    `xml:"Priority" json:"priority"`
}

type MessagePublishRequest struct {
	XMLName           xml.Name           `xml:"Message" json:"-"`
	MessageBody       string             `xml:"MessageBody" json:"message_body"`
	MessageTag        string             `xml:"MessageTag,omitempty" json:"message_tag,omitempty"`
	MessageAttributes *MessageAttributes `xml:"MessageAttributes,omitempty" json:"message_attributes,omitempty"`
}

type MessageAttributes struct {
	XMLName          xml.Name          `xml:"MessageAttributes" json:"-"`
	MailAttributes   *MailAttributes   `xml:"DirectMail,omitempty" json:"direct_mail,omitempty"`
	SMSAttributes    *SMSAttributes    `xml:"DirectSMS,omitempty" json:"direct_sms,omitempty"`
	ExtendAttributes *ExtendAttributes `xml:"ExtendAttributes,omitempty" json:"fc event,omitempty"`
}

type messageAttributesXML struct {
	XMLName          xml.Name `xml:"MessageAttributes"`
	MailAttributes   string   `xml:"DirectMail,omitempty"`
	SMSAttributes    string   `xml:"DirectSMS,omitempty"`
	ExtendAttributes string   `xml:"Extend,omitempty"`
}

func (m *MessageAttributes) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	mailAttributesStr := ""
	smsAttributeStr := ""
	extendAttributesStr := ""
	if m.MailAttributes != nil {
		mailAttributesByte, err := json.Marshal(m.MailAttributes)
		if err != nil {
			return err
		}
		mailAttributesStr = string(mailAttributesByte)
	}
	if m.SMSAttributes != nil {
		smsAttributesByte, err := json.Marshal(m.SMSAttributes)
		if err != nil {
			return err
		}
		smsAttributeStr = string(smsAttributesByte)
	}
	if m.ExtendAttributes != nil {
		extendAttributesByte, err := json.Marshal(m.ExtendAttributes)
		if err != nil {
			return err
		}
		extendAttributesStr = string(extendAttributesByte)
	}
	n := &messageAttributesXML{
		MailAttributes:   mailAttributesStr,
		SMSAttributes:    smsAttributeStr,
		ExtendAttributes: extendAttributesStr,
	}
	return e.Encode(n)
}

type MailAttributes struct {
	Subject        string `json:"Subject"`
	AccountName    string `json:"AccountName"`
	AddressType    int32  `json:"AddressType"`
	IsHtml         bool   `json:"IsHtml"`
	ReplyToAddress int32  `json:"ReplyToAddress"`
}

func (m *MailAttributes) MarshalJSON() ([]byte, error) {
	type Alias MailAttributes
	isHtml := 0
	if m.IsHtml {
		isHtml = 1
	}
	return json.Marshal(&struct {
		IsHtml int `json:"IsHtml"`
		*Alias
	}{
		IsHtml: isHtml,
		Alias:  (*Alias)(m),
	})
}

type SMSAttributes struct {
	FreeSignName string `json:"FreeSignName"`
	TemplateCode string `json:"TemplateCode"`
	Type         string `json:"Type"`
	Receiver     string `json:"Receiver"`
	SmsParams    string `json:"SmsParams"`
}

type ExtendAttributes struct {
	Context string `json:"Context"`
}

type BatchMessageSendRequest struct {
	XMLName  xml.Name             `xml:"Messages"`
	Messages []MessageSendRequest `xml:"Message"`
}

type ReceiptHandles struct {
	XMLName        xml.Name `xml:"ReceiptHandles"`
	ReceiptHandles []string `xml:"ReceiptHandle"`
}

type MessageSubsribeRequest struct {
	XMLName             xml.Name                `xml:"Subscription"`
	Endpoint            string                  `xml:"Endpoint"`
	FilterTag           string                  `xml:"FilterTag,omitempty"`
	NotifyStrategy      NotifyStrategyType      `xml:"NotifyStrategy,omitempty"`
	NotifyContentFormat NotifyContentFormatType `xml:"NotifyContentFormat,omitempty"`
	EndpointServiceName string                  `xml:"EndpointServiceName,omitempty"`
	EndpointRoleName    string                  `xml:"EndpointRoleName,omitempty"`
	EndpointCipherSet   string                  `xml:"EndpointCipherSet,omitempty"`
}

var _ requestIDSetter = &MessageSendResponse{}

type MessageSendResponse struct {
	MessageResponse
	MessageId      string `xml:"MessageId" json:"message_id"`
	MessageBodyMD5 string `xml:"MessageBodyMD5" json:"message_body_md5"`
	// ReceiptHandle is assigned when any DelayMessage is sent
	ReceiptHandle string `xml:"ReceiptHandle,omitempty"`
}

func (r *MessageSendResponse) SetRequestID(reqID string) {
	r.RequestID = reqID
}

type BatchMessageSendEntry struct {
	XMLName        xml.Name `xml:"Message" json:"-"`
	ErrorCode      string   `xml:"ErrorCode,omitempty" json:"error_code,omitempty"`
	ErrorMessage   string   `xml:"ErrorMessage,omitempty" json:"error_messages,omitempty"`
	MessageId      string   `xml:"MessageId,omitempty" json:"message_id,omitempty"`
	MessageBodyMD5 string   `xml:"MessageBodyMD5,omitempty" json:"message_body_md5,omitempty"`
}

var _ requestIDSetter = &BatchMessageSendResponse{}

type BatchMessageSendResponse struct {
	RequestID string
	XMLName   xml.Name                `xml:"Messages" json:"-"`
	Messages  []BatchMessageSendEntry `xml:"Message" json:"messages"`
}

func (r *BatchMessageSendResponse) SetRequestID(reqID string) {
	r.RequestID = reqID
}

type MessageDeleteFailEntry struct {
	XMLName       xml.Name `xml:"Error" json:"-"`
	ErrorCode     string   `xml:"ErrorCode" json:"error_code"`
	ErrorMessage  string   `xml:"ErrorMessage" json:"error_messages"`
	ReceiptHandle string   `xml:"ReceiptHandle,omitempty" json:"receipt_handle"`
}

var _ requestIDSetter = &BatchMessageDeleteErrorResponse{}

type BatchMessageDeleteErrorResponse struct {
	RequestID      string
	XMLName        xml.Name                 `xml:"Errors" json:"-"`
	FailedMessages []MessageDeleteFailEntry `xml:"Error" json:"errors"`
}

func (r *BatchMessageDeleteErrorResponse) SetRequestID(reqID string) {
	r.RequestID = reqID
}

type CreateQueueRequest struct {
	XMLName                xml.Name `xml:"Queue" json:"-"`
	DelaySeconds           int32    `xml:"DelaySeconds" json:"delay_senconds"`
	MaxMessageSize         int32    `xml:"MaximumMessageSize,omitempty" json:"maximum_message_size,omitempty"`
	MessageRetentionPeriod int32    `xml:"MessageRetentionPeriod,omitempty" json:"message_retention_period,omitempty"`
	VisibilityTimeout      int32    `xml:"VisibilityTimeout,omitempty" json:"visibility_timeout,omitempty"`
	PollingWaitSeconds     int32    `xml:"PollingWaitSeconds" json:"polling_wait_secods"`
	Slices                 int32    `xml:"Slices" json:"slices"`
}

type CreateTopicRequest struct {
	XMLName        xml.Name `xml:"Topic" json:"-"`
	MaxMessageSize int32    `xml:"MaximumMessageSize,omitempty" json:"maximum_message_size,omitempty"`
	LoggingEnabled bool     `xml:"LoggingEnabled" json:"logging_enabled"`
}

var _ requestIDSetter = &MessageReceiveResponse{}

type MessageReceiveResponse struct {
	RequestID string
	MessageResponse
	MessageId        string `xml:"MessageId" json:"message_id"`
	ReceiptHandle    string `xml:"ReceiptHandle" json:"receipt_handle"`
	MessageBodyMD5   string `xml:"MessageBodyMD5" json:"message_body_md5"`
	MessageBody      string `xml:"MessageBody" json:"message_body"`
	EnqueueTime      int64  `xml:"EnqueueTime" json:"enqueue_time"`
	NextVisibleTime  int64  `xml:"NextVisibleTime" json:"next_visible_time"`
	FirstDequeueTime int64  `xml:"FirstDequeueTime" json:"first_dequeue_time"`
	DequeueCount     int64  `xml:"DequeueCount" json:"dequeue_count"`
	Priority         int64  `xml:"Priority" json:"priority"`
}

func (r *MessageReceiveResponse) SetRequestID(reqID string) {
	r.RequestID = reqID
}

var _ requestIDSetter = &BatchMessageReceiveResponse{}

type BatchMessageReceiveResponse struct {
	RequestID string
	XMLName   xml.Name                 `xml:"Messages" json:"-"`
	Messages  []MessageReceiveResponse `xml:"Message" json:"messages"`
}

func (r *BatchMessageReceiveResponse) SetRequestID(reqID string) {
	r.RequestID = reqID
}

var _ requestIDSetter = &MessageVisibilityChangeResponse{}

type MessageVisibilityChangeResponse struct {
	RequestID       string
	XMLName         xml.Name `xml:"ChangeVisibility" json:"-"`
	ReceiptHandle   string   `xml:"ReceiptHandle" json:"receipt_handle"`
	NextVisibleTime int64    `xml:"NextVisibleTime" json:"next_visible_time"`
}

func (r *MessageVisibilityChangeResponse) SetRequestID(reqID string) {
	r.RequestID = reqID
}

type QueueAttribute struct {
	XMLName                xml.Name `xml:"Queue" json:"-"`
	QueueName              string   `xml:"QueueName,omitempty" json:"queue_name,omitempty"`
	DelaySeconds           int32    `xml:"DelaySenconds,omitempty" json:"delay_senconds,omitempty"`
	MaxMessageSize         int32    `xml:"MaximumMessageSize,omitempty" json:"maximum_message_size,omitempty"`
	MessageRetentionPeriod int32    `xml:"MessageRetentionPeriod,omitempty" json:"message_retention_period,omitempty"`
	VisibilityTimeout      int32    `xml:"VisibilityTimeout,omitempty" json:"visibility_timeout,omitempty"`
	PollingWaitSeconds     int32    `xml:"PollingWaitSeconds,omitempty" json:"polling_wait_secods,omitempty"`
	ActiveMessages         int64    `xml:"ActiveMessages,omitempty" json:"active_messages,omitempty"`
	InactiveMessages       int64    `xml:"InactiveMessages,omitempty" json:"inactive_messages,omitempty"`
	DelayMessages          int64    `xml:"DelayMessages,omitempty" json:"delay_messages,omitempty"`
	CreateTime             int64    `xml:"CreateTime,omitempty" json:"create_time,omitempty"`
	LastModifyTime         int64    `xml:"LastModifyTime,omitempty" json:"last_modify_time,omitempty"`
}

type TopicAttribute struct {
	XMLName                xml.Name `xml:"Topic" json:"-"`
	TopicName              string   `xml:"TopicName,omitempty" json:"queue_name,omitempty"`
	MaxMessageSize         int32    `xml:"MaximumMessageSize,omitempty" json:"maximum_message_size,omitempty"`
	MessageRetentionPeriod int32    `xml:"MessageRetentionPeriod,omitempty" json:"message_retention_period,omitempty"`
	MessageCount           int64    `xml:"MessageCount,omitempty" json:"message_count,omitempty"`
	CreateTime             int64    `xml:"CreateTime,omitempty" json:"create_time,omitempty"`
	LastModifyTime         int64    `xml:"LastModifyTime,omitempty" json:"last_modify_time,omitempty"`
	LoggingEnabled         bool     `xml:"LoggingEnabled" json:"logging_enabled"`
}

type SubscriptionAttribute struct {
	XMLName             xml.Name                `xml:"Subscription" json:"-"`
	SubscriptionName    string                  `xml:"SubscriptionName,omitempty" json:"queue_name,omitempty"`
	Subscriber          string                  `xml:"Subscriber,omitempty" json:"subscriber,omitempty"`
	TopicOwner          string                  `xml:"TopicOwner,omitempty" json:"topic_owner,omitempty"`
	TopicName           string                  `xml:"TopicName,omitempty" json:"topic_name,omitempty"`
	Endpoint            string                  `xml:"Endpoint,omitempty" json:"endpoint,omitempty"`
	NotifyStrategy      NotifyStrategyType      `xml:"NotifyStrategy,omitempty" json:"notify_strategy,omitempty"`
	NotifyContentFormat NotifyContentFormatType `xml:"NotifyContentFormat,omitempty" json:"notify_content_format,omitempty"`
	FilterTag           string                  `xml:"FilterTag,omitempty" json:"filter_tag,omitempty"`
	CreateTime          int64                   `xml:"CreateTime,omitempty" json:"create_time,omitempty"`
	LastModifyTime      int64                   `xml:"LastModifyTime,omitempty" json:"last_modify_time,omitempty"`
}

type SetSubscriptionAttributesRequest struct {
	XMLName        xml.Name           `xml:"Subscription" json:"-"`
	NotifyStrategy NotifyStrategyType `xml:"NotifyStrategy,omitempty" json:"notify_strategy,omitempty"`
}

type Queue struct {
	QueueURL string `xml:"QueueURL" json:"url"`
}

type Queues struct {
	XMLName    xml.Name `xml:"Queues" json:"-"`
	Queues     []Queue  `xml:"Queue" json:"queues"`
	NextMarker string   `xml:"NextMarker" json:"next_marker"`
}

type Topic struct {
	TopicURL string `xml:"TopicURL" json:"url"`
}

type Topics struct {
	XMLName    xml.Name `xml:"Topics" json:"-"`
	Topics     []Topic  `xml:"Topic" json:"topics"`
	NextMarker string   `xml:"NextMarker" json:"next_marker"`
}

type Subscription struct {
	SubscriptionURL string `xml:"SubscriptionURL" json:"url"`
}

type Subscriptions struct {
	XMLName       xml.Name       `xml:"Subscriptions" json:"-"`
	Subscriptions []Subscription `xml:"Subscription" json:"subscriptions"`
	NextMarker    string         `xml:"NextMarker" json:"next_marker"`
}
