/* основные настройки сайта */
package refersrest

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"refers_rest/pkg/settings"
	storage "refers_rest/pkg/storage/refersdb"
	"refers_rest/pkg/storage/refersdb/postgres"
)

var wait time.Duration

type Rest struct {

	Routes *gin.Engine

	Config *settings.Config

	DB storage.RefersDB

	ctx context.Context

	ExpTokenDay int

	SecretKey   string
}

func New(c *settings.Config) (*Rest, error) {

	db, err := postgres.New(c)
	if err != nil {
		return nil, err
	}

	rest := &Rest{
		Routes:      gin.Default(),
		Config:      c,
		ExpTokenDay: c.ExpTokenDay,
		SecretKey:   c.SecretKey,
		DB:          db,
	}

	rest.initializeRoutes()
	return rest, err
}

func (s *Rest) Start() error {

	connWs := net.JoinHostPort(s.Config.Host, s.Config.Port)
	log.Printf(`merch server start : %s`, connWs)

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	s.ctx = ctx

	srv := &http.Server{
		Addr:           connWs,
		WriteTimeout:   time.Second * 15,
		ReadTimeout:    time.Second * 15,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: 1 << 20,
		Handler:        s.Routes,
		// BaseContext: func(l net.Listener) context.Context {
		// 	return s.ctx
		// },
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	// https://ru.wikipedia.org/wiki/%D0%A1%D0%B8%D0%B3%D0%BD%D0%B0%D0%BB_(Unix)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)

	<-c

	s.DB.Close(ctx)
	srv.Shutdown(ctx)

	log.Println("shutting down")

	return nil
}
