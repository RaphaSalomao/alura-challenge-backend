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
	"github.com/RaphaSalomao/alura-challenge-backend/router"
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
}

func (s *ReceiptControllerSuite) SetupSuite() {
	s.Require().NoError(godotenv.Load("../../test.env"))
	s.Require().NoError(database.Connect())
	s.db = database.DB
	s.m = database.M
	s.client = http.Client{}

	go router.HandleRequests(true)
	time.Sleep(2 * time.Second)
}

func (s *ReceiptControllerSuite) TearDownTest() {
	s.db.Exec("DELETE FROM receipts")
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
	// create request
	expect := model.ReceiptRequest{
		Description: "Salary",
		Value:       3000,
		Date:        "2022-01-25T00:00:00Z",
	}
	request, err := json.Marshal(expect)
	s.Require().NoError(err)
	requestBody := bytes.NewBuffer(request)

	// do request
	resp, err := http.Post("http://localhost:8080/budget-control/api/v1/receipt", "Application/json", requestBody)
	s.Require().NoError(err)

	// get response
	var respBody struct{ Id uuid.UUID }
	json.NewDecoder(resp.Body).Decode(&respBody)
	var got model.Receipt
	s.db.First(&got)

	// assert
	s.Require().Equal(http.StatusCreated, resp.StatusCode)
	s.Require().Equal(got.Id, respBody.Id)
	s.Require().Equal(got.Date, expect.Date)
	s.Require().Equal(got.Value, expect.Value)
	s.Require().Equal(got.Description, strings.ToUpper(expect.Description))
}

func (s *ReceiptControllerSuite) TestCreateReceiptWithSameDescriptionInTheMonth_Fail() {
	// prepare database
	receipt := model.Receipt{
		Description: "SALARY",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
	}
	s.db.Create(&receipt)
	// create request
	expect := model.ReceiptRequest{
		Description: "Salary",
		Value:       3000,
		Date:        "2022-01-25T00:00:00Z",
	}
	request, err := json.Marshal(expect)
	s.Require().NoError(err)

	requestBody := bytes.NewBuffer(request)

	// do request
	resp, err := http.Post("http://localhost:8080/budget-control/api/v1/receipt", "Application/json", requestBody)
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
	s.Require().Equal(respBody.Error, "receipt already created in current month")
	s.Require().Equal(receipt.Id, got.Id)
	s.Require().Equal(receipt.Description, got.Description)
}

func (s *ReceiptControllerSuite) TestFindAllReceipt_Success() {
	// prepare database
	receipt := model.Receipt{
		Description: "SALARY",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
	}
	s.db.Create(&receipt)

	// do request
	resp, err := http.Get("http://localhost:8080/budget-control/api/v1/receipt")
	s.Require().NoError(err)

	var responseBody []model.ReceiptResponse
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(len(responseBody), 1)
	s.Require().Equal(responseBody[0].Description, receipt.Description)
}

func (s *ReceiptControllerSuite) TestFindReceipt_Success() {
	// prepare database
	receipt := model.Receipt{
		Description: "SALARY",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
	}
	s.db.Create(&receipt)

	// do request
	url := fmt.Sprintf("http://localhost:8080/budget-control/api/v1/receipt/%s", receipt.Id.String())
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	s.Require().NoError(err)

	var responseBody model.ReceiptResponse
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(responseBody.Date, receipt.Date)
	s.Require().Equal(responseBody.Description, receipt.Description)
	s.Require().Equal(responseBody.Value, receipt.Value)
}

func (s *ReceiptControllerSuite) TestUpdateReceipt_Success() {
	// prepare database
	receipt := model.Receipt{
		Description: "SALARY",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
	}
	s.db.Create(&receipt)

	// prepare request
	newReceipt := model.ReceiptRequest{
		Description: "New Description",
		Value:       1000,
		Date:        "2022-01-01T00:00:00Z",
	}

	body, err := json.Marshal(newReceipt)
	s.Require().NoError(err)
	requestBody := bytes.NewBuffer(body)
	url := fmt.Sprintf("http://localhost:8080/budget-control/api/v1/receipt/%s", receipt.Id.String())

	// do request
	request, err := http.NewRequest(http.MethodPut, url, requestBody)
	s.Require().NoError(err)
	request.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(request)
	s.Require().NoError(err)
	s.Require().Equal(resp.StatusCode, http.StatusNoContent)

	var got model.Receipt
	s.db.Find(&got, receipt.Id)

	s.Require().Equal(got.Date, newReceipt.Date)
	s.Require().Equal(got.Description, strings.ToUpper(newReceipt.Description))
	s.Require().Equal(got.Value, newReceipt.Value)
}

func (s *ReceiptControllerSuite) TestUpdateReceiptWithSameDescriptionInTheMonth_Fail() {
	// prepare database
	receipt := model.Receipt{
		Description: "SALARY",
		Value:       3000,
		Date:        "2022-01-20T00:00:00Z",
	}
	s.db.Create(&receipt)

	inMonthReceipt := model.Receipt{
		Description: "New Description",
		Value:       5000,
		Date:        "2022-01-15T00:00:00Z",
	}
	s.db.Create(&inMonthReceipt)
	// prepare request
	newReceipt := model.ReceiptRequest{
		Description: "NEW DESCRIPTION",
		Value:       1000,
		Date:        "2022-01-01T00:00:00Z",
	}

	body, err := json.Marshal(newReceipt)
	s.Require().NoError(err)
	requestBody := bytes.NewBuffer(body)
	url := fmt.Sprintf("http://localhost:8080/budget-control/api/v1/receipt/%s", receipt.Id.String())

	// do request
	request, err := http.NewRequest(http.MethodPut, url, requestBody)
	s.Require().NoError(err)
	request.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(request)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnprocessableEntity, resp.StatusCode)

	var responseBody struct {
		Error string
		Id    uuid.UUID
	}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(responseBody.Id, inMonthReceipt.Id)
	s.Require().Equal(responseBody.Error, fmt.Sprintf("receipt %s already created in current month", inMonthReceipt.Description))

	var got model.Receipt
	s.db.Find(&got, receipt.Id)

	s.Require().NotEqual(got.Date, newReceipt.Date)
	s.Require().NotEqual(got.Description, strings.ToUpper(newReceipt.Description))
	s.Require().NotEqual(got.Value, newReceipt.Value)
}

func (s *ReceiptControllerSuite) TestDeleteReceipt_Sucess() {
	// prepare database
	receipt := model.Receipt{
		Description: "SALARY",
		Value:       1000,
		Date:        "2022-01-20T00:00:00Z",
	}
	s.db.Create(&receipt)

	// prepare request
	url := fmt.Sprintf("http://localhost:8080/budget-control/api/v1/receipt/%s", receipt.Id.String())
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	s.Require().NoError(err)

	// do request
	resp, err := s.client.Do(request)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusNoContent, resp.StatusCode)

	tx := s.db.First(&model.Receipt{}, receipt.Id)
	s.Require().Equal(tx.Error, gorm.ErrRecordNotFound)
}

func (s *ReceiptControllerSuite) TestReceiptsByPeriod_Success() {
	// prepare database
	s.db.Create(&[]model.Receipt{
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

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/budget-control/api/v1/receipt/%s/%s", year, month))
	s.Require().NoError(err)

	var responseBody []model.Receipt
	json.NewDecoder(resp.Body).Decode(&responseBody)

	s.Require().Equal(2, len(responseBody))
}
