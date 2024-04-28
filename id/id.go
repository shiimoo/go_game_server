package id

var defaultGen = NewIdGen(0, nil)

func Gen() int {
	return defaultGen.Gen()
}

func GenStr() string {
	return defaultGen.GenStr()
}
