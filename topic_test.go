package ali_mns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

func TestNewMNSTopicWithDecoders(t *testing.T) {
	mMNSClient := &mockMNSClient{}
	fresp := &fasthttp.Response{}
	fresp.SetStatusCode(400)
	fresp.SetBody([]byte(`<Error xmlns="http://mns.aliyuncs.com/doc/v1"><Code>TopicNotExist</Code><Message>test message</Message><RequestId>test-request-id</RequestId><HostId>http://{aid}.mns.cn-shanghai.aliyuncs.com</HostId></Error>`))
	mMNSClient.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fresp, nil)

	mnsTopic := NewMNSTopicWithDecoders("test-name", mMNSClient, NewAliMNSDecoderErrResp())

	resp, err := mnsTopic.PublishMessage(MessagePublishRequest{
		MessageBody: "hello world",
	})
	assert.NotNil(t, err)
	assert.Empty(t, resp.RequestID)

	errResp := err.(ErrorResponse)
	assert.Equal(t, "test-request-id", errResp.RequestId)
	assert.Equal(t, "test message", errResp.Message)
	assert.Equal(t, "http://{aid}.mns.cn-shanghai.aliyuncs.com", errResp.HostId)
	assert.Equal(t, "TopicNotExist", errResp.Code)
	assert.Equal(t, "code: TopicNotExist, message: test message, requestId: test-request-id, hostId http://{aid}.mns.cn-shanghai.aliyuncs.com", errResp.Error())
}
