package eventbus

import (
	"runtime/debug"
	"sync"

	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/logger"
	"github.com/lance4117/gofuse/once"
)

type Topic int

type Event struct {
	Topic Topic
	Data  any
}

// subscriber 包装订阅者函数和唯一ID
type subscriber struct {
	id int
	fn func(event *Event)
}

type EventBus struct {
	subs   map[Topic][]subscriber
	mu     sync.RWMutex // 使用读写锁提高并发性能
	nextID int          // 用于生成唯一订阅者ID
}

var GetEventBus = once.Do(func() *EventBus {
	return &EventBus{
		subs: make(map[Topic][]subscriber),
	}
})

// Subscribe 订阅主题，返回取消订阅函数
func (eb *EventBus) Subscribe(topic Topic, fn func(event *Event)) func() {
	eb.mu.Lock()
	id := eb.nextID
	eb.nextID++
	eb.subs[topic] = append(eb.subs[topic], subscriber{id: id, fn: fn})
	eb.mu.Unlock()

	// 返回取消订阅函数
	return func() {
		eb.mu.Lock()
		defer eb.mu.Unlock()
		subs := eb.subs[topic]
		for i, s := range subs {
			if s.id == id {
				eb.subs[topic] = append(subs[:i], subs[i+1:]...)
				return
			}
		}
	}
}

// Publish 发布事件到指定主题
func (eb *EventBus) Publish(topic Topic, data any) {
	eb.mu.RLock()
	subs := make([]subscriber, len(eb.subs[topic]))
	copy(subs, eb.subs[topic]) // 复制一份避免长时间持锁
	eb.mu.RUnlock()

	if len(subs) == 0 {
		return
	}

	event := &Event{Topic: topic, Data: data}

	var wg sync.WaitGroup
	wg.Add(len(subs))
	for _, s := range subs {
		go func(sub subscriber) {
			defer func() {
				if r := recover(); r != nil {
					// 订阅者 panic 不影响其他订阅者
					logger.Errorf("%s: topic=%d subscriber_id=%d panic=%v stack=%s",
						errs.ErrSubsPanic.Error(), topic, sub.id, r, debug.Stack())
				}
				wg.Done()
			}()
			sub.fn(event)
		}(s)
	}
	wg.Wait()
}

// PublishAsync 异步发布事件，不等待处理完成
func (eb *EventBus) PublishAsync(topic Topic, data any) {
	go eb.Publish(topic, data)
}

// Unsubscribe 取消订阅整个主题
func (eb *EventBus) Unsubscribe(topic Topic) {
	eb.mu.Lock()
	delete(eb.subs, topic)
	eb.mu.Unlock()
}

// HasSubscribers 检查主题是否有订阅者
func (eb *EventBus) HasSubscribers(topic Topic) bool {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.subs[topic]) > 0
}

// SubscriberCount 返回主题的订阅者数量
func (eb *EventBus) SubscriberCount(topic Topic) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.subs[topic])
}
