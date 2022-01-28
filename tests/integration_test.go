package test_test

import (
	"encoding/json"
	"net/http"
	"os"
	"syscall"
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
	s.Require().NoError(godotenv.Load("test.env"))
	s.Require().NoError(database.Connect())
	s.db = database.DB
	s.m = database.M

	go router.HandleRequests()
}

func (s *ReceiptServiceTestSuite) TearDownTest() {
	s.db.Exec("DELETE FROM receipt")
}

func (s *ReceiptServiceTestSuite) TearDownSuite() {
	s.m.Down()
	s.db.Exec("DROP TABLE schema_migrations")
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGINT)
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
	resp, err := http.Get("http://localhost:8080/budget-control/api/v1/health")
	s.Require().NoError(err)
	defer resp.Body.Close()

	expect := struct{ Online bool }{true}
	var got struct{ Online bool }
	json.NewDecoder(resp.Body).Decode(&got)
	s.Require().Equal(expect, got)
}
