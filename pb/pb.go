package pb

var _msgMgr *MsgSet

func getMsgSet() *MsgSet {
	if _msgMgr == nil {
		_msgMgr = new(MsgSet)
		_msgMgr._set = make(map[string]*msgTemp)
	}
	return _msgMgr
}

func Encode(pbName string, msg Message) ([]byte, error) {
	return getMsgSet().Encode(pbName, msg)
}

func Decode(pbName string, bs []byte) (Message, error) {
	return getMsgSet().Decode(pbName, bs)
}
