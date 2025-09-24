package dbStore

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/logger"
	"xorm.io/xorm"
)

var (
	engines   = make(map[string]*xorm.Engine)
	enginesMu sync.RWMutex
)

// Config 数据库配置结构体
// 一个 Name 只能对应一种配置，后续修改其他的内容也不会生效
type Config struct {
	Name            string        // 数据库连接缓存标识（name=>engine）
	Driver          string        // 数据库驱动名称 (mysql, postgres, sqlite ...)
	DSN             string        // data source name (连接字符串)
	MaxOpenConns    int           // 无限制：0 （在高并发场景下可能会把数据库拖死）
	MaxIdleConns    int           // 最大空闲连接数, 设置过小=>频繁建连/断连；设置过大=>浪费连接资源。
	ConnMaxLifetime time.Duration // 超过这个时间后，即使连接还在用，也会被丢弃并重新建立。MySQL 默认 8 小时
	ShowSQL         bool          // 开发阶段确认SQL 语句使用，生产环境不建议打开
}

// Repo 通用数据仓库结构体，使用泛型支持不同数据类型
type Repo[T any] struct {
	engine  *xorm.Engine
	session *xorm.Session
}

// GetOrCreateDB 获取数据库连接
func GetOrCreateDB(cfg Config) *xorm.Engine {
	enginesMu.RLock()
	eng, ok := engines[cfg.Name]
	enginesMu.RUnlock()
	if ok {
		return eng
	}

	enginesMu.Lock()
	defer enginesMu.Unlock()
	// double-check
	if eng, ok := engines[cfg.Name]; ok {
		return eng
	}

	e, err := xorm.NewEngine(cfg.Driver, cfg.DSN)
	if err != nil {
		logger.Fatal(err, errs.ErrNewStoreEngineFail)
	}

	e.SetMaxOpenConns(cfg.MaxOpenConns)
	e.SetMaxIdleConns(cfg.MaxIdleConns)
	e.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	e.ShowSQL(cfg.ShowSQL)

	if err := e.Ping(); err != nil {
		logger.Fatal(err, errs.ErrNewStoreEngineFail)
	}

	engines[cfg.Name] = e
	return e
}

// NewRepo 创建一个新的数据仓库实例
func NewRepo[T any](cfg Config) *Repo[T] {
	eng := GetOrCreateDB(cfg)
	return &Repo[T]{engine: eng}
}

// GetByID 根据ID获取单条记录
func (r *Repo[T]) GetByID(id int64) (*T, bool, error) {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			err := sess.Close()
			if err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	var m T
	has, err := sess.ID(id).Get(&m)
	return &m, has, err
}

// GetAll 获取所有记录
func (r *Repo[T]) GetAll() ([]T, error) {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			err := sess.Close()
			if err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	var m []T
	err := sess.Find(&m)
	return m, err
}

// Get 分页获取记录
func (r *Repo[T]) Get(limit, start int) ([]T, error) {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			err := sess.Close()
			if err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	var m []T
	err := sess.Limit(limit, start).Find(&m)
	return m, err
}

// Insert 插入新记录
func (r *Repo[T]) Insert(m *T) error {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			err := sess.Close()
			if err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	_, err := sess.Insert(m)
	return err
}

// Update 更新记录
func (r *Repo[T]) Update(condition string, m *T) error {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			err := sess.Close()
			if err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	_, err := sess.Where(condition).Update(m)
	return err
}

// UpdateById 根据ID更新记录
func (r *Repo[T]) UpdateById(m *T, id int64) error {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			err := sess.Close()
			if err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	_, err := sess.ID(id).Update(m)
	return err
}

// Delete 删除记录
func (r *Repo[T]) Delete(m *T) error {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			err := sess.Close()
			if err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	_, err := sess.Delete(m)
	return err
}

// DeleteByID 根据ID删除记录
func (r *Repo[T]) DeleteByID(id int64) error {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			err := sess.Close()
			if err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	var m T
	_, err := sess.ID(id).Delete(&m)
	return err
}

// Session 获取当前库的Session
// 记得手动关闭
func (r *Repo[T]) Session() *xorm.Session {
	sess, _ := r.getSession()
	return sess
}

// DoTx 执行一个事务，fn 中出现错误会自动回滚。
func (r *Repo[T]) DoTx(fn func(txRepo *Repo[T]) error) error {
	// 新建一个session
	// 不会污染原始 Repo，事务结束后 r 还能正常用。
	sess := r.engine.NewSession()
	defer func(sess *xorm.Session) {
		err := sess.Close()
		if err != nil {
			logger.Error(err)
		}
	}(sess)

	if err := sess.Begin(); err != nil {
		return err
	}

	// 事务中用一个新的 Repo，绑定 session
	txRepo := &Repo[T]{engine: r.engine, session: sess}

	if err := fn(txRepo); err != nil {
		_ = sess.Rollback()
		return err
	}
	return sess.Commit()
}

// getSession 优先返回 session，否则新建一个临时的
func (r *Repo[T]) getSession() (*xorm.Session, bool) {
	if r.session != nil {
		return r.session, false
	}
	return r.engine.NewSession(), true
}
