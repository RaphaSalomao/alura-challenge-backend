package test_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/model/enum"
	"github.com/RaphaSalomao/alura-challenge-backend/router"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ExpenseControllerSuite struct {
	suite.Suite
	db     *gorm.DB
	m      *migrate.Migrate
	client http.Client
}

func (s *ExpenseControllerSuite) SetupSuite() {
	s.Require().NoError(godotenv.Load("../../test.env"))
	s.Require().NoError(database.Connect())
	s.db = database.DB
	s.m = database.M
	s.client = http.Client{}
	go router.HandleRequests(true)
	time.Sleep(2 * time.Second)
}

func (s *ExpenseControllerSuite) TearDownTest() {
	s.db.Exec("DELETE FROM expenses")
}

func (s *ExpenseControllerSuite) TearDownSuite() {
	s.m.Down()
	s.db.Exec("DROP TABLE schema_migrations")
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGINT)
}

func TestExpenseControllerSuite(t *testing.T) {
	suite.Run(t, new(ExpenseControllerSuite))
}

func (s *ExpenseControllerSuite) TestCreateExpense_Success() {
	// create request
	expect := model.ExpenseRequest{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-25T00:00:00Z",
		Category:    enum.CategoryHome,
	}
	request, err := json.Marshal(expect)
	s.Require().NoError(err)
	requestBody := bytes.NewBuffer(request)

	// do request
	resp, err := http.Post("http://localhost:8080/budget-control/api/v1/expense", "Application/json", requestBody)
	s.Require().NoError(err)

	// get response
	var respBody struct{ Id uuid.UUID }
	json.NewDecoder(resp.Body).Decode(&respBody)
	var got model.Expense
	s.db.First(&got)

	// assert
	s.Require().Equal(http.StatusCreated, resp.StatusCode)
	s.Require().Equal(respBody.Id, got.Id)
	s.Require().Equal(expect.Date, got.Date)
	s.Require().Equal(expect.Value, got.Value)
	s.Require().Equal(expect.Category, got.Category)
	s.Require().Equal(strings.ToUpper(expect.Description), got.Description)
}

func (s *ExpenseControllerSuite) TestCreateExpenseWithoutCategory_Success() {
	// create request
	expect := model.ExpenseRequest{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-25T00:00:00Z",
	}
	request, err := json.Marshal(expect)
	s.Require().NoError(err)
	requestBody := bytes.NewBuffer(request)

	// do request
	resp, err := http.Post("http://localhost:8080/budget-control/api/v1/expense", "Application/json", requestBody)
	s.Require().NoError(err)

	// get response
	var respBody struct{ Id uuid.UUID }
	json.NewDecoder(resp.Body).Decode(&respBody)
	var got model.Expense
	s.db.First(&got)

	// assert
	s.Require().Equal(http.StatusCreated, resp.StatusCode)
	s.Require().Equal(respBody.Id, got.Id)
	s.Require().Equal(expect.Date, got.Date)
	s.Require().Equal(expect.Value, got.Value)
	s.Require().Equal(enum.CategoryOthers, got.Category)
	s.Require().Equal(strings.ToUpper(expect.Description), got.Description)
}

func (s *ExpenseControllerSuite) TestCreateExpensetWithSameDescriptionInTheMonth_Fail() {
	// prepare database
	expense := model.Expense{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
	}
	s.db.Create(&expense)
	// create request
	expect := model.ExpenseRequest{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-25T00:00:00Z",
	}
	request, err := json.Marshal(expect)
	s.Require().NoError(err)

	requestBody := bytes.NewBuffer(request)

	// do request
	resp, err := http.Post("http://localhost:8080/budget-control/api/v1/expense", "Application/json", requestBody)
	s.Require().NoError(err)

	// get response
	var respBody struct {
		Error string
		Id    uuid.UUID
	}
	json.NewDecoder(resp.Body).Decode(&respBody)
	var got model.Expense
	s.db.First(&got, respBody.Id)

	// assert
	s.Require().Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.Require().Equal(respBody.Error, "expense already created in current month")
	s.Require().Equal(expense.Id, got.Id)
	s.Require().Equal(expense.Description, got.Description)
}

func (s *ExpenseControllerSuite) TestFindAllExpense_Success() {
	// prepare database
	expense := model.Expense{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
	}
	s.db.Create(&expense)

	// do request
	resp, err := http.Get("http://localhost:8080/budget-control/api/v1/expense")
	s.Require().NoError(err)

	var responseBody []model.ExpenseResponse
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(len(responseBody), 1)
	s.Require().Equal(responseBody[0].Description, expense.Description)
	s.Require().Equal(enum.CategoryOthers, expense.Category)
}

func (s *ExpenseControllerSuite) TestFindExpense_Success() {
	// prepare database
	expense := model.Expense{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
		Category:    enum.CategoryHome,
	}
	s.db.Create(&expense)

	// do request
	url := fmt.Sprintf("http://localhost:8080/budget-control/api/v1/expense/%s", expense.Id.String())
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	s.Require().NoError(err)

	var responseBody model.ExpenseResponse
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(expense.Date, responseBody.Date)
	s.Require().Equal(expense.Description, responseBody.Description)
	s.Require().Equal(expense.Value, responseBody.Value)
	s.Require().Equal(expense.Category, responseBody.Category)
}

func (s *ExpenseControllerSuite) TestUpdateExpense_Success() {
	// prepare database
	expense := model.Expense{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
	}
	s.db.Create(&expense)

	// prepare request
	newExpense := model.ExpenseRequest{
		Description: "Food",
		Value:       1000,
		Date:        "2022-01-01T00:00:00Z",
		Category:    enum.CategoryFood,
	}

	body, err := json.Marshal(newExpense)
	s.Require().NoError(err)
	requestBody := bytes.NewBuffer(body)
	url := fmt.Sprintf("http://localhost:8080/budget-control/api/v1/expense/%s", expense.Id.String())

	// do request
	request, err := http.NewRequest(http.MethodPut, url, requestBody)
	s.Require().NoError(err)
	request.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(request)
	s.Require().NoError(err)
	s.Require().Equal(resp.StatusCode, http.StatusNoContent)

	var got model.Expense
	s.db.Find(&got, expense.Id)

	s.Require().Equal(newExpense.Date, got.Date)
	s.Require().Equal(strings.ToUpper(newExpense.Description), got.Description)
	s.Require().Equal(newExpense.Value, got.Value)
	s.Require().Equal(newExpense.Category, got.Category)
}

func (s *ExpenseControllerSuite) TestUpdateExpenseWithSameDescriptionInTheMonth_Fail() {
	// prepare database
	expense := model.Expense{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
	}
	s.db.Create(&expense)

	inMonthExpense := model.Expense{
		Description: "New Description",
		Value:       5000,
		Date:        "2022-01-15T00:00:00Z",
	}
	s.db.Create(&inMonthExpense)
	// prepare request
	newExpense := model.ExpenseRequest{
		Description: "NEW DESCRIPTION",
		Value:       1000,
		Date:        "2022-01-01T00:00:00Z",
		Category:    enum.CategoryOthers,
	}

	body, err := json.Marshal(newExpense)
	s.Require().NoError(err)
	requestBody := bytes.NewBuffer(body)
	url := fmt.Sprintf("http://localhost:8080/budget-control/api/v1/expense/%s", expense.Id.String())

	// do request
	request, err := http.NewRequest(http.MethodPut, url, requestBody)
	s.Require().NoError(err)
	request.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(request)
	s.Require().NoError(err)
	s.Require().Equal(resp.StatusCode, http.StatusUnprocessableEntity)

	var responseBody struct {
		Error string
		Id    uuid.UUID
	}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(inMonthExpense.Id, responseBody.Id)
	s.Require().Equal(fmt.Sprintf("expense %s already created in current month", inMonthExpense.Description), responseBody.Error)

	var got model.Expense
	s.db.Find(&got, expense.Id)

	s.Require().NotEqual(got.Date, newExpense.Date)
	s.Require().NotEqual(got.Description, strings.ToUpper(newExpense.Description))
	s.Require().NotEqual(got.Value, newExpense.Value)
}

func (s *ExpenseControllerSuite) TestDeleteExpense_Sucess() {
	// prepare database
	expense := model.Expense{
		Description: "Taxes",
		Value:       1000,
		Date:        "2022-01-20T00:00:00Z",
	}
	s.db.Create(&expense)

	// prepare request
	url := fmt.Sprintf("http://localhost:8080/budget-control/api/v1/expense/%s", expense.Id.String())
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	s.Require().NoError(err)

	// do request
	resp, err := s.client.Do(request)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusNoContent, resp.StatusCode)

	tx := s.db.First(&model.Expense{}, expense.Id)
	s.Require().Equal(tx.Error, gorm.ErrRecordNotFound)
}

func (s *ExpenseControllerSuite) TestExpensesByPeriod_Success() {
	// prepare database
	s.db.Create(&[]model.Expense{
		{
			Description: "DESC1",
			Value:       1000,
			Date:        "2022-01-01T00:00:00Z",
		},
		{
			Description: "DESC2",
			Value:       1000,
			Date:        "2022-01-01T00:00:00Z",
		},
		{
			Description: "DESC1",
			Value:       1000,
			Date:        "2022-02-01T00:00:00Z",
		},
		{
			Description: "DESC2",
			Value:       1000,
			Date:        "2022-02-01T00:00:00Z",
		},
	})

	// prepare request
	year, month := "2022", "01"

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/budget-control/api/v1/expense/%s/%s", year, month))
	s.Require().NoError(err)

	var responseBody []model.Expense
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(2, len(responseBody))
}
