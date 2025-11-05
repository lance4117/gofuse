package gen

import (
	"math/rand/v2"

	"github.com/brianvoe/gofakeit/v7"
)

// IntN 生成 [0,n) 的随机整数；n<=0 时返回 0
func IntN(max int) int {
	if max <= 0 {
		return 0
	}
	return rand.IntN(max)
}

// IntRange 生成 [min, max] 的整数，若 max<min 自动互换
func IntRange(min, max int) int { return gofakeit.Number(min, max) }

// FloatRange 生成 [min, max] 的浮点数
func FloatRange(min, max float64) float64 { return gofakeit.Float64Range(min, max) }

// Shuffle 对切片做原地洗牌
func Shuffle[T any](arr []T) { gofakeit.ShuffleAnySlice(arr) }
