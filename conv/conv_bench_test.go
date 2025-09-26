package conv

import (
	"fmt"
	"testing"
)

func BenchmarkFmtInt(b *testing.B) {
	v := 12345
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%v", v)
	}
}

func BenchmarkAnyInt(b *testing.B) {
	v := 12345
	for i := 0; i < b.N; i++ {
		_ = AnyToString(v)
	}
}
