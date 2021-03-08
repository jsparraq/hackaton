package controller

import (
	"encoding/json"
	"net/http"

	"github.com/jsparraq/api-rest/entity"
	"github.com/jsparraq/api-rest/errors"
	"github.com/jsparraq/api-rest/service"
)

type controller struct{}

var (
	postService service.PostService = service.NewPostService()
)

// PostController controller
type PostController interface {
	AddPost(response http.ResponseWriter, request *http.Request)
}

// NewPostController func
func NewPostController() PostController {
	return &controller{}
}

func (*controller) AddPost(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "application/json")
	var post entity.Post
	errDecode := json.NewDecoder(request.Body).Decode(&post)
	if errDecode != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Error marshalling the request"})
		return
	}
	errValidation := postService.Validate(&post)
	if errValidation != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: errValidation.Error()})
		return
	}
	result, errCreation := postService.Create(&post)
	if errCreation != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Error saving the post"})
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}
