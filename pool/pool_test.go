package pool

import (
	"fmt"
	"testing"
	"time"
)

func task(val int) TaskFunc {
	return func() (any, error) {
		time.Sleep(500 * time.Millisecond)
		return fmt.Sprintf("job-%d done", val), nil
	}
}

func TestWorkerPool(t *testing.T) {
	// 创建一个 5 worker 的池，开启结果收集
	wp, _ := New(5)

	// 提交 10 个任务
	for i := 0; i < 10; i++ {
		_ = wp.Submit(task(i))
	}

	// 异步收集结果
	go func() {
		for r := range wp.Results() {
			if r.Err != nil {
				fmt.Println("task error:", r.Err)
			} else {
				fmt.Println("task result:", r.Value)
			}
		}
	}()

	// 等待完成
	wp.Wait()
	wp.Release()
}
