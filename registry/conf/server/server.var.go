package server

import (
	"fmt"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/conf"
)

var _ conf.IVarConf = &VarConf{}

//VarConf 变量信息
type VarConf struct {
	varConfPath  string
	varVersion   int32
	varNodeConfs map[string]conf.JSONConf
	registry     registry.IRegistry
}

//NewVarConf 构建服务器配置缓存
func NewVarConf(varConfPath string, rgst registry.IRegistry) (s *VarConf, err error) {
	s = &VarConf{
		varConfPath:  varConfPath,
		registry:     rgst,
		varNodeConfs: make(map[string]conf.JSONConf),
	}
	if err = s.load(); err != nil {
		return
	}
	return s, nil
}

//load 加载所有配置项
func (c *VarConf) load() (err error) {

	//检查跟路径是否存在
	if b, err := c.registry.Exists(c.varConfPath); err == nil && !b {
		return nil
	}

	//获取第一级目录
	var varfirstNodes []string
	varfirstNodes, c.varVersion, err = c.registry.GetChildren(c.varConfPath)
	if err != nil {
		return err
	}

	for _, p := range varfirstNodes {

		//获取第二级目录
		firstNodePath := registry.Join(c.varConfPath, p)
		varSecondChildren, _, err := c.registry.GetChildren(firstNodePath)
		if err != nil {
			return err
		}

		//获取二级目录的值
		for _, node := range varSecondChildren {
			nodePath := registry.Join(firstNodePath, node)
			data, version, err := c.registry.GetValue(nodePath)
			if err != nil {
				return err
			}
			rdata, err := decrypt(data)
			if err != nil {
				return err
			}
			varConf, err := conf.NewJSONConf(rdata, version)
			if err != nil {
				err = fmt.Errorf("%s配置有误:%v", nodePath, err)
				return err
			}
			c.varNodeConfs[registry.Join(p, node)] = *varConf
		}
	}
	return nil
}

//GetVersion 获取数据版本号
func (c *VarConf) GetVersion() int32 {
	return c.varVersion
}

//GetConf 指定配置文件名称，获取var配置信息
func (c *VarConf) GetConf(tp string, name string) (*conf.JSONConf, error) {
	if v, ok := c.varNodeConfs[registry.Join(tp, name)]; ok {
		return &v, nil
	}
	return nil, conf.ErrNoSetting
}

//GetClone 获取配置拷贝
func (c *VarConf) GetClone() conf.IVarConf {
	s := &VarConf{
		varConfPath:  c.varConfPath,
		registry:     c.registry,
		varNodeConfs: make(map[string]conf.JSONConf),
	}
	for k, v := range c.varNodeConfs {
		s.varNodeConfs[k] = v
	}
	return s
}

//Has 是否存在配置项
func (c *VarConf) Has(tp string, name string) bool {
	_, ok := c.varNodeConfs[registry.Join(tp, name)]
	return ok
}

//Iter 迭代所有子配置
func (c *VarConf) Iter(f func(path string, conf *conf.JSONConf) bool) {
	for path, v := range c.varNodeConfs {
		if !f(path, &v) {
			break
		}
	}
}