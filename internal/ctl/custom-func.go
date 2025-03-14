package ctl

import tmpl "text/template"

/*
[REVIEW REQUIRED] Custom functions can pose security risks depending on their implementation, especially if the function is not pure (i.e., a function whose output is uniquely determined for each input).

For reviewers:
Determine whether the added or modified function is pure.
If it is not pure function, consider how the change could affect not only the application but also external and OS through I/O.

https://en.wikipedia.org/wiki/Pure_function
*/


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
