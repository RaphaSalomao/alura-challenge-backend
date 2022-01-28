package test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/RaphaSalomao/alura-challenge-backend/database"
	"github.com/RaphaSalomao/alura-challenge-backend/model"
	"github.com/RaphaSalomao/alura-challenge-backend/router"
	"github.com/RaphaSalomao/alura-challenge-backend/service"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ReceiptServiceTestSuite struct {
	suite.Suite
	db *gorm.DB
	m  *migrate.Migrate
}

func (s *ReceiptServiceTestSuite) SetupSuite() {
	go router.HandleRequests()
	s.Require().NoError(godotenv.Load("test.env"))
	s.Require().NoError(database.Connect())
	s.db = database.DB
	s.m = database.M
}

func (s *ReceiptServiceTestSuite) TearDownTest() {
	s.db.Exec("DELETE FROM receipt")
}

func (s *ReceiptServiceTestSuite) TearDownSuite() {
	s.m.Down()
	s.db.Exec("DROP TABLE schema_migrations")
}

func TestReceiptServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReceiptServiceTestSuite))
}

func (s *ReceiptServiceTestSuite) TestFindReceipts_Fail() {
	receiptService := service.RreceiptService{
		DB: *s.db,
	}
	var receiptResponse model.ReceiptResponse
	id := uuid.New()
	response := receiptService.FindReceipt(&receiptResponse, id)
	s.Require().Error(response)
	s.Require().Equal(response.Error(), "receipt not found")
}

func (s *ReceiptServiceTestSuite) TestHealthCheck_Success() {
	req, err := http.NewRequest(http.MethodGet, ("http://localhost:8080/budget-control/api/v1/health"), strings.NewReader(""))
	if err != nil {
		fmt.Println("ERROR AT REQUEST")
	}
	fmt.Println(req.Response.StatusCode)
}
