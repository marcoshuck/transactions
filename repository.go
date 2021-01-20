package main

import (
	"gorm.io/gorm"
	"sync"
)

type TransactionManager interface {
	Begin()
	Commit() error
	Rollback() error
}

type Entity struct {
	gorm.Model
	Name string
	Age  int
}

type Repository interface {
	TransactionManager
	Get(id int) (*Entity, error)
	Create(entity Entity) (*Entity, error)
	Update(id int, entity Entity) (*Entity, error)
}

type repository struct {
	db     *gorm.DB
	tx     *gorm.DB
	txLock sync.RWMutex
}

func (r *repository) Begin() {
	r.txLock.Lock()
	r.tx = r.db.Begin()
}

func (r *repository) Commit() error {
	defer r.txLock.Unlock()
	err := r.tx.Commit().Error
	r.tx = nil
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Rollback() error {
	defer r.txLock.Unlock()
	err := r.tx.Rollback().Error
	r.tx = nil
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Get(id int) (*Entity, error) {
	var q *gorm.DB
	if r.tx != nil {
		q = r.tx
	} else {
		q = r.db
	}

	var e Entity
	err := q.Model(&Entity{}).Where("id = ?", id).First(&e).Error
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) Create(entity Entity) (*Entity, error) {
	var q *gorm.DB
	if r.tx != nil {
		q = r.tx
	} else {
		q = r.db
	}

	err := q.Model(&Entity{}).Create(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *repository) Update(id int, entity Entity) (*Entity, error) {
	var q *gorm.DB
	if r.tx != nil {
		q = r.tx
	} else {
		q = r.db
	}

	err := q.Model(&Entity{}).Where("id = ?", id).Save(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func NewSQLRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}
