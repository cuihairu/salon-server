package utils

import (
	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(item *T) error
	FindByID(id uint) (*T, error)
	FindAll() ([]T, error)
	FindByField(field string, value interface{}) ([]T, error)
	Update(item *T) error
	Delete(id uint) error
	ExecuteInTransaction(txFunc func(tx *gorm.DB) error) error
}

// Repository 是 Repository 的实现
type gormRepository[T any] struct {
	db *gorm.DB
}

// NewRepository 创建一个新的 gormRepository
func NewRepository[T any](db *gorm.DB) Repository[T] {
	return &gormRepository[T]{db: db}
}

func (r *gormRepository[T]) Create(item *T) error {
	return r.db.Create(item).Error
}

func (r *gormRepository[T]) FindByID(id uint) (*T, error) {
	var item T
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *gormRepository[T]) FindAll() ([]T, error) {
	var items []T
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *gormRepository[T]) FindByField(field string, value interface{}) ([]T, error) {
	var items []T
	if err := r.db.Where(field+" = ?", value).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *gormRepository[T]) Update(item *T) error {
	return r.db.Save(item).Error
}

func (r *gormRepository[T]) Delete(id uint) error {
	var t T
	return r.db.Delete(&t, id).Error
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
