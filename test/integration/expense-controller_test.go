package test_test

import (
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
	"github.com/RaphaSalomao/alura-challenge-backend/test/factory"
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
	port  string
}

func (s *ExpenseControllerSuite) SetupSuite() {
	s.Require().NoError(godotenv.Load("../../test.env"))
	s.Require().NoError(database.Connect())
	s.db = database.DB
	s.m = database.M
	s.client = http.Client{}
	s.port = os.Getenv("SRV_PORT")

	go router.HandleRequests()
	time.Sleep(2 * time.Second)
}

func (s *ExpenseControllerSuite) TearDownTest() {
	s.db.Exec("DELETE FROM expenses")
	s.db.Exec("DELETE FROM users")
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
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Body: model.ExpenseRequest{
			Description: "Taxes",
			Value:       3000,
			Date:        "2022-01-25T00:00:00Z",
			Category:    enum.CategoryHome,
		},
		Path:   "/budget-control/api/v1/expense",
		Method: http.MethodPost,
		DB:     s.db,
		Client: s.client,
		Port:   s.port,
	}
	r.SaveUser()
	// create request
	expect := r.Body.(model.ExpenseRequest)

	// do request
	resp, err := r.DoRequest()
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
	// preapare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Body: model.ExpenseRequest{
			Description: "Taxes",
			Value:       3000,
			Date:        "2022-01-25T00:00:00Z",
		},
		Path:   "/budget-control/api/v1/expense",
		Method: http.MethodPost,
		DB:     s.db,
		Client: s.client,
		Port:  s.port,
	}
	r.SaveUser()
	// create request
	expect := r.Body.(model.ExpenseRequest)
	// do request
	resp, err := r.DoRequest()
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
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Body: model.ExpenseRequest{
			Description: "Taxes",
			Value:       3000,
			Date:        "2022-01-25T00:00:00Z",
		},
		Path:   "/budget-control/api/v1/expense",
		Method: http.MethodPost,
		DB:     s.db,
		Client: s.client,
		Port:   s.port,
	}
	r.SaveUser()
	expense := model.Expense{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&expense)

	// do request
	resp, err := r.DoRequest()
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
	s.Require().Equal("expense already created in current month", respBody.Error)
	s.Require().Equal(expense.Id, got.Id)
	s.Require().Equal(expense.Description, got.Description)
}

func (s *ExpenseControllerSuite) TestFindAllExpense_Success() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Path:   "/budget-control/api/v1/expense",
		Method: http.MethodGet,
		DB:     s.db,
		Client: s.client,
		Port:   s.port,
	}
	r.SaveUser()
	expense := model.Expense{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&expense)

	// do request
	resp, err := r.DoRequest()
	s.Require().NoError(err)

	var responseBody []model.ExpenseResponse
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(1, len(responseBody))
	s.Require().Equal(responseBody[0].Description, expense.Description)
	s.Require().Equal(enum.CategoryOthers, expense.Category)
	s.Require().Equal(responseBody[0].Value, expense.Value)
	s.Require().Equal(responseBody[0].Date, expense.Date)
	s.Require().Equal(responseBody[0].Id, expense.Id.String())
	s.Require().Equal(responseBody[0].UserId, expense.UserId.String())
}

func (s *ExpenseControllerSuite) TestFindExpense_Success() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Method: http.MethodGet,
		DB:     s.db,
		Client: s.client,
		Port:   s.port,
	}
	r.SaveUser()
	expense := model.Expense{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
		Category:    enum.CategoryHome,
		UserId:      r.User.Id,
	}
	s.db.Create(&expense)
	r.Path = fmt.Sprintf("/budget-control/api/v1/expense/%s", expense.Id.String())
	// do request
	resp, err := r.DoRequest()
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
	s.Require().Equal(expense.Id.String(), responseBody.Id)
	s.Require().Equal(expense.UserId.String(), responseBody.UserId)
}

func (s *ExpenseControllerSuite) TestUpdateExpense_Success() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Body: model.ExpenseRequest{
			Description: "Food",
			Value:       1000,
			Date:        "2022-01-01T00:00:00Z",
			Category:    enum.CategoryFood,
		},
		Method: http.MethodPut,
		DB:     s.db,
		Client: s.client,
		Port:   s.port,
	}
	r.SaveUser()

	expect := r.Body.(model.ExpenseRequest)

	expense := model.Expense{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&expense)

	// do request
	r.Path = fmt.Sprintf("/budget-control/api/v1/expense/%s", expense.Id.String())
	resp, err := r.DoRequest()
	s.Require().NoError(err)
	s.Require().Equal(resp.StatusCode, http.StatusNoContent)

	var got model.Expense
	s.db.Find(&got, expense.Id)

	s.Require().Equal(expect.Date, got.Date)
	s.Require().Equal(strings.ToUpper(expect.Description), got.Description)
	s.Require().Equal(expect.Value, got.Value)
	s.Require().Equal(expect.Category, got.Category)
}

func (s *ExpenseControllerSuite) TestUpdateExpenseWithSameDescriptionInTheMonth_Fail() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Body: model.ExpenseRequest{
			Description: "New Description",
			Value:       1000,
			Date:        "2022-01-01T00:00:00Z",
		},
		Method: http.MethodPut,
		DB:     s.db,
		Client: s.client,
		Port:   s.port,
	}
	r.SaveUser()

	expense := model.Expense{
		Description: "Taxes",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&expense)

	inMonthExpense := model.Expense{
		Description: "New Description",
		Value:       5000,
		Date:        "2022-01-15T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&inMonthExpense)
	// prepare request
	newExpense := model.ExpenseRequest{
		Description: "NEW DESCRIPTION",
		Value:       1000,
		Date:        "2022-01-01T00:00:00Z",
		Category:    enum.CategoryOthers,
	}
	// do request
	r.Path = fmt.Sprintf("/budget-control/api/v1/expense/%s", expense.Id.String())
	resp, err := r.DoRequest()
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

	s.Require().NotEqual(newExpense.Date, got.Date)
	s.Require().NotEqual(strings.ToUpper(newExpense.Description), got.Description)
	s.Require().NotEqual(newExpense.Value, got.Value)
}

func (s *ExpenseControllerSuite) TestDeleteExpense_Sucess() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Method: http.MethodDelete,
		DB:     s.db,
		Client: s.client,
		Port:   s.port,
	}
	r.SaveUser()

	expense := model.Expense{
		Description: "Taxes",
		Value:       1000,
		Date:        "2022-01-20T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&expense)

	// prepare request
	r.Path = fmt.Sprintf("/budget-control/api/v1/expense/%s", expense.Id.String())
	// do request
	resp, err := r.DoRequest()
	s.Require().NoError(err)

	s.Require().Equal(http.StatusNoContent, resp.StatusCode)

	tx := s.db.First(&model.Expense{}, expense.Id)
	s.Require().Equal(tx.Error, gorm.ErrRecordNotFound)
}

func (s *ExpenseControllerSuite) TestExpensesByPeriod_Success() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Method: http.MethodGet,
		DB:     s.db,
		Client: s.client,
		Port:   s.port,
	}
	r.SaveUser()

	s.db.Create(&[]model.Expense{
		{
			Description: "DESC1",
			Value:       1000,
			Date:        "2022-01-01T00:00:00Z",
			UserId:      r.User.Id,
		},
		{
			Description: "DESC2",
			Value:       1000,
			Date:        "2022-01-01T00:00:00Z",
			UserId:      r.User.Id,
		},
		{
			Description: "DESC1",
			Value:       1000,
			Date:        "2022-02-01T00:00:00Z",
			UserId:      r.User.Id,
		},
		{
			Description: "DESC2",
			Value:       1000,
			Date:        "2022-02-01T00:00:00Z",
			UserId:      r.User.Id,
		},
	})

	// prepare request
	year, month := "2022", "01"
	r.Path = fmt.Sprintf("/budget-control/api/v1/expense/%s/%s", year, month)
	resp, err := r.DoRequest()
	s.Require().NoError(err)

	var responseBody []model.Expense
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(2, len(responseBody))
}
