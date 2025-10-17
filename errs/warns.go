package errs

import "fmt"

var (
	WarnConfigLoadNil        = fmt.Sprintf("Warn Config Value is Empty")
	WarnRegisterModulesEmpty = fmt.Sprintf("Warn RegisterModules is empty, defult regist bank and auth")
)
