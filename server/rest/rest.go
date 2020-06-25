package rest

import (
	"fmt"
	"github.com/bahadrix/corpushub/repository"
	"github.com/bahadrix/corpushub/server/operator"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const VERSION = "0.1.0"

func setupRouter(host string, port int, op *operator.Operator) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")

	v1.GET("/ping", func(context *gin.Context) {
		context.String(200, "PONGv%s %s %d", VERSION, host, port)
	})

	v1.GET("/repos", func(context *gin.Context) {
		repos, err := op.GetRepos()

		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})

			return
		}

		context.JSON(200, repos)
	})

	v1.PUT("/repos", func(context *gin.Context) {
		var opts repository.RepoOptions

		if err := context.ShouldBindJSON(&opts); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := op.AddRepo(&opts)

		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		context.Status(200)
	})

	v1.GET("/repo/sync", func(context *gin.Context) {
		repoURI := context.Query("uri")

		err := op.SyncRepo(repoURI)

		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		context.Status(200)
	})

	v1.GET("/search", func(context *gin.Context) {
		q := context.Query("q")
		result, err := op.Search(q)
		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		context.JSON(200, result)
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
