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
	"github.com/RaphaSalomao/alura-challenge-backend/test/factory"
	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ControllerSuite struct {
	suite.Suite
	db   *gorm.DB
	m    *migrate.Migrate
	port string
}

func (s *ControllerSuite) SetupSuite() {
	s.Require().NoError(godotenv.Load("../../test.env"))
	s.Require().NoError(database.Connect())
	s.db = database.DB
	s.m = database.M
	s.port = os.Getenv("SRV_PORT")
	go router.HandleRequests()
	time.Sleep(2 * time.Second)
}

func (s *ControllerSuite) TearDownTest() {
	s.db.Exec("DELETE FROM receipts")
	s.db.Exec("DELETE FROM expenses")
	s.db.Exec("DELETE FROM users")
}

func (s *ControllerSuite) TearDownSuite() {
	s.m.Down()
	s.db.Exec("DROP TABLE schema_migrations")
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGINT)
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerSuite))
}
func (s *ControllerSuite) TestHealthCheck_Success() {
	resp, err := http.Get("http://localhost:5000/budget-control/api/v1/health")
	s.Require().NoError(err)
	defer resp.Body.Close()

	expect := struct{ Online bool }{true}
	var got struct{ Online bool }
	json.NewDecoder(resp.Body).Decode(&got)

	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal(expect, got)
}

func (s *ControllerSuite) TestMonthBalanceSumary_Success() {
	// prepare database
	r := factory.Request{
		User: model.User{
			Email:    "email@email.com",
			Password: "password",
		},
		Method: http.MethodGet,
		DB:     s.db,
		Client: http.Client{},
		Port:   s.port,
	}
	r.SaveUser()

	receipts := []model.Receipt{
		{
			Description: "Receipt 1",
			Value:       1100,
			Date:        "2020-01-01T00:00:00Z",
			UserId:      r.User.Id,
		},
		{
			Description: "Receipt 2",
			Value:       1200,
			Date:        "2020-01-02T00:00:00Z",
			UserId:      r.User.Id,
		},
		{
			Description: "Receipt 3",
			Value:       1300,
			Date:        "2020-01-03T00:00:00Z",
			UserId:      r.User.Id,
		},
	}
	expenses := []model.Expense{
		{
			Description: "Expense 1",
			Value:       1100,
			Date:        "2020-01-01T00:00:00Z",
			Category:    enum.CategoryFood,
			UserId:      r.User.Id,
		},
		{
			Description: "Expense 2",
			Value:       250,
			Date:        "2020-01-02T00:00:00Z",
			Category:    enum.CategoryHealth,
			UserId:      r.User.Id,
		},
		{
			Description: "Expense 3",
			Value:       100,
			Date:        "2020-01-03T00:00:00Z",
			Category:    enum.CategoryHealth,
			UserId:      r.User.Id,
		},
		{
			Description: "Expense 4",
			Value:       1000,
			Date:        "2020-01-04T00:00:00Z",
			Category:    enum.CategoryFood,
			UserId:      r.User.Id,
		},
	}
	s.db.Create(&receipts)
	s.db.Create(&expenses)

	// prepare expected response
	totalReceipt := 0.0
	totalExpense := 0.0
	categoryBalance := map[enum.Category]float64{}

	for _, receipt := range receipts {
		totalReceipt += receipt.Value
	}
	for _, expense := range expenses {
		totalExpense += expense.Value
		categoryBalance[expense.Category] += expense.Value
	}
	monthBalance := totalReceipt - totalExpense

	// do request
	year, month := "2020", "01"
	r.Path = fmt.Sprintf("/budget-control/api/v1/summary/%s/%s", year, month)
	resp, err := r.DoRequest()
	s.Require().NoError(err)

	// assert response
	var bs model.BalanceSumaryResponse
	json.NewDecoder(resp.Body).Decode(&bs)

	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal(monthBalance, bs.MonthBalance)
	s.Require().Equal(totalReceipt, bs.TotalReceipt)
	s.Require().Equal(totalExpense, bs.TotalExpense)
	s.Require().Equal(categoryBalance, bs.CategoryBalance)
}

func (s *ControllerSuite) TestCreateUser_Success() {
	// prepare request
	expect := model.UserRequest{
		Email:    "email@email.com",
		Password: "password",
	}

	request, err := json.Marshal(expect)
	s.Require().NoError(err)
	requestBody := bytes.NewBuffer(request)

	// do request
	resp, err := http.Post("http://localhost:5000/budget-control/api/v1/user", "application/json", requestBody)

	// assert response
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var user model.User
	s.db.Where("email = ?", expect.Email).First(&user)
	s.Require().Equal(expect.Email, user.Email)
	s.Require().Equal(true, utils.ValidadeHashAndPassword(strings.ToLower(expect.Password), user.Password))
}

func (s *ControllerSuite) TestAuthenticate_Success() {
	// prepare database
	password, err := utils.HashPassword("password")
	s.Require().NoError(err)
	user := model.User{
		Email:    "email@email.com",
		Password: password,
	}
	s.db.Create(&user)

	// prepare request
	userRequest := model.UserRequest{
		Email:    "email@email.com",
		Password: "password",
	}

	request, err := json.Marshal(userRequest)
	s.Require().NoError(err)
	requestBody := bytes.NewBuffer(request)

	// do request
	resp, err := http.Post("http://localhost:5000/budget-control/api/v1/authenticate", "application/json", requestBody)
	s.Require().NoError(err)

	var tokenResponse struct{ Token string }
	json.NewDecoder(resp.Body).Decode(&tokenResponse)
	userId, err := utils.ParseToken(tokenResponse.Token)

	// assert response
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().NoError(err)
	s.Require().Equal(user.Id.String(), userId)
}

func (s *ControllerSuite) TestAuthenticate_Fail() {
	// prepare request
	userRequest := model.UserRequest{
		Email:    "email@email.com",
		Password: "password",
	}

	request, err := json.Marshal(userRequest)
	s.Require().NoError(err)
	requestBody := bytes.NewBuffer(request)

	// do request
	resp, err := http.Post("http://localhost:5000/budget-control/api/v1/authenticate", "application/json", requestBody)

	var tokenResponse struct{ Token string }
	json.NewDecoder(resp.Body).Decode(&tokenResponse)

	// assert response
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)
}
