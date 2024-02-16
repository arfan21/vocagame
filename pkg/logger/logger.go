package logger

import (
	"context"
	"os"
	"sync"

	"github.com/arfan21/vocagame/config"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"
)

var LoggerInstance zerolog.Logger
var once sync.Once

func Log(ctx context.Context) *zerolog.Logger {
	once.Do(func() {
		multi := zerolog.MultiLevelWriter(os.Stdout)
		LoggerInstance = zerolog.New(multi).With().Timestamp().Logger()

		if config.GetConfig().Env == "dev" {
			LoggerInstance = LoggerInstance.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

	})

	LoggerInstance = LoggerInstance.With().Ctx(ctx).Logger()

	if req_id, ok := ctx.Value(requestid.ConfigDefault.ContextKey).(string); ok {
		LoggerInstance = LoggerInstance.With().Str(fiberzerolog.FieldRequestID, req_id).Logger()
	}

	return &LoggerInstance
}
