package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/syukur91/ischool-monitor/api/controller"
	"github.com/syukur91/ischool-monitor/service"

	"github.com/arifsetiawan/go-common/env"
	Middleware "github.com/syukur91/ischool-monitor/pkg/middleware"
	"gopkg.in/go-playground/validator.v9"
)

func init() {
	// Seed random number
	rand.Seed(time.Now().Unix())
}

type (
	// CustomValidator is
	CustomValidator struct {
		validator *validator.Validate
	}
)

// Validate is
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {

	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"disableCaller": true,
		"disableStacktrace": true,
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"initialFields": {
			"app_name": "` + env.Getenv("APP_NAME", "ischool-monitor") + `",
			"type":"` + env.Getenv("APP_TYPE", "api") + `",
			"version":"` + env.Getenv("VERSION", "0.1.0") + `"
		},
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "timeKey": "time",
		  "timeEncoder": "ISO8601",
		  "levelEncoder": "lowercase"
		}
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		log.Fatalf("Failed to initialize zap logger: %v\n", err)
	}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v\n", err)
	}
	defer logger.Sync()

	// @
	// initialize database connection
	if len(os.Getenv("DB_CONNECTION_STR")) == 0 {
		log.Fatalf("Database connection string is not set. Set DB_CONNECTION_STR in environment\n")
	}

	db, err := sqlx.Connect(env.Getenv("DB_DRIVER", "postgres"), os.Getenv("DB_CONNECTION_STR"))
	if err != nil {
		log.Fatalf("Failed to make database connection: %v\n", err)
	}
	defer db.Close()

	// @
	// Initialize echo
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HTTPErrorHandler = Middleware.ErrorHandler(logger)

	loggerConfig := Middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().RequestURI, "/public/") ||
				strings.Contains(c.Request().RequestURI, "favicon") ||
				strings.Contains(c.Request().RequestURI, "/js/") {
				return true
			}
			return false
		},
		AppName: "ischool-monitor",
		Logger:  logger,
	}

	e.Use(Middleware.LoggerWithConfig(loggerConfig))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// @
	// Create services
	mataPelajaranService := service.NewMata_PelajaranService(db)

	// @
	// Routes
	r := e.Group("/:tenant")

	// Mandatory hello world
	r.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello "+c.Param("tenant")+"! This is API version: "+os.Getenv("VERSION"))
	})

	// @
	// Handlers
	// Products API
	mataPelajaranHandler := &controller.Mata_PelajaranHandler{
		Mata_PelajaranService: mataPelajaranService,
	}
	mataPelajaranHandler.SetRoutes(r)

	// @
	// Start app
	servicePort := env.Getenv("PORT", ":6200")
	log.Println("ischool-monitor started at port " + servicePort)
	e.Start(servicePort)

}
