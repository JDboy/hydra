package context

import (
	"testing"

	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/errs"
)

func Test_Auth_Response(t *testing.T) {

	c := &ctx.Auth{} //构建对象

	//测试参数为空
	gotNil := c.Response()
	assert.Equal(t, nil, gotNil, "参数为空")

	//测试参数不为空
	response := c.Response(1)
	assert.Equal(t, 1, response, "参数不为空")

	//测试获取response
	gotResponse := c.Response()
	assert.Equal(t, 1, gotResponse, "获取response")
}

func Test_Auth_Request(t *testing.T) {

	c := &ctx.Auth{} //构建对象

	//测试参数为空
	gotNil := c.Request()
	assert.Equal(t, nil, gotNil, "参数为空")

	//测试参数不为空
	request := c.Request(1)
	assert.Equal(t, 1, request, "参数不为空")

	//获取request
	gotRequest := c.Request()
	assert.Equal(t, 1, gotRequest, "获取request")

}

func Test_Auth_Bind(t *testing.T) {
	type result struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	type obj struct{}
	var r *obj

	tests := []struct {
		name      string
		request   interface{}
		out       interface{}
		want      interface{}
		wantError bool
		errCode   int
		errStr    string
	}{
		{name: "参数非地址", request: 1, out: map[string]string{}, wantError: true, errStr: "输入参数非指针 map"},
		{name: "request为空指针", out: &map[string]string{}, request: r, wantError: true, errCode: 401, errStr: "请求中未包含用户信息,用户未登录"},
		{name: "request为空", out: &map[string]string{}, wantError: true, errCode: 401, errStr: "请求中未包含用户信息,用户未登录"},
		{name: "request为空字符串", out: &map[string]string{}, request: "", wantError: true, errCode: 401, errStr: "请求中未包含用户信息,用户未登录"},
		{name: "request为func返回空值", request: func() interface{} { return nil }, wantError: true, out: &map[string]string{},
			errCode: 401, errStr: "请求中未包含用户信息,用户未登录"},
		{name: "request为func返回非空值", request: func() interface{} { return result{Key: "1", Value: "1"} },
			out: &result{}, want: &result{Key: "1", Value: "1"}},
		{name: "request非json字符串", request: `str`, out: &map[string]string{}, wantError: true,
			errStr: "将用户信息反序化为对象时失败:invalid character 's' looking for beginning of value"},
		{name: "request为json字符串", request: `{"key":"1","value":"1"}`, out: &result{}, want: &result{Key: "1", Value: "1"}},
		{name: "默认情况", request: map[string]string{"key": "value"}, out: &map[string]string{}, want: &map[string]string{"key": "value"}},
	}
	for _, tt := range tests {
		c := &ctx.Auth{}
		c.Request(tt.request)

		err := c.Bind(tt.out)
		assert.Equal(t, tt.wantError, err != nil, tt.name)
		if tt.wantError {
			if e, ok := err.(*errs.Error); ok {
				assert.Equal(t, tt.errCode, e.GetCode(), tt.name)
				assert.Equal(t, tt.errStr, e.Error(), tt.name)
			} else {
				assert.Equal(t, tt.errStr, err.Error(), tt.name)
			}
			continue
		}
		assert.Equal(t, tt.want, tt.out, tt.name)
	}
}
