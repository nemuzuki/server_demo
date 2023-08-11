package routes

import(
	"github.com/gin-gonic/gin"
	"server_demo/controller"
)
func ApiRoutes(r *gin.Engine) {

	apiRouter := r.Group("server_demo")

	apiRouter.GET("/getTopic/", controller.QueryTopicById)
	apiRouter.GET("/getPost/", controller.QueryPostById)
}