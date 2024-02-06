package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

// userIdentity ...
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{"empty auth header"})
		return
	}

	if len(header) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{"token is empty"})
		return
	}

	//userId, err := h.service.Authorization.ParseToken(headerParts[1])
	//if err != nil {
	//	newErrorResponse(c, http.StatusUnauthorized, err.Error())
	//	return
	//}

	c.Set(userCtx, uint64(1))
}

// getUserId ...
func getUserId(c *gin.Context) (uint64, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(uint64)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
