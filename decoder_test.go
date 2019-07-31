package ali_mns

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAliMNSDecoderErrRespDecode(t *testing.T) {
	d := NewAliMNSDecoderErrResp()

	errResp := ErrorResponse{}
	bodyReader := bytes.NewReader([]byte(`<Error xmlns="http://mns.aliyuncs.com/doc/v1"><Code>InvalidQueueName</Code><Message>test message</Message><RequestId>test-request-id</RequestId><HostId>http://{aid}.mns.cn-shanghai.aliyuncs.com</HostId></Error>`))
	err := d.Decode(bodyReader, &errResp)
	assert.Nil(t, err)

	assert.Equal(t, "test-request-id", errResp.RequestId)
	assert.Equal(t, "test message", errResp.Message)
	assert.Equal(t, "http://{aid}.mns.cn-shanghai.aliyuncs.com", errResp.HostId)
	assert.Equal(t, "InvalidQueueName", errResp.Code)
	assert.Equal(t, "code: InvalidQueueName, message: test message, requestId: test-request-id, hostId http://{aid}.mns.cn-shanghai.aliyuncs.com", errResp.Error())
}

func TestAliMNSDecoderErrRespDecodeError(t *testing.T) {
	d := NewAliMNSDecoderErrResp()

	dErr, err := d.DecodeError([]byte(`<Error xmlns="http://mns.aliyuncs.com/doc/v1"><Code>InvalidQueueName</Code><Message>test message</Message><RequestId>test-request-id</RequestId><HostId>http://{aid}.mns.cn-shanghai.aliyuncs.com</HostId></Error>`), "test-resource")
	assert.Nil(t, err)
	errResp := dErr.(ErrorResponse)
	assert.Equal(t, "test-request-id", errResp.RequestId)
	assert.Equal(t, "test message", errResp.Message)
	assert.Equal(t, "http://{aid}.mns.cn-shanghai.aliyuncs.com", errResp.HostId)
	assert.Equal(t, "InvalidQueueName", errResp.Code)
}

func TestBatchOpDecoderErrRespDecode(t *testing.T) {
	resp := BatchMessageSendRequest{}
	d := NewBatchOpDecoderErrResp(&resp)

	errResp := ErrorResponse{}
	bodyReader := bytes.NewReader([]byte(`<Error xmlns="http://mns.aliyuncs.com/doc/v1"><Code>InvalidQueueName</Code><Message>test message</Message><RequestId>test-request-id</RequestId><HostId>http://{aid}.mns.cn-shanghai.aliyuncs.com</HostId></Error>`))
	err := d.Decode(bodyReader, &errResp)
	assert.Nil(t, err)

	assert.Equal(t, "test-request-id", errResp.RequestId)
	assert.Equal(t, "test message", errResp.Message)
	assert.Equal(t, "http://{aid}.mns.cn-shanghai.aliyuncs.com", errResp.HostId)
	assert.Equal(t, "InvalidQueueName", errResp.Code)
	assert.Equal(t, "code: InvalidQueueName, message: test message, requestId: test-request-id, hostId http://{aid}.mns.cn-shanghai.aliyuncs.com", errResp.Error())
}

func TestBatchOpDecoderErrResp(t *testing.T) {
	resp := BatchMessageSendRequest{}
	d := NewBatchOpDecoderErrResp(&resp)
	dErr, err := d.DecodeError([]byte(`<Error xmlns="http://mns.aliyuncs.com/doc/v1"><Code>InvalidQueueName</Code><Message>test message</Message><RequestId>test-request-id</RequestId><HostId>http://{aid}.mns.cn-shanghai.aliyuncs.com</HostId></Error>`), "test-resource")
	assert.Nil(t, err)
	assert.Nil(t, err)
	errResp := dErr.(ErrorResponse)
	assert.Equal(t, "test-request-id", errResp.RequestId)
	assert.Equal(t, "test message", errResp.Message)
	assert.Equal(t, "http://{aid}.mns.cn-shanghai.aliyuncs.com", errResp.HostId)
	assert.Equal(t, "InvalidQueueName", errResp.Code)
}
