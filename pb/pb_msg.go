package pb

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

// 消息体

type msgTemp struct {
	pkg  string // proto package name
	desc *desc.MessageDescriptor
}

func (m *msgTemp) GetName() string {
	return fmt.Sprintf("%s.%s", m.pkg, m.desc.GetName())
}

func (m *msgTemp) NewMessage() *dynamic.Message {
	return dynamic.NewMessage(m.desc)
}

func (m *msgTemp) GetFieldInfo(name string) string {
	// fmt.Println(m.desc.FindFieldByName(name).IsMap())
	return m.desc.FindFieldByName(name).String()
}

type Message map[string]interface{}
