package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praadit/dating-apps/requests"
	"github.com/praadit/dating-apps/response"
	"github.com/praadit/dating-apps/utils"
)

func (c *Controller) Explore(ctx *gin.Context) {
	user, err := c.service.AuthUser(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	selUser, err := c.service.Explore(ctx.Request.Context(), user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, selUser.ToResponse())
}

func (c *Controller) Swipe(ctx *gin.Context) {
	user, err := c.service.AuthUser(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	req := &requests.SwipeRequest{}
	if err := ctx.Bind(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: utils.FilterError(err, "Failed to parse request body").Error(),
		})
		return
	}

	if err := utils.ValidateStruct(c.validate, req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := c.service.Swipe(ctx.Request.Context(), user, req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status: true,
	})
}
