package mqc

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/conf/server/queue"
)

const DefMethod = "GET"

//Request 处理任务请求
type Request struct {
	queue *queue.Queue
	mq.IMQCMessage
	method string
	form   map[string]interface{}
	header map[string]string
}

//NewRequest 构建任务请求
func NewRequest(queue *queue.Queue, m mq.IMQCMessage) (r *Request, err error) {
	r = &Request{
		IMQCMessage: m,
		queue:       queue,
		method:      DefMethod,
		form:        make(map[string]interface{}),
		header:      make(map[string]string),
	}

	//将消息原串转换为map
	input := make(map[string]interface{})
	if err = json.Unmarshal([]byte(m.GetMessage()), &input); err != nil {
		return nil, fmt.Errorf("队列%s中存放的数据不是有效的json:%s %w", queue.Queue, m.GetMessage(), err)
	}

	//检查是否包含头信息
	if v, ok := input["__header__"].(map[string]interface{}); ok {
		for n, m := range v {
			r.header[n] = fmt.Sprint(m)
		}
	}

	//将所有非"__""参数加到form列表
	for k, v := range input {
		if !strings.HasPrefix(k, "__") {
			r.form[k] = v
		}
	}

	//检查是否有专门存储数据的节点，并覆盖外部节点
	if v, ok := input["__raw__"].(map[string]interface{}); ok {
		r.form = v
	}

	r.form["__body_"] = m.GetMessage()
	return r, nil
}

//GetName 获取任务名称
func (m *Request) GetName() string {
	return m.queue.Queue
}

//GetService 服务名
func (m *Request) GetService() string {
	return m.queue.Service
}

//GetMethod 方法名
func (m *Request) GetMethod() string {
	return m.method
}

//GetForm 输入参数
func (m *Request) GetForm() map[string]interface{} {
	return m.form
}

//GetHeader 头信息
func (m *Request) GetHeader() map[string]string {
	return m.header
}
