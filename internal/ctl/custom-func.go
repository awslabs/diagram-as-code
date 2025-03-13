package ctl

import tmpl "text/template"

func customSeq(count int) []int {
	result := make([]int, count)
	for i := 0; i < count; i++ {
		result[i] = i
	}
	return result
}

func customAdd(x, y int) int {
	return x + y
}

func customMul(x, y int) int {
	return x * y
}

func customMkarr(args ...interface{}) []interface{} {
	return args
}

var funcMap = tmpl.FuncMap{
	"seq": customSeq,
	"add": customAdd,
	"mul": customMul,
	"mkarr": customMkarr,
}
