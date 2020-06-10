package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
)

type dispCtx struct {
	*dispatcher.Context
}

//
func (g *dispCtx) GetRouterPath() string {
	return g.Context.Request.GetService()
}
func (g *dispCtx) GetBody() io.ReadCloser {
	return nil
}
func (g *dispCtx) GetMethod() string {
	return g.Context.Request.GetMethod()
}
func (g *dispCtx) GetURL() *url.URL {
	u, _ := url.ParseRequestURI(g.Context.Request.GetService())
	return u
}
func (g *dispCtx) GetHeaders() http.Header {
	hd := http.Header{}
	for k, v := range g.Context.Request.GetHeader() {
		hd[k] = []string{v}
	}
	return hd
}
func (g *dispCtx) GetCookies() []*http.Cookie {
	return nil
}
func (g *dispCtx) GetQuery(k string) (string, bool) {
	v, ok := g.Context.Request.GetForm()[k]
	return fmt.Sprint(v), ok
}
func (g *dispCtx) GetFormValue(k string) (string, bool) {
	v, ok := g.Context.Request.GetForm()[k]
	return fmt.Sprint(v), ok
}

func (g *dispCtx) GetForm() url.Values {
	values := url.Values{}
	for k, v := range g.Context.Request.GetForm() {
		values.Set(k, fmt.Sprint(v))
	}
	return values
}

func (g *dispCtx) WStatus(s int) {
	g.Context.Writer.WriteHeader(s)
}
func (g *dispCtx) Status() int {
	return g.Context.Writer.Status()
}
func (g *dispCtx) Written() bool {
	return g.Context.Writer.Written()
}
func (g *dispCtx) WHeader(k string) string {
	return g.Context.Writer.Header().Get(k)
}
func (g *dispCtx) ClientIP() string {
	return g.Context.GetClientIP()
}
func (g *dispCtx) ContentType() string {
	return g.Context.GetHeader("Content-Type")
}
func (g *dispCtx) File(name string) error {
	ff, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	body := base64.StdEncoding.EncodeToString(ff)
	_, err = g.Writer.WriteString(body)
	return err
}
func (g *dispCtx) ShouldBind(v interface{}) error {
	js, err := json.Marshal(g.Context.Request.GetForm)
	if err != nil {
		return fmt.Errorf("ShouldBind将输入的信息转换为JSON时失败 %w", err)
	}
	return json.Unmarshal(js, v)
}