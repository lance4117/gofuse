package times

import (
	"testing"
)

func TestTime(t *testing.T) {
	now := NowMilli()
	t.Log(now)
	nowafter := NowAfter(-Hour)
	t.Log(nowafter)
	t.Log(ToDateTime(now))
	t.Log(ToDateTime(nowafter))
	t.Log(ToDateOnly(nowafter))

	t.Log(NowDateTime())
	t.Log(NowDateOnly())

}
