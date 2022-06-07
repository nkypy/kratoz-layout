package server

import (
	"io"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"

	"switch_data_center_go/api/common"
)

const (
	baseHttpContentType      string = "application"
	httpRequestUUIDHeader    string = "RequestUUID"
	httpRequestUUIDParamName string = "RequestUUID"
)

// http api state values
const (
	HttpRespStateOK  string = "ok"
	HttpRespStateERR string = "err"
)

// returns the content-type with base prefix.
func httpContentType(subtype string) string {
	return strings.Join([]string{baseHttpContentType, subtype}, "/")
}

// HttpRequestDecoder decodes the request body to object.
func HttpRequestDecoder(r *http.Request, v interface{}) error {
	codec, ok := http.CodecForRequest(r, "Content-Type")
	if !ok {
		return errors.BadRequest("CODEC", r.Header.Get("Content-Type"))
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	if len(data) == 0 {
		return nil
	}

	if err = codec.Unmarshal(data, v); err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	// TODO set httpRequestUUIDHeader to header
	r.Header.Set(httpRequestUUIDHeader, cast.ToString(funk.Get(v, httpRequestUUIDParamName, funk.WithAllowZero())))
	return nil
}

// HttpResponseEncoder encodes the object to the HTTP response.
func HttpResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if v == nil {
		_, err := w.Write(nil)
		return err
	}

	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(dataToHttpResp(v, r.Header.Get(httpRequestUUIDHeader)))
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", httpContentType(codec.Name()))
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// HttpErrorEncoder encodes the error to the HTTP response.
func HttpErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := errors.FromError(err)
	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(errToHttpResp(se, r.Header.Get(httpRequestUUIDHeader)))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", httpContentType(codec.Name()))
	w.WriteHeader(int(se.Code))
	_, _ = w.Write(body)
}

// kratos error 转标准 Response (客户API文档指定)
func errToHttpResp(err *errors.Error, requestUUID string) common.Resp {
	return common.Resp{
		State:       HttpRespStateERR,
		RespTime:    common.RespTime(),
		RequestUUID: requestUUID,
		Data:        nil,
		Error: &common.RespError{
			Code:     err.Code,
			Reason:   err.Reason,
			Message:  err.Message,
			Metadata: err.Metadata,
		},
	}
}

// data 转标准 Response (客户API文档指定)
func dataToHttpResp(data interface{}, requestUUID string) common.Resp {
	return common.Resp{
		State:       HttpRespStateOK,
		RespTime:    common.RespTime(),
		RequestUUID: requestUUID,
		Data:        data,
		Error:       nil,
	}
}
