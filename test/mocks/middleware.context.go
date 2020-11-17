package mocks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	extcontext "github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
)

var _ middleware.IMiddleContext = &MiddleContext{}

type MiddleContext struct {
	MockNext     func()
	MockMeta     conf.IMeta
	MockUser     *MockUser
	MockTFuncs   extcontext.TFuncs
	MockRequest  extcontext.IRequest
	MockResponse extcontext.IResponse
	HttpRequest  *http.Request
	HttpResponse http.ResponseWriter
	MockAPPConf  app.IAPPConf
}

func (ctx *MiddleContext) Next() {
	if ctx.MockNext != nil {
		ctx.MockNext()
	}
}

func (ctx *MiddleContext) Meta() conf.IMeta {
	return ctx.MockMeta
}

//Request 请求信息
func (ctx *MiddleContext) Request() extcontext.IRequest {
	return ctx.MockRequest
}

//Response 响应信息
func (ctx *MiddleContext) Response() extcontext.IResponse {
	return ctx.MockResponse
}

//Context 控制超时的Context
func (ctx *MiddleContext) Context() context.Context {
	return context.Background()
}

//APPConf 服务器配置
func (ctx *MiddleContext) APPConf() app.IAPPConf {
	return ctx.MockAPPConf
}

//TmplFuncs 模板函数列表
func (ctx *MiddleContext) TmplFuncs() extcontext.TFuncs {
	return ctx.MockTFuncs
}

//User 用户信息
func (ctx *MiddleContext) User() extcontext.IUser {
	return ctx.MockUser
}

//Log 日志组件
func (ctx *MiddleContext) Log() logger.ILogger {
	return logger.GetSession(ctx.MockAPPConf.GetServerConf().GetServerName(), ctx.User().GetRequestID())
}

//Close 关闭并释放资源
func (ctx *MiddleContext) Close() {}

func (ctx *MiddleContext) Trace(...interface{}) {}

//GetHttpReqResp GetHttpReqResp
func (ctx *MiddleContext) GetHttpReqResp() (*http.Request, http.ResponseWriter) {
	return ctx.HttpRequest, ctx.HttpResponse
}

var _ extcontext.IUser = &MockUser{}

type MockUser struct {
	MockClientIP  string
	MockRequestID string
	MockAuth      extcontext.IAuth
}

func (u *MockUser) GetGID() string {
	return ""
}

//GetClientIP 获取客户端请求IP
func (u *MockUser) GetClientIP() string {
	return u.MockClientIP
}

//GetRequestID 获取请求编号
func (u *MockUser) GetRequestID() string {
	return u.MockRequestID
}

//Auth 认证信息
func (u *MockUser) Auth() extcontext.IAuth {
	return u.MockAuth
}

var _ extcontext.IPath = &MockPath{}

type MockPath struct {
	MockMethod        string
	MockRequestPath   string
	MockURL           string
	MockCookies       map[string]string
	MockHeader        http.Header
	MockIsLimit       bool
	MockAllowFallback bool
	MockRouter        *router.Router
}

//GetMethod 获取服务请求方法GET POST PUT DELETE 等
func (p *MockPath) GetMethod() string {
	return p.MockMethod
}

func (p *MockPath) GetEncoding() string {
	return ""
}

//GetRouter 获取当前请求对应的路由信息
func (p *MockPath) GetRouter() (*router.Router, error) {
	if p.MockRouter == nil {
		return nil, fmt.Errorf("路由信息不存在")
	}
	return p.MockRouter, nil
}

//GetRequestPath 获取请求路径
func (p *MockPath) GetRequestPath() string {
	return p.MockRequestPath
}

//GetURL 获取请求的URL信息
func (p *MockPath) GetURL() string {
	return p.MockURL
}

//GetCookie 获取请求Cookie
func (p *MockPath) GetCookie(name string) (string, bool) {
	v, ok := p.MockCookies[name]
	return v, ok
}

//GetHeader 获取头信息
func (p *MockPath) GetHeader(name string) string {
	return p.MockHeader.Get(name)
}

//GetHeaders 获取请求头
func (p *MockPath) GetHeaders() http.Header {
	return p.MockHeader
}

//GetCookies 获取cookie信息
func (p *MockPath) GetCookies() map[string]string {
	return p.MockCookies
}

//Limit 设置限流信息
func (p *MockPath) Limit(isLimit bool, fallback bool) {
	p.MockIsLimit = isLimit
	p.MockAllowFallback = fallback
}

//IsLimited 是否已限流
func (p *MockPath) IsLimited() bool {
	return p.MockIsLimit
}

//AllowFallback 是否允许降级
func (p *MockPath) AllowFallback() bool {
	return p.MockAllowFallback
}

var _ extcontext.IRequest = &MockRequest{}

type MockRequest struct {
	SpecialList  []string
	MockPath     extcontext.IPath
	MockBindObj  interface{}
	MockParamMap map[string]string
	MockQueryMap map[string]interface{}
	MockBodyMap  map[string]interface{}
	extcontext.IGetter
	extcontext.IFile
}

//Path 地址、头、cookie相关信息
func (r *MockRequest) Path() extcontext.IPath {
	return r.MockPath
}

//Param 路由参数
func (r *MockRequest) Param(name string) string {
	return r.MockParamMap[name]
}

//Bind 将请求的参数绑定到对象
func (r *MockRequest) Bind(obj interface{}) error {
	obj = &r.MockBindObj
	return nil
}

//Check 检查指定的字段是否有值
func (r *MockRequest) Check(field ...string) error {
	if len(field) == 0 {
		return nil
	}
	for _, str := range field {
		res, ok := r.Get(str)
		if !ok || types.IsEmpty(res) {
			return fmt.Errorf("[%s]数据不存在", str)
		}
	}
	return nil
}

//GetMap 将当前请求转换为map并返回
func (r *MockRequest) GetMap() (map[string]interface{}, error) {
	if r.MockQueryMap == nil {
		return nil, fmt.Errorf("人工制造错误")
	}
	return r.MockQueryMap, nil
}

//GetRawBody 获取请求的body参数
func (r *MockRequest) GetRawBody(encoding ...string) (string, error) {
	return "", nil
}

//GetBody 获取请求的body参数
func (r *MockRequest) GetBody(encoding ...string) (string, error) {
	bytes, _ := json.Marshal(r.MockBodyMap)
	return string(bytes), nil
}

//GetBodyMap 将body转换为map
func (r *MockRequest) GetBodyMap(encoding ...string) (map[string]interface{}, error) {
	return r.MockBodyMap, nil
}

//GetBodyMap 将body转换为map
func (r *MockRequest) GetRawBodyMap(encoding ...string) (map[string]interface{}, error) {
	return nil, nil
}

//GetTrace 获取请求的trace信息
func (r *MockRequest) GetPlayload() string {
	return ""
}

//GetKeys 获取字段名称
func (r *MockRequest) GetKeys() []string {
	keys := make([]string, 0, 1)
	keyMap := map[string]string{}
	for k, _ := range r.MockParamMap {
		if _, ok := keyMap[k]; !ok {
			keyMap[k] = k
		}
	}
	for k, _ := range r.MockQueryMap {
		if _, ok := keyMap[k]; !ok {
			keyMap[k] = k
		}
	}
	for k, _ := range r.MockBodyMap {
		if _, ok := keyMap[k]; !ok {
			keyMap[k] = k
		}
	}

	for _, v := range keyMap {
		keys = append(keys, v)
	}

	return keys
}

//Get 获取字段的值
func (r *MockRequest) Get(name string) (result string, ok bool) {

	v, b := r.MockParamMap[name]
	if b {
		return fmt.Sprint(v), b
	}

	q, b := r.MockQueryMap[name]
	if b {
		return fmt.Sprint(q), b
	}

	m, b := r.MockBodyMap[name]
	if b {
		return fmt.Sprint(m), b
	}

	return "", false
}

//GetString 获取字符串
func (r *MockRequest) GetString(name string, def ...string) string {
	if v, ok := r.Get(name); ok {
		return v
	}
	return types.GetStringByIndex(def, 0, "")
}

func (r *MockRequest) GetInt(name string, def ...int) int {
	v, _ := r.Get(name)
	return types.GetInt(v, def...)
}

func (r *MockRequest) GetMax(name string, o ...int) int {
	v := r.GetInt(name, o...)
	return types.GetMax(v, o...)
}
func (r *MockRequest) GetMin(name string, o ...int) int {
	v := r.GetInt(name, o...)
	return types.GetMin(v, o...)
}
func (r *MockRequest) GetInt64(name string, def ...int64) int64 {
	v, _ := r.Get(name)
	return types.GetInt64(v, def...)
}
func (r *MockRequest) GetFloat32(name string, def ...float32) float32 {
	v, _ := r.Get(name)
	return types.GetFloat32(v, def...)
}
func (r *MockRequest) GetFloat64(name string, def ...float64) float64 {
	v, _ := r.Get(name)
	return types.GetFloat64(v, def...)
}
func (r *MockRequest) GetBool(name string, def ...bool) bool {
	v, _ := r.Get(name)
	return types.GetBool(v, def...)
}
func (r *MockRequest) GetDatetime(name string, format ...string) (time.Time, error) {
	v, _ := r.Get(name)
	return types.GetDatetime(v, format...)
}
func (r *MockRequest) IsEmpty(name string) bool {
	_, ok := r.Get(name)
	return ok
}

//SaveFile 保存上传文件到指定路径
func (r *MockRequest) SaveFile(fileKey, dst string) error {
	return nil
}

//GetFileSize 获取上传文件大小
func (r *MockRequest) GetFileSize(fileKey string) (int64, error) {
	return 0, nil
}

//GetFileName 获取上传文件名称
func (r *MockRequest) GetFileName(fileKey string) (string, error) {
	return "", nil
}

//GetFileBody 获取上传文件内容
func (r *MockRequest) GetFileBody(fileKey string) (io.ReadCloser, error) {
	return nil, nil
}

var _ extcontext.IResponse = &MockResponse{}

type MockResponse struct {
	SpecialList     []string
	MockHeader      map[string][]string
	MockRaw         interface{}
	MockStatus      int
	MockContent     string
	MockError       error
	MockContentType string
}

//AddSpecial 添加特殊标记，用于在打印响应内容时知道当前请求进行了哪些特殊处理
func (res *MockResponse) AddSpecial(t string) {
	res.SpecialList = append(res.SpecialList, t)
}

//GetSpecials 获取特殊标识字段串，多个标记用"|"分隔
func (res *MockResponse) GetSpecials() string {
	return strings.Join(res.SpecialList, "|")
}

//Header 设置响应头
func (res *MockResponse) Header(key string, val string) {
	res.MockHeader[key] = []string{val}
}

//GetHeaders 设置响应头
func (res *MockResponse) GetHeaders() map[string][]string {
	return res.MockHeader
}

//GetRaw 获取未经处理的响应内容
func (res *MockResponse) GetRaw() interface{} {
	return res.MockRaw
}

//StatusCode 设置状态码
func (res *MockResponse) StatusCode(code int) {
	res.MockStatus = code
}

//ContentType 设置Content-Type响应头
func (res *MockResponse) ContentType(v string) {
	res.MockContentType = v
	res.Header("Content-Type", v)
}

//NoNeedWrite 无需写入响应数据到缓存
func (res *MockResponse) NoNeedWrite(status int) {
	res.MockStatus = status
}

//WriteFinal 修改最终渲染内容
func (res *MockResponse) WriteFinal(status int, content string, ctp string) {
	res.MockStatus = status
	res.MockContent = content
	res.ContentType(ctp)
	return
}

//Write 向响应流中写入状态码与内容(不会立即写入)
func (res *MockResponse) Write(s int, v interface{}) error {
	res.MockStatus = s
	res.MockContent = fmt.Sprint(v)
	return nil
}

//WriteAny 向响应流中写入内容,状态码根据内容进行判断(不会立即写入)
func (res *MockResponse) WriteAny(v interface{}) error {
	switch t := v.(type) {
	case errs.IError:
		res.MockStatus = t.GetCode()
	case error:
		res.MockStatus = 500
	default:
		res.MockContent = types.GetString(v)
	}
	return nil
}

//File 向响应流中写入文件(立即写入)
func (res *MockResponse) File(path string) {
	res.MockContent = path
}

//Abort 终止当前请求继续执行
func (res *MockResponse) Abort(code int, err error) {
	res.MockStatus = code
	switch v := err.(type) {
	case errs.IError:
		res.MockContent = v.GetError().Error()
	case error:
		res.MockContent = v.Error()
	}
}

//Stop 停止当前服务执行
func (res *MockResponse) Stop(code int) {
	res.MockStatus = code
}

//GetRawResponse 获取原始响应状态码与内容
func (res *MockResponse) GetRawResponse() (int, interface{}) {
	return res.MockStatus, res.MockRaw
}

//GetFinalResponse 获取最终渲染的响应状态码与内容
func (res *MockResponse) GetFinalResponse() (int, string) {
	return res.MockStatus, res.MockContent
}

//Flush 将当前内容写入响应流
func (res *MockResponse) Flush() {
}

var _ http.ResponseWriter = &MockResponseWriter{}

type MockResponseWriter struct {
	ResponseHeader http.Header
	ContentBytes   []byte
	StatusCode     int
}

func (w *MockResponseWriter) Header() http.Header {
	return w.ResponseHeader
}
func (w *MockResponseWriter) Write(bytes []byte) (int, error) {
	w.ContentBytes = bytes
	return w.StatusCode, nil
}
func (w *MockResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

const (
	noWritten     = -1
	defaultStatus = 200
)

type MockResponseWriter2 struct {
	size   int
	status int
	data   []byte
	header http.Header
}

func (w *MockResponseWriter2) Copy() *MockResponseWriter2 {
	var cp = *w
	return &cp
}

func (w *MockResponseWriter2) Reset() {
	w.size = noWritten
	w.header = make(map[string][]string)
	w.status = defaultStatus
	w.data = nil
}

func (w *MockResponseWriter2) WriteHeader(code int) {
	if code > 0 && w.status != code {
		w.status = code
	}
}

func (w *MockResponseWriter2) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
	}
}

func (w *MockResponseWriter2) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	w.data = data
	w.size += len(data)
	return w.size, nil
}
func (w *MockResponseWriter2) Header() http.Header {
	return w.header
}
func (w *MockResponseWriter2) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	w.data = []byte(s)
	w.size += len(w.data)
	return w.size, nil
}
func (w *MockResponseWriter2) Status() int {
	return w.status
}
func (w *MockResponseWriter2) Data() []byte {
	return w.data
}
func (w *MockResponseWriter2) Size() int {
	return w.size
}

func (w *MockResponseWriter2) Written() bool {
	return w.size != noWritten
}
