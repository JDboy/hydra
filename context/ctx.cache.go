package context

import (
	"bytes"
	"context"
	"runtime"
	"strconv"
	"sync"

	"github.com/micro-plat/hydra/global"
)

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:32]
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

var ctxMap sync.Map

//Cache 将当前上下文配置保存到当前线程编号对应的缓存
func Cache(s IContext) uint64 {
	tid := getGID()
	ctxMap.LoadOrStore(tid, s)
	return tid
}

//GetContextWithDefault 获取可用的context.Context
func GetContextWithDefault() context.Context {
	if c, ok := ctxMap.Load(getGID()); ok {
		return c.(IContext).Context()
	}
	return context.WithValue(context.Background(), "X-Request-Id", global.Def.Log().GetSessionID())
}

//Current 从缓存中获取请求上下文配置
func Current() IContext {
	if c, ok := ctxMap.Load(getGID()); ok {
		return c.(IContext)
	}
	panic("未获取到当前线程的请求上下文")
}

//Del 删除当前线程的请求上下文缓存
func Del(tid uint64) {
	ctxMap.Delete(tid)
}
