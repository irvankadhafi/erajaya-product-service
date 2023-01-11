package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/golang/mock/gomock"
	"github.com/irvankadhafi/erajaya-product-service/internal/config"
	"github.com/irvankadhafi/erajaya-product-service/internal/model/mock"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
	"testing"
)

func initializeTest() {
	config.GetConf()
	setupLogger()
}

func setupLogger() {
	formatter := runtime.Formatter{
		ChildFormatter: &logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		},
		Line: true,
		File: true,
	}

	logrus.SetFormatter(&formatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)

	verbose, _ := strconv.ParseBool(os.Getenv("VERBOSE"))
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

type repoTestKit struct {
	dbmock          sqlmock.Sqlmock
	db              *gorm.DB
	ctrl            *gomock.Controller
	mockProductRepo *mock.MockProductRepository
}

func initializeRepoTestKit(t *testing.T) (kit *repoTestKit, close func()) {
	dbConn, dbMock, err := sqlmock.New()
	if err != nil {
		logrus.Fatal(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: dbConn}), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	productRepo := mock.NewMockProductRepository(ctrl)
	tk := &repoTestKit{
		ctrl:            ctrl,
		dbmock:          dbMock,
		db:              gormDB,
		mockProductRepo: productRepo,
	}

	return tk, func() {
		if conn, _ := tk.db.DB(); conn != nil {
			_ = conn.Close()
		}
	}
}
