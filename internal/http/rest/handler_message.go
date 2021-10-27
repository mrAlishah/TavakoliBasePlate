package rest

import (
	"GolangTraining/internal/subscription"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (h *Handler) SignIn(c *gin.Context) {
	var request subscription.Request

	if err := c.ShouldBind(&request); err != nil {
		if vErr, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, GetFailResponseFromValidationErrors(vErr))
		}
		return
	}

	msg, err := h.SubscriptionService.Test(request)
	if err != nil{
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, msg)
}

//func (h *Handler) GetMessages(c *gin.Context) {
//	msg, err := h.SubscriptionService.Test(c, "Omid1")
//	if err != nil{
//		c.JSON(http.StatusInternalServerError, err)
//		return
//	}
//	c.JSON(http.StatusOK, msg)
//}