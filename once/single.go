package once

import "sync"

// Do 单例模式
func Do[T any](fn func() T) func() T {
	var (
		once sync.Once
		s    T
	)
	return func() T {
		once.Do(func() {
			s = fn()
		})
		return s
	}
}

// DoWithErr 单例模式, 返回泛型和错误
func DoWithErr[T any](fn func() (T, error)) func() (T, error) {
	var (
		once sync.Once
		s    T
		err  error
	)
	return func() (T, error) {
		once.Do(func() {
			s, err = fn()
		})
		return s, err
	}
}

// DoWithParam 单例带参数模式 泛型函数
func DoWithParam[T any, P any](fn func(P) T) func(P) T {
	var (
		once sync.Once
		s    T
	)
	return func(param P) T {
		once.Do(func() {
			s = fn(param)
		})
		return s
	}
}
