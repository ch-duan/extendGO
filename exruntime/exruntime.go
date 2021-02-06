package exruntime

import (
	"runtime"
	"strings"
)

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	fn := f.Name()
	if strings.LastIndex(fn, ".") != -1 {
		if strings.LastIndex(fn, ".")+1 <= len(fn) {
			s := fn[strings.LastIndex(fn, ".")+1:]
			if s != "" {
				return s
			}
		}
	}
	return f.Name()
}
