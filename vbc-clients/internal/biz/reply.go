package biz

import (
	"encoding/json"
	"vbc/lib"
)

// 成功的编码
const Reply_code_success = 200
const Reply_code_BadRequest = 400
const Reply_code_Unauthorized = 401
const Reply_code_internal_error = 500 // 内部出错了
const Reply_code_common_error = 1001
const Reply_code_jwt_expired = 1002             // JWT过期了
const Reply_code_data_validation_failure = 1003 // 数据验证失败

const Reply_code_waiting_password = 2000 // 需要密码访问
const Reply_code_password_error = 2001   // 密码验证失败

type Reply map[string]interface{}

const (
	Reply_key_code    = "code"
	Reply_key_message = "message"
)

func CreateReply() Reply {

	return map[string]interface{}{
		Reply_key_code:    -1,
		Reply_key_message: "Error",
	}
}

func (c Reply) CommonStrError(err string) {
	c[Reply_key_code] = Reply_code_common_error
	c[Reply_key_message] = err
}

func (c Reply) CommonError(err error) {
	c[Reply_key_code] = Reply_code_common_error
	c[Reply_key_message] = err.Error()
}

func (c Reply) InternalError(err error) {
	c[Reply_key_code] = Reply_code_internal_error
	c[Reply_key_message] = err.Error()
}

func (c Reply) Merge(typeMap lib.TypeMap) {
	if typeMap != nil {
		for k, v := range typeMap {
			c[k] = v
		}
	}
}

func (c Reply) Success() {
	c[Reply_key_code] = Reply_code_success
	c[Reply_key_message] = "Success"
}

func (c Reply) Update(code int, message string) {
	c[Reply_key_code] = code
	c[Reply_key_message] = message
}

func (c Reply) Json() []byte {
	json, _ := json.Marshal(&c)
	return json
}
