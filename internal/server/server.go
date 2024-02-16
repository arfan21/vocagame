package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arfan21/vocagame/config"
	"github.com/arfan21/vocagame/internal/middleware"
	"github.com/arfan21/vocagame/pkg/exception"
	"github.com/arfan21/vocagame/pkg/logger"
	"github.com/arfan21/vocagame/pkg/pkgutil"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	ctxTimeout = 5
)

type Server struct {
	app     *fiber.App
	db      *pgxpool.Pool
	dbRedis *redis.Client
}

func New(
	db *pgxpool.Pool,
	dbRedis *redis.Client,
) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: exception.FiberErrorHandler,
	})

	timeout := time.Duration(config.GetConfig().Service.Timeout) * time.Second

	app.Use(middleware.Timeout(timeout))
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger.Log(context.Background()),
		Fields: append(fiberzerolog.ConfigDefault.Fields, fiberzerolog.FieldRequestID),
	}))

	app.Use(cors.New())
	app.Use(requestid.New())

	app.Use(middleware.RequestIdUser())
	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	return &Server{
		app:     app,
		db:      db,
		dbRedis: dbRedis,
	}
}

func (s *Server) Run() error {
	s.Routes()
	ctx := context.Background()
	go func() {
		if err := s.app.Listen(pkgutil.GetPort()); err != nil {
			logger.Log(ctx).Fatal().Err(err).Msg("failed to start server")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	_, shutdown := context.WithTimeout(ctx, ctxTimeout*time.Second)
	defer shutdown()

	logger.Log(ctx).Info().Msg("shutting down server")
	return s.app.Shutdown()
}
