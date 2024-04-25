package pb

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

// 消息体

type Msg struct {
	msgDesc *desc.MessageDescriptor // 消息描述
}

func (m *Msg) NewMessage() {
	dynamic.NewMessage(msg)
}
