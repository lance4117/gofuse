package eventbus

import (
	"sync/atomic"
	"testing"
	"time"
)

const (
	TopicTest Topic = iota + 1
	TopicTest2
)

func TestEventbus(t *testing.T) {
	bus := GetEventBus()
	defer bus.Unsubscribe(TopicTest)

	var counter int32

	// Subscribe 返回取消订阅函数
	cancel1 := bus.Subscribe(TopicTest, func(e *Event) {
		atomic.AddInt32(&counter, 1)
		t.Logf("subscriber1: %v", e.Data)
	})
	_ = bus.Subscribe(TopicTest, func(e *Event) {
		atomic.AddInt32(&counter, 1)
		t.Logf("subscriber2: %v", e.Data)
	})

	// 发布事件，现在只需要传 data
	bus.Publish(TopicTest, "hello world")

	if atomic.LoadInt32(&counter) != 2 {
		t.Errorf("expected 2 calls, got %d", counter)
	}

	// 测试取消单个订阅者
	cancel1()
	atomic.StoreInt32(&counter, 0)
	bus.Publish(TopicTest, "after cancel")

	if atomic.LoadInt32(&counter) != 1 {
		t.Errorf("expected 1 call after cancel, got %d", counter)
	}
}

func TestSubscriberCount(t *testing.T) {
	bus := GetEventBus()
	defer bus.Unsubscribe(TopicTest2)

	if bus.SubscriberCount(TopicTest2) != 0 {
		t.Error("expected 0 subscribers")
	}

	cancel := bus.Subscribe(TopicTest2, func(e *Event) {})
	if bus.SubscriberCount(TopicTest2) != 1 {
		t.Error("expected 1 subscriber")
	}

	cancel()
	if bus.SubscriberCount(TopicTest2) != 0 {
		t.Error("expected 0 subscribers after cancel")
	}
}

func TestPublishAsync(t *testing.T) {
	bus := GetEventBus()
	const TopicAsync Topic = 100
	defer bus.Unsubscribe(TopicAsync)

	var called int32
	bus.Subscribe(TopicAsync, func(e *Event) {
		time.Sleep(10 * time.Millisecond)
		atomic.StoreInt32(&called, 1)
	})

	// 异步发布不阻塞
	bus.PublishAsync(TopicAsync, "async data")

	// 立即检查应该还没完成
	if atomic.LoadInt32(&called) == 1 {
		t.Log("async completed immediately (may happen on fast systems)")
	}

	// 等待完成
	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Error("async handler was not called")
	}
}

func TestPanicRecovery(t *testing.T) {
	bus := GetEventBus()
	const TopicPanic Topic = 101
	defer bus.Unsubscribe(TopicPanic)

	var called int32

	// 第一个订阅者会 panic
	bus.Subscribe(TopicPanic, func(e *Event) {
		panic("test panic")
	})

	// 第二个订阅者应该正常执行
	bus.Subscribe(TopicPanic, func(e *Event) {
		atomic.StoreInt32(&called, 1)
	})

	// 不应该 panic
	bus.Publish(TopicPanic, "test")

	if atomic.LoadInt32(&called) != 1 {
		t.Error("second subscriber should be called even if first panics")
	}
}
