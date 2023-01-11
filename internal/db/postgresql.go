package db

import (
	"github.com/irvankadhafi/erajaya-product-service/internal/config"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

var (
	// PostgreSQL represents gorm DB
	PostgreSQL *gorm.DB

	// StopTickerCh signal for closing ticker channel
	StopTickerCh chan bool
)

// InitializePostgresConn :nodoc:
func InitializePostgresConn() {
	conn, err := openPostgresConn(config.DatabaseDSN())
	if err != nil {
		log.WithField("databaseDSN", config.DatabaseDSN()).Fatal("failed to connect postgresql database: ", err)
	}

	PostgreSQL = conn
	StopTickerCh = make(chan bool)

	go checkConnection(time.NewTicker(config.DatabasePingInterval()))

	switch config.LogLevel() {
	case "error":
		PostgreSQL.Logger = PostgreSQL.Logger.LogMode(gormLogger.Error)
	case "warn":
		PostgreSQL.Logger = PostgreSQL.Logger.LogMode(gormLogger.Warn)
	default:
		PostgreSQL.Logger = PostgreSQL.Logger.LogMode(gormLogger.Info)

	}

	log.Info("Connection to PostgreSQL Server success...")
}

func checkConnection(ticker *time.Ticker) {
	for {
		select {
		case <-StopTickerCh:
			ticker.Stop()
			return
		case <-ticker.C:
			if _, err := PostgreSQL.DB(); err != nil {
				reconnectPostgresConn()
			}
		}
	}
}

func reconnectPostgresConn() {
	b := backoff.Backoff{
		Factor: 2,
		Jitter: true,
		Min:    100 * time.Millisecond,
		Max:    1 * time.Second,
	}

	postgresRetryAttempts := config.DatabaseRetryAttempts()

	for b.Attempt() < postgresRetryAttempts {
		conn, err := openPostgresConn(config.DatabaseDSN())
		if err != nil {
			log.WithField("databaseDSN", config.DatabaseDSN()).Error("failed to connect postgresql database: ", err)
		}

		if conn != nil {
			PostgreSQL = conn
			break
		}
		time.Sleep(b.Duration())
	}

	if b.Attempt() >= postgresRetryAttempts {
		log.Fatal("maximum retry to connect database")
	}
	b.Reset()
}

func openPostgresConn(dsn string) (*gorm.DB, error) {
	psqlDialector := postgres.Open(dsn)
	db, err := gorm.Open(psqlDialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	conn, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMaxIdleConns(config.DatabaseMaxIdleConns())
	conn.SetMaxOpenConns(config.DatabaseMaxOpenConns())
	conn.SetConnMaxLifetime(config.DatabaseConnMaxLifetime())

	return db, nil
}
