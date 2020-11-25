package run

import (
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/urfave/cli"
)

//getFlags 获取运行时的参数
func getFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, cli.StringFlag{
		Name:        "trace,t",
		Destination: &global.Def.Trace,
		Usage:       `-性能分析。支持:cpu,mem,block,mutex,web`,
	})
	flags = append(flags, cli.StringFlag{
		Name:        "tport,tp",
		Destination: &global.Def.TracePort,
		Usage:       `-性能分析服务端口号。用于trace为web模式时的端口号。默认：19999`,
	})
	flags = append(flags, global.RunCli.GetFlags()...)
	return flags
}
