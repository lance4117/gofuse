package errs

import "fmt"

var (
	WarnRegisterModulesEmpty = fmt.Sprintf("Warn RegisterModules is empty, defult regist bank and auth")
)

func WarnConfigLoadNil(key string) string {
	return fmt.Sprintf("Warn Config %s is Empty", key)
}
