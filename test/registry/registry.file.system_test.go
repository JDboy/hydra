package registry

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/registry/filesystem"
	"github.com/micro-plat/hydra/test/assert"
	r "github.com/micro-plat/lib4go/registry"
)

var fscases = []struct {
	name  string
	path  string
	value string
}{
	{name: "一段路径", path: "hydra", value: "1"},
	{name: "路径中有数字", path: "1231222", value: "2"},
	{name: "数字在前", path: "123hydra", value: "3"},
	{name: "路径中有特殊字符", path: "1232hydra#$%", value: "4"},
	{name: "路径中只有特殊字符", path: "#$%", value: "5"},
	{name: "路径中有数字", path: "/123123", value: "6"},
	{name: "带/线", path: "/hydra#$%xee", value: "7"},
	{name: "前后/", path: "/hydra#$%/", value: "8"},
	{name: "以段以上路径", path: "/hydra/abc/", value: "18"},
	{name: "多段有数字", path: "/hydra/454/", value: "17"},
	{name: "多段有特殊字符", path: "/hydra/#$#%/", value: "189"},
	{name: "较长路径", path: "/hydraabcefgjijkfsnopqrstuvwxyz", value: "1445"},
	{name: "较长分段", path: "/hydra/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/", value: "1255"},
}

func createfsRegistry() []registry.IRegistry {
	rgs := make([]registry.IRegistry, 0, 1)
	f, err := filesystem.NewFileSystem(".")
	fmt.Println(err)
	rgs = append(rgs, f)
	return rgs

}

func TestFSCreateTempNode(t *testing.T) {

	//构建所有注册中心
	rgs := createfsRegistry()

	//按注册中心进行测试
	for _, fs := range rgs {
		//创建节点
		for _, c := range fscases {
			err := fs.CreateTempNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)
		}

		//检查节点值是否正确，是否有被覆盖等
		for _, c := range fscases {
			data, v, err := fs.GetValue(c.path)
			assert.Equal(t, nil, err, c.name)
			assert.NotEqual(t, v, int32(0), c.name)
			assert.Equal(t, string(data), c.value, c.name)
		}

		//删除文件
		for _, c := range fscases {
			fs.Delete(c.path)
		}
	}

}

func TestFSUpdateNode(t *testing.T) {
	fscases := []struct {
		name   string
		path   string
		value  string
		nvalue string
	}{
		{name: "更新为数字", path: "/hydra", value: "1", nvalue: "2333"},
		{name: "更新为字符", path: "1231222", value: "2", nvalue: "sdfd"},
		{name: "更新为中文", path: "123hydra", value: "3", nvalue: "研发"},
		{name: "更新为特殊字符", path: "1232hydra#$%", value: "4", nvalue: "研发12312@#@"},
		{name: "更新为json", path: "#$%", value: "5", nvalue: `{"abc":"ef",age:[10,20]}`},
		{name: "更新为xml", path: "/hydra/apiserver/api/conf", value: "5", nvalue: `<xml><node id="abc"/></xml>`},
	}
	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, fs := range rgs {
		//创建节点,更新节点
		for _, c := range fscases {
			err := fs.CreateTempNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)

			err = fs.Update(c.path, c.nvalue)
			assert.Equal(t, nil, err, c.name)
		}

		//检查节点值是否正确
		for _, c := range fscases {
			data, v, err := fs.GetValue(c.path)
			fmt.Println(c.path, data)
			assert.Equal(t, nil, err, c.name)
			assert.NotEqual(t, v, int32(0), c.name)
			assert.Equal(t, string(data), c.nvalue, c.name)
		}
	}
}

func TestFSExists(t *testing.T) {

	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, fs := range rgs {
		for _, c := range fscases {

			//节点不存在
			exists := false
			b, err := fs.Exists(c.path)
			assert.Equal(t, nil, err, c.name)
			assert.Equal(t, b, exists, c.name)

			//创建节点
			err = fs.CreateTempNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)

			//节点应存在
			exists = true
			b, err = fs.Exists(c.path)
			assert.Equal(t, nil, err, c.name)
			assert.Equal(t, b, exists, c.name)
		}
	}
}

func TestFsDelete(t *testing.T) {

	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, fs := range rgs {
		exists := false
		for _, c := range fscases {
			//创建节点
			err := fs.CreateTempNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)

			//删除节点
			err = fs.Delete(c.path)
			assert.Equal(t, nil, err, c.name)

			//是否存在
			b, err := fs.Exists(c.path)
			assert.Equal(t, nil, err, c.name)
			assert.Equal(t, b, exists, c.name)
		}
	}
}

func TestFSChildren(t *testing.T) {

	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, fs := range rgs {
		fscases := []struct {
			name     string
			path     string
			children []string
			value    string
		}{
			{name: "一个", path: "/hydra1", value: "1", children: []string{"efg"}},
			{name: "多个", path: "/hydra2", value: "1", children: []string{"abc", "efg", "efss", "12", "!@#"}},
			{name: "空", path: "/hydra3", value: "1"},
		}

		for _, c := range fscases {

			//创建节点
			for _, ch := range c.children {
				err := fs.CreateTempNode(registry.Join(c.path, ch), c.value)
				assert.Equal(t, nil, err, c.name)
			}

			//获取子节点
			paths, v, err := fs.GetChildren(c.path)
			assert.Equal(t, nil, err, c.name)
			assert.NotEqual(t, v, 0, c.name)
			assert.Equal(t, len(paths), len(c.children), paths)

			if len(c.children) == 0 {
				continue
			}

			//排序列表
			sort.Strings(paths)
			sort.Strings(c.children)
			assert.Equal(t, paths, c.children, c.name)
		}
	}
}

func TestFSWatchValue(t *testing.T) {
	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, fs := range rgs {
		fscases := []struct {
			name   string
			path   string
			value  string
			nvalue string
		}{
			{name: "一个", path: "/hydra1", value: "1", nvalue: "2"},
			{name: "一个", path: "/hydra1/hydra1-1", value: "1", nvalue: "2"},
			{name: "一个", path: "/hydra2", value: "2", nvalue: "234"},
		}
		for _, c := range fscases {

			//创建临时节点
			err := fs.CreateTempNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)

			//监控值变化
			notify, err := fs.WatchValue(c.path)
			assert.Equal(t, nil, err, c.name)
			//此时值未变化不应收到通知
			go func(c chan r.ValueWatcher, name, nvalue string) {
				select {
				case v := <-c:
					fmt.Println("V:", v)
					value, version := v.GetValue()
					assert.NotEqual(t, version, int32(0), name)
					assert.Equal(t, nvalue, string(value), name)
				case <-time.After(time.Second):
					t.Error("测试未通过")
				}
			}(notify, c.name, c.nvalue)
		}

		//更新值
		for _, c := range fscases {
			err := fs.Update(c.path, c.nvalue)
			assert.Equal(t, nil, err, c.name)
		}

		time.Sleep(time.Second)
	}
}
