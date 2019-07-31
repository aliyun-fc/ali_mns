package ali_mns

import (
	"bytes"
	"encoding/xml"
	"io"
)

var _ MNSDecoder = new(aliMNSDecoder)
var _ MNSDecoder = new(aliMNSDecoderErrResp)
var _ MNSDecoder = new(batchOpDecoder)
var _ MNSDecoder = new(batchOpDecoderErrResp)

type MNSDecoder interface {
	Decode(reader io.Reader, v interface{}) (err error)
	DecodeError(bodyBytes []byte, resource string) (decodedError error, err error)

	Test() bool
}

type aliMNSDecoder struct {
}

type batchOpDecoder struct {
	v interface{}
}

func NewAliMNSDecoder() MNSDecoder {
	return &aliMNSDecoder{}
}

func (p *aliMNSDecoder) Test() bool {
	return false
}

func (p *batchOpDecoder) Test() bool {
	return true
}

func (p *aliMNSDecoder) Decode(reader io.Reader, v interface{}) (err error) {
	decoder := xml.NewDecoder(reader)
	err = decoder.Decode(&v)

	return
}

func (p *aliMNSDecoder) DecodeError(bodyBytes []byte, resource string) (decodedError error, err error) {
	bodyReader := bytes.NewReader(bodyBytes)
	errResp := ErrorResponse{}
	decoder := xml.NewDecoder(bodyReader)
	err = decoder.Decode(&errResp)
	if err == nil {
		decodedError = ParseError(errResp, resource)
	}
	return
}

type BatchOpDecoderFactory func(v interface{}) MNSDecoder

func NewBatchOpDecoder(v interface{}) MNSDecoder {
	return &batchOpDecoder{v: v}
}

func (p *batchOpDecoder) Decode(reader io.Reader, v interface{}) (err error) {
	decoder := xml.NewDecoder(reader)
	err = decoder.Decode(&v)

	if err == io.EOF {
		err = nil
	}

	return
}

func (p *batchOpDecoder) DecodeError(bodyBytes []byte, resource string) (decodedError error, err error) {
	bodyReader := bytes.NewReader(bodyBytes)

	decoder := xml.NewDecoder(bodyReader)
	err = decoder.Decode(&p.v)
	if err != nil {
		bodyReader.Seek(0, 0)
		errResp := ErrorResponse{}
		err = decoder.Decode(&errResp)
		if err == nil {
			decodedError = ParseError(errResp, resource)
		}
	} else {
		decodedError = ERR_MNS_BATCH_OP_FAIL.New()
	}
	return
}

type batchOpDecoderErrResp struct {
	v interface{}
}

func NewBatchOpDecoderErrResp(v interface{}) MNSDecoder {
	return &batchOpDecoderErrResp{v: v}
}

func (p *batchOpDecoderErrResp) Test() bool {
	return false
}

func (p *batchOpDecoderErrResp) Decode(reader io.Reader, v interface{}) (err error) {
	decoder := xml.NewDecoder(reader)
	err = decoder.Decode(&v)

	if err == io.EOF {
		err = nil
	}

	return
}

func (p *batchOpDecoderErrResp) DecodeError(bodyBytes []byte, resource string) (decodedError error, err error) {
	bodyReader := bytes.NewReader(bodyBytes)

	decoder := xml.NewDecoder(bodyReader)
	err = decoder.Decode(&p.v)
	if err != nil {
		bodyReader.Seek(0, 0)
		errResp := ErrorResponse{}
		err = decoder.Decode(&errResp)
		if err == nil {
			decodedError = errResp
		}
	} else {
		decodedError = ERR_MNS_BATCH_OP_FAIL.New()
	}
	return
}

type aliMNSDecoderErrResp struct {
}

func NewAliMNSDecoderErrResp() MNSDecoder {
	return &aliMNSDecoderErrResp{}
}

func (p *aliMNSDecoderErrResp) Test() bool {
	return false
}

func (p *aliMNSDecoderErrResp) Decode(reader io.Reader, v interface{}) (err error) {
	decoder := xml.NewDecoder(reader)
	err = decoder.Decode(&v)

	return
}

func (p *aliMNSDecoderErrResp) DecodeError(bodyBytes []byte, resource string) (decodedError error, err error) {
	bodyReader := bytes.NewReader(bodyBytes)
	errResp := &ErrorResponse{}
	decoder := xml.NewDecoder(bodyReader)
	err = decoder.Decode(errResp)
	decodedError = *errResp

	return
}
