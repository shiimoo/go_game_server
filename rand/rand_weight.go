package random

import (
	"errors"
	"math/rand"
)

var (
	ErrWeightEleListLengthZero = errors.New("eles length must > 0")
	ErrWeightEleListTotalZero  = errors.New("eles at least one of them is > 0")
)

// Ele 权重随机元素接口
type Ele interface {
	Weight() int // 获取权重值
}

func WeightN(eles []Ele, num int, un_repeat bool) []Ele {
	eLen := len(eles)
	if eLen == 0 {
		panic(ErrWeightEleListLengthZero)
	}
	if num <= 0 {
		panic("weight rand num must > 0")
	}
	if un_repeat { // 不能重复
		if eLen == num {
			return eles
		} else if eLen < num {
			panic("when un_repeat is true , eles length must result num")
		}
	}
	gears := make([]int, 0)
	total := 0
	for _, ele := range eles {
		total += ele.Weight()
		gears = append(gears, total)
	}
	if total <= 0 {
		panic(ErrWeightEleListTotalZero)
	}
	res := make([]Ele, num)
	for i := 0; i < num; i++ {
		rNum := rand.Intn(total)
		for index, count := range gears {
			if rNum <= count {
				res[i] = eles[index]
				break
			}
		}
	}
	return res
}

func WeightOne(eles []Ele) Ele {
	eLen := len(eles)
	if eLen == 0 {
		panic(ErrWeightEleListLengthZero)
	}
	gears := make([]int, 0)
	total := 0
	for _, ele := range eles {
		total += ele.Weight()
		gears = append(gears, total)
	}
	if total <= 0 {
		panic(ErrWeightEleListTotalZero)
	}
	rNum := rand.Intn(total)
	for index, count := range gears {
		if rNum < count {
			return eles[index]
		}
	}
	return nil
}
