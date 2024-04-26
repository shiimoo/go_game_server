package random

import "math/rand"

func ProbaN(p, max, num int) int {
	if num <= 0 {
		panic("Proba rand num must > 0")
	}
	res := 0
	if p <= 0 || max <= 0 {
		return res
	}
	if p <= 0 {
		return res
	}

	for i := 0; i < num; i++ {
		if rand.Intn(max) <= p {
			res += 1
		}
	}
	return res
}

func ProbaOne(p, max int) bool {
	if p <= 0 || max <= 0 {
		return false
	}
	if p <= 0 {
		return false
	}
	return rand.Intn(max) <= p
}
