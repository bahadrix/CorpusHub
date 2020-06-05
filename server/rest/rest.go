package rest

import (
	"fmt"
	"github.com/bahadrix/corpushub/server/operator"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const VERSION = "0.1.0"


func setupRouter(host string, port int, op *operator.Operator) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")

	v1.GET("/ping", func(context *gin.Context) {
		context.String(200, "PONGv%s %s %d", VERSION, host, port)
	})

	return router
}

func Start(host string, port int, debugMode bool, op *operator.Operator) error {

	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		log.Info("Starting in release mode at ", host, ":", port)
		gin.SetMode(gin.ReleaseMode)
	}

	router := setupRouter(host, port, op)
	return router.Run(fmt.Sprintf("%s:%d", host, port))

}
