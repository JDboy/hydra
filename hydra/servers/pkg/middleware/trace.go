package middleware

import "strings"

//Trace 系统跟踪日志
func Trace() Handler {
	return func(ctx IMiddleContext) {

		enable := ctx.ServerConf().GetMainConf().GetMainConf().GetBool("trace")
		if !enable {
			ctx.Next()
			return
		}

		ctx.Response().AddSpecial("trace")
		//1. 打印请求参数
		if strings.ToLower(ctx.Request().Path().GetMethod()) != "get" {
			input, _ := ctx.Request().GetMap()
			ctx.Log().Info("> trace.request:", input)
		}

		//2. 业务处理
		ctx.Next()

		//3. 打印响应参数
		s, c := ctx.Response().GetFinalResponse()
		ctx.Log().Info("> trace.response:", s, c)

	}
}