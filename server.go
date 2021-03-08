package main

import (
	"fmt"
	"net/http"

	"github.com/jsparraq/api-rest/controller"
	router "github.com/jsparraq/api-rest/http"
)

var (
	httpRouter     router.Router             = router.NewMuxRouter()
	postController controller.PostController = controller.NewPostController()
)

func main() {
	const port string = ":8000"
	httpRouter.GET("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(res, "Up and running ...")
	})

	httpRouter.POST("/posts", postController.AddPost)
	httpRouter.SERVE(port)
}
