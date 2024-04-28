package id

import "strconv"

/* id生成/管理 */
type IdGen struct {
	id    int
	toStr func(int) string
}

func (i *IdGen) Gen() int {
	i.id += 1
	return i.id
}

func (i *IdGen) GenStr() string {
	if i.toStr != nil {
		return i.toStr(i.Gen())
	}
	return strconv.Itoa(i.Gen())
}

func NewIdGen(start int, toStr func(int) string) *IdGen {
	if start < 0 {
		start = 0
	}
	return &IdGen{
		id:    start,
		toStr: toStr,
	}
}
