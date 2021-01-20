package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Entity{})

	r := NewSQLRepository(db)

	r.Create(Entity{
		Name: "Test1",
		Age:  21,
	})
	r.Create(Entity{
		Name: "Test2",
		Age:  22,
	})
	r.Create(Entity{
		Name: "Test2",
		Age:  23,
	})

	r.Begin()
	e, err := r.Get(1)
	if err != nil {
		panic(err)
	}

	e.Age++

	_, err = r.Update(1, *e)

	err = r.Commit()
	if err != nil {
		r.Rollback()
		panic(err)
	}
}
