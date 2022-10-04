package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func putUsers(c *gin.Context) {
	users := make([]*UserRequest, 0)
	err := c.ShouldBindJSON(&users)
	if err != nil {
		c.Error(err)
		return
	}
	usersOut := make([]*UserResponse, 0)

	for _, u := range users {
		uo, err := u.ToResponse()
		if err != nil {
			c.Error(err)
			return
		}
		usersOut = append(usersOut, uo)
	}
	c.IndentedJSON(http.StatusOK, usersOut)
}

func RegisterHandlers(group *gin.RouterGroup) {
	users := group.Group("users")
	users.PUT("", putUsers)
}
