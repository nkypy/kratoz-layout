package common

import (
	"time"

	"github.com/go-kratos/kratos-layout/pkg/util"
)

func RespTime() string {
	return util.FormatDateTime(time.Now())
}

// Req 通用请求参数
type Req struct {
	// 请求 UUID
	RequestUUID string `json:"RequestUUID"`
	// 请求时间
	RequestTime string `json:"RequestTime"`
}

// Resp 通用响应参数
type Resp struct {
	// 请求状态:   ok：成功，err：失败
	State string `json:"State"`
	// 响应时间
	RespTime string `json:"RespTime"`
	// 请求时的UUID  原路返回 用于前端校验匹配关系对
	RequestUUID string `json:"RequestUUID"`
	// 业务数据
	Data interface{} `json:"Data"`
	// 错误信息（State 为 err 时）
	Error *RespError `json:"Error"`
}

type RespError struct {
	// 错误码，跟 http-status 一致
	Code int32 `json:"Code"`
	// 错误原因，定义为业务判定错误码
	Reason string `json:"Reason"`
	// 错误信息，为用户可读的信息，可作为用户提示内容
	Message string `json:"Message"`
	// 错误元信息，为错误添加附加可扩展信息
	Metadata map[string]string `json:"Metadata"`
}
