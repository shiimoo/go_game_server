package pb

import (
	"sync"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

type MsgSet struct {
	mu   sync.RWMutex
	_set map[string]*msgTemp
}

func (m *MsgSet) RegisterMsg(pkg string, desc *desc.MessageDescriptor) {
	msg := new(msgTemp)
	msg.pkg = pkg
	msg.desc = desc

	m.mu.Lock()
	m._set[msg.GetName()] = msg
	m.mu.Unlock()
}

func (m *MsgSet) getFieldInfo(pbName, fName string) {
	m.mu.RLock()
	msg, ok := m._set[pbName]
	m.mu.RUnlock()
	if ok {
		msg.GetFieldInfo(fName)
	}
}

func (m *MsgSet) newMessage(pbName string) *dynamic.Message {
	m.mu.RLock()
	msg, ok := m._set[pbName]
	m.mu.RUnlock()
	if ok {
		return msg.NewMessage()
	}
	return nil
}

func (m *MsgSet) Encode(pbName string, msg Message) ([]byte, error) {
	msgObj := getMsgSet().newMessage(pbName)
	for k, v := range msg {
		getMsgSet().getFieldInfo(pbName, k)
		if err := msgObj.TrySetFieldByName(k, v); err != nil {
			// fmt.Println(k, v, err, msgObj.AddRepeatedField())
			return nil, err
		}
	}
	return msgObj.Marshal()
}

func (m *MsgSet) Decode(pbName string, bs []byte) (Message, error) {
	msgObj := getMsgSet().newMessage(pbName)
	if err := msgObj.Unmarshal(bs); err != nil {
		return nil, err
	}
	msg := make(Message)
	for _, filed := range msgObj.GetKnownFields() {
		msg[filed.GetName()] = msgObj.GetFieldByName(filed.GetName())
	}
	return msg, nil
}
