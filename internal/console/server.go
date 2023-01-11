package console

import (
	"context"
	"fmt"
	"github.com/irvankadhafi/erajaya-product-service/cache"
	"github.com/irvankadhafi/erajaya-product-service/internal/config"
	"github.com/irvankadhafi/erajaya-product-service/internal/db"
	"github.com/irvankadhafi/erajaya-product-service/internal/delivery/httpsvc"
	"github.com/irvankadhafi/erajaya-product-service/internal/helper"
	"github.com/irvankadhafi/erajaya-product-service/internal/repository"
	"github.com/irvankadhafi/erajaya-product-service/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var runServerCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Long:  `This subcommand start the server`,
	Run:   runServer,
}

func init() {
	RootCmd.AddCommand(runServerCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	// Initiate all connection like db, redis, etc
	db.InitializePostgresConn()
	pgDB, err := db.PostgreSQL.DB()
	continueOrFatal(err)
	defer helper.WrapCloser(pgDB.Close)

	cacheManager := cache.NewCache()

	redisConn, err := db.NewRedisConnPool(config.RedisHost())
	continueOrFatal(err)
	defer helper.WrapCloser(redisConn.Close)

	cacheManager.SetConnectionPool(redisConn)
	cacheManager.SetDefaultTTL(config.RedisCacheTTL())

	productRepository := repository.NewProductRepository(db.PostgreSQL, cacheManager)
	productUsecase := usecase.NewProductUsecase(productRepository)

	httpServer := echo.New()

	httpServer.Pre(middleware.AddTrailingSlash())
	httpServer.Use(middleware.Logger())
	httpServer.Use(middleware.Recover())
	httpServer.Use(middleware.CORS())

	apiGroup := httpServer.Group("/api")
	httpsvc.RouteService(apiGroup, productUsecase)

	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	quitCh := make(chan bool, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		for {
			select {
			case <-sigCh:
				gracefulShutdown(httpServer)
				quitCh <- true
			case e := <-errCh:
				logrus.Error(e)
				gracefulShutdown(httpServer)
				quitCh <- true
			}
		}
	}()

	go func() {
		// Start HTTP server
		if err := httpServer.Start(fmt.Sprintf(":%s", config.HTTPPort())); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	<-quitCh
	log.Info("exiting")
}

func continueOrFatal(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func gracefulShutdown(httpSvr *echo.Echo) {
	db.StopTickerCh <- true

	if httpSvr != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpSvr.Shutdown(ctx); err != nil {
			httpSvr.Logger.Fatal(err)
		}
	}
}
