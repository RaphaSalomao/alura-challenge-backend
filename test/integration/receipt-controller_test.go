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
	"github.com/RaphaSalomao/alura-challenge-backend/router"
	"github.com/RaphaSalomao/alura-challenge-backend/test/factory"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ReceiptControllerSuite struct {
	suite.Suite
	db     *gorm.DB
	m      *migrate.Migrate
	client http.Client
	port   string
}

func (s *ReceiptControllerSuite) SetupSuite() {
	s.Require().NoError(godotenv.Load("../../test.env"))
	s.Require().NoError(database.Connect())
	s.db = database.DB
	s.m = database.M
	s.client = http.Client{}
	s.port = os.Getenv("SRV_PORT")

	go router.HandleRequests()
	time.Sleep(2 * time.Second)
}

func (s *ReceiptControllerSuite) TearDownTest() {
	s.db.Exec("DELETE FROM receipts")
	s.db.Exec("DELETE FROM users")
}

func (s *ReceiptControllerSuite) TearDownSuite() {
	s.m.Down()
	s.db.Exec("DROP TABLE schema_migrations")
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGINT)
}

func TestReceiptControllerSuite(t *testing.T) {
	suite.Run(t, new(ReceiptControllerSuite))
}

func (s *ReceiptControllerSuite) TestCreateReceipt_Success() {
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Body: model.ReceiptRequest{
			Description: "Salary",
			Value:       3000,
			Date:        "2022-01-25T00:00:00Z",
		},
		Path:   "/budget-control/api/v1/receipt",
		Method: http.MethodPost,
		Client: s.client,
		DB:     s.db,
		Port:   s.port,
	}
	r.SaveUser()
	expect := r.Body.(model.ReceiptRequest)

	resp, err := r.DoRequest()
	s.Require().NoError(err)

	var respBody struct{ Id uuid.UUID }
	json.NewDecoder(resp.Body).Decode(&respBody)

	var got model.Receipt
	s.db.First(&got)

	// assert
	s.Require().Equal(http.StatusCreated, resp.StatusCode)
	s.Require().Equal(respBody.Id, got.Id)
	s.Require().Equal(expect.Date, got.Date)
	s.Require().Equal(expect.Value, got.Value)
	s.Require().Equal(strings.ToUpper(expect.Description), got.Description)
	s.Require().Equal(r.User.Id, got.UserId)
}

func (s *ReceiptControllerSuite) TestCreateReceiptWithSameDescriptionInTheMonth_Fail() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Body: model.ReceiptRequest{
			Description: "Salary",
			Value:       3000,
			Date:        "2022-01-25T00:00:00Z",
		},
		Method: http.MethodPost,
		Client: s.client,
		DB:     s.db,
		Path:   "/budget-control/api/v1/receipt",
		Port:   s.port,
	}
	r.SaveUser()

	receipt := model.Receipt{
		Description: "Salary",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&receipt)

	resp, err := r.DoRequest()
	s.Require().NoError(err)

	// get response
	var respBody struct {
		Error string
		Id    uuid.UUID
	}
	json.NewDecoder(resp.Body).Decode(&respBody)
	var got model.Receipt
	s.db.First(&got, respBody.Id)

	// assert
	s.Require().Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.Require().Equal("receipt already created in current month", respBody.Error)
	s.Require().Equal(receipt.Id, got.Id)
	s.Require().Equal(receipt.Description, got.Description)
}

func (s *ReceiptControllerSuite) TestFindAllReceipt_Success() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Method: http.MethodGet,
		Client: s.client,
		DB:     s.db,
		Path:   "/budget-control/api/v1/receipt",
		Port:   s.port,
	}

	// do request
	r.SaveUser()
	receipt := model.Receipt{
		Description: "SALARY",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&receipt)

	resp, err := r.DoRequest()
	s.Require().NoError(err)

	var responseBody []model.ReceiptResponse
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(1, len(responseBody))
	s.Require().Equal(receipt.Description, responseBody[0].Description)
	s.Require().Equal(receipt.Value, responseBody[0].Value)
	s.Require().Equal(receipt.Date, responseBody[0].Date)
	s.Require().Equal(receipt.Id.String(), responseBody[0].Id)
	s.Require().Equal(r.User.Id.String(), responseBody[0].UserId)
}

func (s *ReceiptControllerSuite) TestFindReceipt_Success() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Method: http.MethodGet,
		Client: s.client,
		DB:     s.db,
		Port:   s.port,
	}
	r.SaveUser()
	receipt := model.Receipt{
		Description: "SALARY",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&receipt)
	r.Path = fmt.Sprintf("/budget-control/api/v1/receipt/%s", receipt.Id.String())

	// do request
	resp, err := r.DoRequest()
	s.Require().NoError(err)

	var responseBody model.ReceiptResponse
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(receipt.Date, responseBody.Date)
	s.Require().Equal(receipt.Description, responseBody.Description)
	s.Require().Equal(receipt.Value, responseBody.Value)
	s.Require().Equal(receipt.Id.String(), responseBody.Id)
	s.Require().Equal(receipt.UserId.String(), responseBody.UserId)
}

func (s *ReceiptControllerSuite) TestUpdateReceipt_Success() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Method: http.MethodPut,
		Client: s.client,
		DB:     s.db,
		Body: model.ReceiptRequest{
			Description: "SALARY",
			Value:       3000,
			Date:        "2022-01-20T00:00:00Z",
		},
		Port: s.port,
	}
	r.SaveUser()
	receipt := model.Receipt{
		Description: "BONUS",
		Value:       3000,
		Date:        "2022-01-25T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&receipt)

	r.Path = fmt.Sprintf("/budget-control/api/v1/receipt/%s", receipt.Id.String())
	resp, err := r.DoRequest()
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, resp.StatusCode)

	var got model.Receipt
	s.db.Find(&got, receipt.Id)

	s.Require().Equal(r.Body.(model.ReceiptRequest).Date, got.Date)
	s.Require().Equal(strings.ToUpper(r.Body.(model.ReceiptRequest).Description), got.Description)
	s.Require().Equal(r.Body.(model.ReceiptRequest).Value, got.Value)
}

func (s *ReceiptControllerSuite) TestUpdateReceiptWithSameDescriptionInTheMonth_Fail() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Method: http.MethodPut,
		Client: s.client,
		DB:     s.db,
		Body: model.ReceiptRequest{
			Description: "New Description",
			Value:       3100,
			Date:        "2022-01-22T00:00:00Z",
		},
		Port:  s.port,
	}
	r.SaveUser()

	receipt := model.Receipt{
		Description: "SALARY",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&receipt)

	inMonthReceipt := model.Receipt{
		Description: "New Description",
		Value:       5000,
		Date:        "2022-01-15T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&inMonthReceipt)

	// do request
	r.Path = fmt.Sprintf("/budget-control/api/v1/receipt/%s", receipt.Id.String())
	resp, err := r.DoRequest()
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnprocessableEntity, resp.StatusCode)

	var responseBody struct {
		Error string
		Id    uuid.UUID
	}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(inMonthReceipt.Id, responseBody.Id)
	s.Require().Equal(fmt.Sprintf("receipt %s already created in current month", inMonthReceipt.Description), responseBody.Error)

	var got model.Receipt
	s.db.Find(&got, receipt.Id)

	s.Require().NotEqual(r.Body.(model.ReceiptRequest).Date, got.Date)
	s.Require().NotEqual(strings.ToUpper(r.Body.(model.ReceiptRequest).Description), got.Description)
	s.Require().NotEqual(r.Body.(model.ReceiptRequest).Value, got.Value)
}

func (s *ReceiptControllerSuite) TestDeleteReceipt_Sucess() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Method: http.MethodDelete,
		Client: s.client,
		DB:     s.db,
		Port:  s.port,
	}
	r.SaveUser()
	receipt := model.Receipt{
		Description: "SALARY",
		Value:       1000,
		Date:        "2022-01-20T00:00:00Z",
		UserId:      r.User.Id,
	}
	s.db.Create(&receipt)

	// prepare request
	r.Path = fmt.Sprintf("/budget-control/api/v1/receipt/%s", receipt.Id.String())

	// do request
	resp, err := r.DoRequest()
	s.Require().NoError(err)

	s.Require().Equal(http.StatusNoContent, resp.StatusCode)

	tx := s.db.First(&model.Receipt{}, receipt.Id)
	s.Require().Equal(tx.Error, gorm.ErrRecordNotFound)
}

func (s *ReceiptControllerSuite) TestReceiptsByPeriod_Success() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Method: http.MethodGet,
		Client: s.client,
		DB:     s.db,
		Port:  s.port,
	}
	r.SaveUser()

	s.db.Create(&[]model.Receipt{
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
	r.Path = fmt.Sprintf("/budget-control/api/v1/receipt/%s/%s", year, month)
	resp, err := r.DoRequest()
	s.Require().NoError(err)

	var responseBody []model.Receipt
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal(2, len(responseBody))
}
