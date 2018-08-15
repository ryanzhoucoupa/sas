package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/dovadi/dbconfig"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// auto-generated build constants
var (
	version = "VERSION"
	commit  = "COMMIT"
)

var (
	ratePath   = "/rates"
	healthPath = "/health"
)

// mysql connections
var (
	connectionString = dbconfig.MysqlConnectionString("settings.json")
)

var (
	rengine *gin.Engine
)

var (
	db *sql.DB
)

var (
	// The flag package provides command line configuration options
	// You can see the options using the command line option --help which shows the descriptions below
	configurationFlag = flag.String("configuration-path", "conf.json", "Loads configuration file")
	maxStackTraceSize = 4096
	log               = logrus.New()
	isAppEngine       = true
)

// Configuration values for the JSON config file
type Configuration struct {
	BindAddress string `env:"SAS_BIND_ADDRESS"`
	Verbose     string `env:"SAS_VERBOSE"`
}

// Basic health check. Check to see if db connection is still there
func healthContext(c *gin.Context) {
	status := "ok"
	db, err := sql.Open("mysql", connectionString)

	_, err = db.Query("SELECT 1 from users limit 1")
	if err != nil {
		status = fmt.Sprintf("Error connecting with db with error: %v", err)
		log.Error(status)
	}

	c.JSON(http.StatusOK, gin.H{"status": status, "version": version, "commit": commit})
}

func buildRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET(healthPath, healthContext)
	}
}

func getHTTPClient(r *http.Request) *http.Client {
	return &http.Client{}
}

func loadConfiguration(configFile string) (*Configuration, error) {
	file, err := os.Open(configFile)
	if err != nil {
		log.Error(fmt.Sprintf("Error opening config file:%s error:%v", configFile, err))
		return nil, err
	}
	decoder := json.NewDecoder(file)
	var configuration Configuration
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Error(fmt.Sprintf("Error decoding config file:%s error:%v", configFile, err))
		return nil, err
	}
	configType := reflect.TypeOf(&configuration).Elem()
	configValue := reflect.ValueOf(&configuration).Elem()
	for i := 0; i < configType.NumField(); i++ {
		configField := configType.Field(i)
		envValue := os.Getenv(configField.Tag.Get("env"))
		if envValue != "" {
			configValue.FieldByName(configField.Name).SetString(envValue)
		}
	}
	return &configuration, nil
}

func ginErrorHandler(message string, err error, c *gin.Context, printStack bool, sendAirbrake bool) {
	w := gin.DefaultWriter
	w.Write([]byte(fmt.Sprintf("%s error:%v", message, err)))
	if printStack {
		trace := make([]byte, maxStackTraceSize)
		runtime.Stack(trace, false)
		w.Write([]byte(fmt.Sprintf("stack trace--\n%s\n--", trace)))
	}
	c.AbortWithError(http.StatusInternalServerError, err)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	rengine = gin.New()
	rengine.Use(gin.Logger())

	flag.Parse()
	log.Formatter = new(logrus.JSONFormatter)
	conf, err := loadConfiguration(*configurationFlag)
	if err != nil {
		log.Error(fmt.Sprintf("Error loading configuration: %v", err))
		return
	}

	buildRoutes(rengine)
	rengine.Run(conf.BindAddress)
}
