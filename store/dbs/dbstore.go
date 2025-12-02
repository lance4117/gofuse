package dbs

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
// 一份 Name 只对应一份配置，后续修改同名配置也不会生效
type Config struct {
	Name            string        // 数据库连接缓存标识，name=>engine
	Driver          string        // 数据库类型 (mysql, postgres, sqlite ...)
	DSN             string        // data source name (连接字符串)
	MaxOpenConns    int           // 并发量，0 在高并发下可能导致数据库连接耗尽
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接存活时长
	ShowSQL         bool          // 调试阶段打印SQL
}

// Repo 通用数据仓储结构，使用泛型支持不同实体
type Repo[T any] struct {
	engine  *xorm.Engine
	session *xorm.Session
}

// GetOrCreateDB 获取数据库实例
func GetOrCreateDB(cfg Config) (*xorm.Engine, error) {
	enginesMu.RLock()
	eng, ok := engines[cfg.Name]
	enginesMu.RUnlock()
	if ok {
		return eng, nil
	}

	enginesMu.Lock()
	defer enginesMu.Unlock()
	// double-check
	if eng, ok := engines[cfg.Name]; ok {
		return eng, nil
	}

	e, err := xorm.NewEngine(cfg.Driver, cfg.DSN)
	if err != nil {
		logger.Error(err, errs.ErrNewStoreEngineFail)
		return nil, err
	}

	e.SetMaxOpenConns(cfg.MaxOpenConns)
	e.SetMaxIdleConns(cfg.MaxIdleConns)
	e.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	e.ShowSQL(cfg.ShowSQL)

	if err := e.Ping(); err != nil {
		logger.Error(err, errs.ErrNewStoreEngineFail)
		return nil, err
	}

	engines[cfg.Name] = e
	return e, nil
}

// NewRepo 创建一个新的数据仓储实例
func NewRepo[T any](cfg Config) (*Repo[T], error) {
	eng, err := GetOrCreateDB(cfg)
	if err != nil {
		return nil, err
	}
	return &Repo[T]{engine: eng}, nil
}

// GetByID 根据ID获取单条记录
func (r *Repo[T]) GetByID(id int64) (*T, bool, error) {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			if err := sess.Close(); err != nil {
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
			if err := sess.Close(); err != nil {
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
			if err := sess.Close(); err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	var m []T
	err := sess.Limit(limit, start).Find(&m)
	return m, err
}

// Insert 新增记录
func (r *Repo[T]) Insert(m *T) error {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			if err := sess.Close(); err != nil {
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
			if err := sess.Close(); err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	_, err := sess.Where(condition).Update(m)
	return err
}

// UpdateById 按ID更新记录
func (r *Repo[T]) UpdateById(m *T, id int64) error {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			if err := sess.Close(); err != nil {
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
			if err := sess.Close(); err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	_, err := sess.Delete(m)
	return err
}

// DeleteByID 按ID删除记录
func (r *Repo[T]) DeleteByID(id int64) error {
	sess, created := r.getSession()
	if created {
		defer func(sess *xorm.Session) {
			if err := sess.Close(); err != nil {
				logger.Error(err)
			}
		}(sess)
	}
	var m T
	_, err := sess.ID(id).Delete(&m)
	return err
}

// Session 获取当前的Session
// 记得手动关闭
func (r *Repo[T]) Session() *xorm.Session {
	sess, _ := r.getSession()
	return sess
}

// DoTx 执行一个事务fn 出现错误会自动回滚
func (r *Repo[T]) DoTx(fn func(txRepo *Repo[T]) error) error {
	// 新建一个session
	// 不污染原始 Repo，避免后续 r 继续复用。
	sess := r.engine.NewSession()
	defer func(sess *xorm.Session) {
		if err := sess.Close(); err != nil {
			logger.Error(err)
		}
	}(sess)

	if err := sess.Begin(); err != nil {
		return err
	}

	// 构建新的 Repo 使用该 session
	txRepo := &Repo[T]{engine: r.engine, session: sess}

	if err := fn(txRepo); err != nil {
		_ = sess.Rollback()
		return err
	}
	return sess.Commit()
}

// getSession 确保拿到 session，必要时新建一个
func (r *Repo[T]) getSession() (*xorm.Session, bool) {
	if r.session != nil {
		return r.session, false
	}
	return r.engine.NewSession(), true
}
