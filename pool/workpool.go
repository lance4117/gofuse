package pool

import (
	"sync"

	"github.com/panjf2000/ants/v2"
)

// TaskFunc 定义任务函数
type TaskFunc func() (any, error)

// Result 表示任务执行结果
type Result struct {
	Value any
	Err   error
}

// Pool 封装的 worker pool
type Pool struct {
	pool   *ants.PoolWithFunc
	wg     sync.WaitGroup
	result chan Result
}

// New 创建一个新的 worker pool
// size: 池子大小
func New(size int) (*Pool, error) {
	wp := &Pool{}
	var err error

	wp.result = make(chan Result, size*2) // 结果缓冲

	wp.pool, err = ants.NewPoolWithFunc(size, func(task interface{}) {
		if fn, ok := task.(TaskFunc); ok {
			val, err := fn()
			if wp.result != nil {
				wp.result <- Result{Value: val, Err: err}
			}
			wp.wg.Done()
		}
	})
	if err != nil {
		return nil, err
	}

	return wp, nil
}

// Submit 提交一个任务
func (wp *Pool) Submit(task TaskFunc) error {
	wp.wg.Add(1)
	return wp.pool.Invoke(task)
}

// Results 获取结果通道
func (wp *Pool) Results() <-chan Result {
	return wp.result
}

// Wait 等待所有任务完成
func (wp *Pool) Wait() {
	wp.wg.Wait()
}

// Release 释放资源
func (wp *Pool) Release() {
	wp.pool.Release()
	close(wp.result)
}
