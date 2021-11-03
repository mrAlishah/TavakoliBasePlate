package rest

import (
	"GolangTraining/internal/subscription"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (h *Handler) CreateUser(c *gin.Context) {
	var request subscription.Request

	if err := c.ShouldBind(&request); err != nil {
		if vErr, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, GetFailResponseFromValidationErrors(vErr))
		}
		return
	}

	ctx := context.Background()
	user, err := h.SubscriptionService.Create(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
