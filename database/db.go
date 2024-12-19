package database

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/guilhermeonrails/api-go-gin/models"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var Repo *repository

type repository struct {
	DB      *gorm.DB
	SqlMock sqlmock.Sqlmock
}

func NewRepo(testMock bool) {
	var (
		sqlDB   *sql.DB
		sqlMock sqlmock.Sqlmock
		err     error
	)

	if testMock {
		sqlDB, sqlMock, err = sqlmock.New()
		if err != nil {
			panic(err)
		}
	} else {
		sqlDB, err = sql.Open("postgres", "host=localhost user=postgres password=123 dbname=curso_alura port=5432 sslmode=disable")
		if err != nil {
			panic(err)
		}
	}

	db, err := gorm.Open("postgres", sqlDB)
	if err != nil {
		panic(err)
	}

	if !testMock {
		db.AutoMigrate(&models.Aluno{})
	}

	Repo = &repository{db, sqlMock}
}
