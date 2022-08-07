package app

import (
	config "aspro/cmd/configs"
	"aspro/internal/domain/handler"
	"aspro/internal/domain/repository"
	"aspro/internal/domain/service"
	"aspro/pkg/client/postgresql"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	router     *httprouter.Router
	httpServer *http.Server
	handlers   *handler.Handler
}

func NewApp() (App, error) {
	logrus.Println("router initializing")

	user := viper.GetString("db.username")
	pass := viper.GetString("db.password")
	router := httprouter.New()
	router.GET("/", handler.Index)
	router.GET("/protected/", handler.BasicAuth(handler.Protected, user, pass))

	logrus.Println("swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logrus.Println("heartbeat metric initializing")
	metricHandler := handler.Handler{}
	metricHandler.Register(router)

	pgConfig, err := postgresql.NewConfig(&postgresql.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Database: viper.GetString("db.database"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	pgClient := repository.NewRepository(pgConfig)
	services := service.NewService(pgClient)
	handlers := handler.NewHandler(services)

	return App{
		router:   router,
		handlers: handlers,
	}, nil
}

func (a *App) Run() {
	a.startHTTP()
}

func (a *App) startHTTP() {
	logrus.Info("start HTTP")

	var listener net.Listener

	if viper.GetString("listen.type") == config.LISTEN_TYPE_SOCK {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logrus.Fatal(err)
		}
		socketPath := path.Join(appDir, viper.GetString("listen.socket_file"))
		logrus.Infof("socket path: %s", socketPath)

		logrus.Info("created and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		logrus.Infof("bind application to host: %s and port: %s", viper.GetString("listen.bind_ip"), viper.GetString("listen.port"))
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", viper.GetString("listen.bind_ip"), viper.GetString("listen.port")))
		if err != nil {
			logrus.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Location", "Authorization", "Content-Disposition"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:        handler,
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logrus.Println("application completely initialized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logrus.Warn("server shutdown")
		default:
			logrus.Fatal(err)
		}
	}
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
