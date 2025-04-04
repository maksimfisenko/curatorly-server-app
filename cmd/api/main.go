package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"os"
	"runtime"
	"time"

	_ "github.com/lib/pq"
	"github.com/maksimfisenko/curatorly-server-app/internal/data"
	"github.com/maksimfisenko/curatorly-server-app/internal/jsonlog"

	_ "github.com/maksimfisenko/curatorly-server-app/docs"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	jwt struct {
		secret string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

//	@title			Curatorly Server App API
//	@version		1.0
//	@description	This is API for the Curaturly server app.

// @host		localhost:8000
// @BasePath	/
func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "api server port")
	flag.StringVar(&cfg.env, "env", "development", "environment (development |  staging | production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "postgresql dsn")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "postgresql max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "postgresql max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "postgresql max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "enable rate limiter")

	flag.StringVar(&cfg.jwt.secret, "jwt-secret", "", "jwt secret")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
