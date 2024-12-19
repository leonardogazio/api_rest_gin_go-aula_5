package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/guilhermeonrails/api-go-gin/controllers"
	"github.com/guilhermeonrails/api-go-gin/database"
	"github.com/guilhermeonrails/api-go-gin/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var engine *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	database.NewRepo(true) // Inicia a conexão usando o SQL Mock "github.com/DATA-DOG/go-sqlmock"
	engine = gin.Default()
}

func TestEAiFulanoBeleza(t *testing.T) {
	engine.GET("/:nome", controllers.Saudacao)
	req, _ := http.NewRequest(http.MethodGet, "/Leo", nil)
	res := httptest.NewRecorder()

	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	resBody, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, `{"API diz:":"E ai Leo, tudo beleza?"}`, string(resBody))
}

func TestListaAlunos(t *testing.T) {
	engine.GET("/alunos", controllers.ExibeTodosAlunos)
	httpReq, err := http.NewRequest(http.MethodGet, "/alunos", nil)
	require.NoError(t, err)

	httpRes := httptest.NewRecorder()

	expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "nome", "cpf", "rg"}).
		AddRow(1, time.Unix(1734566894, 0), time.Unix(1734566951, 0), nil, "Tom Araya", "12345678901", "112224449").
		AddRow(2, time.Unix(1734567947, 0), time.Unix(1734567955, 0), nil, "Marty Friedman", "11178956903", "226587930")

	expectedResult := []models.Aluno{
		{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Unix(1734566894, 0),
				UpdatedAt: time.Unix(1734566951, 0),
			},
			Nome: "Tom Araya",
			CPF:  "12345678901",
			RG:   "112224449",
		},
		{
			Model: gorm.Model{
				ID:        2,
				CreatedAt: time.Unix(1734567947, 0),
				UpdatedAt: time.Unix(1734567955, 0),
			},
			Nome: "Marty Friedman",
			CPF:  "11178956903",
			RG:   "226587930",
		},
	}

	database.Repo.SqlMock.ExpectQuery(`SELECT (.+) FROM "alunos"`).WillReturnRows(expectedRows)

	engine.ServeHTTP(httpRes, httpReq)

	assert.Equal(t, http.StatusOK, httpRes.Code)

	resBody, _ := ioutil.ReadAll(httpRes.Body)

	var res []models.Aluno
	err = json.Unmarshal(resBody, &res)
	require.NoError(t, err)

	assert.Equal(t, expectedResult, res)
}

func TestGetAlunoPeloCPF(t *testing.T) {
	tt := []struct {
		testName,
		reqCPF string
		expectedRow []driver.Value
		expectedStt int
		expectedRes models.Aluno
	}{
		{
			testName:    "Caso #1 - Tem que retornar os dados do Tom Araya.",
			reqCPF:      "12345678901",
			expectedRow: []driver.Value{1, time.Unix(1734566894, 0), time.Unix(1734566951, 0), nil, "Tom Araya", "12345678901", "112224449"},
			expectedStt: 200,
			expectedRes: models.Aluno{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Unix(1734566894, 0),
					UpdatedAt: time.Unix(1734566951, 0),
				},
				Nome: "Tom Araya",
				CPF:  "12345678901",
				RG:   "112224449",
			},
		},
		{
			testName:    "Caso #2 - Tem que retornar os dados do Marty Friedman.",
			reqCPF:      "11178956903",
			expectedRow: []driver.Value{2, time.Unix(1734567947, 0), time.Unix(1734567955, 0), nil, "Marty Friedman", "11178956903", "226587930"},
			expectedStt: 200,
			expectedRes: models.Aluno{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Unix(1734567947, 0),
					UpdatedAt: time.Unix(1734567955, 0),
				},
				Nome: "Marty Friedman",
				CPF:  "11178956903",
				RG:   "226587930",
			},
		},
		{
			testName:    "Caso #3 - Tem que retornar status HTTP 404, CPF não encontrado.",
			reqCPF:      "11111111111",
			expectedRow: nil,
			expectedStt: 404,
			expectedRes: models.Aluno{},
		},
	}

	engine.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)

	for _, testCase := range tt {
		t.Run(testCase.testName, func(t *testing.T) {
			t.Logf("Executando o teste: %s", testCase.testName)

			httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/alunos/cpf/%s", testCase.reqCPF), nil)
			require.NoError(t, err)

			httpRes := httptest.NewRecorder()

			if testCase.expectedRow != nil {
				database.Repo.SqlMock.ExpectQuery(`SELECT (.+) FROM "alunos"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "nome", "cpf", "rg"}).AddRow(testCase.expectedRow...))
			}

			engine.ServeHTTP(httpRes, httpReq)

			assert.Equal(t, testCase.expectedStt, httpRes.Code)

			resBody, _ := ioutil.ReadAll(httpRes.Body)

			var res models.Aluno
			err = json.Unmarshal(resBody, &res)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedRes, res)
		})
	}
}
