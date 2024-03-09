package myhttp

import (
	"github.com/gin-gonic/gin"
	"infrastructure/myerror"
	"net/http"
)

type responseJSONSuccess struct {
	Data interface{} `json:"data"`
}

type responseJSONError struct {
	Error string `json:"error"`
}

func EncodeJSONSuccess(c *gin.Context, response interface{}) {
	res := &responseJSONSuccess{
		Data: response,
	}

	c.JSON(http.StatusOK, res)
}

func EncodeJSONError(c *gin.Context, err error) {
	parsedErr := myerror.GetParsedError(err)

	res := &responseJSONError{
		Error: parsedErr.Message,
	}

	c.JSON(getHTTPCode(parsedErr), res)
}

func getHTTPCode(err error) int {
	switch myerror.GetParsedError(err).Type {
	case myerror.BadRequestError:
		return http.StatusBadRequest
	case myerror.NotFoundError:
		return http.StatusNotFound
	case myerror.ForbiddenError:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
