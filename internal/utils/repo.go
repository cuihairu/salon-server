package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Repository[T any] interface {
	Create(item *T) error
	FindByID(id uint) (*T, error)
	FindAll() ([]T, error)
	FindByField(field string, value interface{}) ([]T, error)
	FindByFields(fields map[string]interface{}) ([]T, error)
	Update(item *T) error
	Delete(id uint) error
	ExecuteInTransaction(txFunc func(tx *gorm.DB) error) error
	FindWithPaging(paging *Paging) (*Paginated[T], error)
}

type Paginated[T any] struct {
	Total int64
	List  []T
}

func (p *Paginated[T]) HasValue() bool {
	return !p.IsEmpty()
}

func (p *Paginated[T]) IsEmpty() bool {
	return p == nil || len(p.List) == 0
}

// Repository 是 Repository 的实现
type gormRepository[T any] struct {
	db    *gorm.DB
	table *gorm.DB
}

// NewRepository 创建一个新的 gormRepository
func NewRepository[T any](db *gorm.DB) Repository[T] {
	var table T
	return &gormRepository[T]{db: db, table: db.Model(&table)}
}

func (r *gormRepository[T]) Create(item *T) error {
	return r.table.Create(item).Error
}

func (r *gormRepository[T]) FindByID(id uint) (*T, error) {
	var item T
	if err := r.table.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *gormRepository[T]) FindAll() ([]T, error) {
	var items []T
	if err := r.table.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *gormRepository[T]) FindByField(field string, value interface{}) ([]T, error) {
	var items []T
	if err := r.table.Where(field+" = ?", value).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *gormRepository[T]) FindByFields(fields map[string]interface{}) ([]T, error) {
	var items []T
	if err := r.table.Where(fields).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *gormRepository[T]) Update(item *T) error {
	if item == nil {
		return fmt.Errorf("item is nil")
	}
	return r.table.Updates(item).Error
}

func (r *gormRepository[T]) Delete(id uint) error {
	var t T
	return r.table.Delete(&t, id).Error
}

func (r *gormRepository[T]) ExecuteInTransaction(txFunc func(tx *gorm.DB) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := txFunc(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *gormRepository[T]) FindWithPaging(paging *Paging) (*Paginated[T], error) {
	var items []T
	var total int64
	if err := r.table.Count(&total).Error; err != nil {
		return nil, err
	}
	if paging == nil {
		paging = NewPagingWithDefault()
	} else {
		if paging.Page <= 0 {
			paging.Page = 1
		}
		if paging.PageSize <= 0 {
			paging.PageSize = 10
		}
	}
	if err := r.table.Limit(paging.Limit()).Offset(paging.Offset()).Find(&items).Error; err != nil {
		return nil, err
	}
	return &Paginated[T]{
		Total: total,
		List:  items,
	}, nil
}

// JsonField 是一个 json 字段
type JsonField[T any] struct {
	data *T
}

func (g *JsonField[T]) SetData(data *T) {
	if g != nil && data != nil {
		g.data = data
	}
}

func (g *JsonField[T]) Scan(value interface{}) error {
	if g == nil {
		return nil
	}
	if value == nil {
		g.data = nil
		return nil
	}
	if b, ok := value.([]byte); ok {
		var t T
		if len(b) == 0 {
			g.data = &t
			return nil
		}
		str := strings.TrimSpace(string(b)) // string(b)
		if str == "{}" || str == "null" || str == "NULL" || str == "[]" {
			g.data = &t
			return nil
		}
		err := json.Unmarshal(b, &t)
		if err != nil {
			return err
		}
		g.data = &t
		return nil
	}
	return fmt.Errorf("failed to scan Json Field: %v", value)
}

func (g *JsonField[T]) Value() (driver.Value, error) {
	if g == nil {
		return nil, nil
	}
	if g.data == nil {
		var emptyValue T
		g.data = &emptyValue
	}
	return json.Marshal(g.data)
}

func (g *JsonField[T]) HasValue() bool {
	return g != nil && g.data != nil
}

func (g *JsonField[T]) Data() *T {
	return g.data
}

func NewJsonField[T any](data *T) *JsonField[T] {
	return &JsonField[T]{data: data}
}

func (g *JsonField[T]) MarshalJSON() ([]byte, error) {
	if g == nil || g.data == nil {
		return []byte("null"), nil
	}
	return json.Marshal(g.data)
}

func (g *JsonField[T]) UnmarshalJSON(data []byte) error {
	if g == nil || len(data) == 0 || string(data) == "null" {
		return nil
	}
	var t T
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}
	g.data = &t
	return nil
}

func (g *JsonField[T]) String() string {
	if g == nil || g.data == nil {
		return ""
	}
	return fmt.Sprintf("%v", g.data)
}

func (g *JsonField[T]) IsNil() bool {
	return g == nil || g.data == nil
}
