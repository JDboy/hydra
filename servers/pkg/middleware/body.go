package middleware

import "github.com/micro-plat/hydra/servers/pkg/swap"

//Body 处理请求的body参数
func Body() swap.Handler {
	return func(ctx swap.IContext) {
		if body, ok := ctx.Request().GetBody(); ok {
			ctx.Set("__body_", body)
		}
		ctx.Next()
	}
}