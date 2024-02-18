package controller

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
		errMessage := err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	selUser, err := c.service.Explore(ctx.Request.Context(), user)
	if err != nil {
		errMessage := err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	ctx.JSON(http.StatusOK, response.FormatRequest(selUser.ToResponse(), nil))
}

func (c *Controller) Swipe(ctx *gin.Context) {
	user, err := c.service.AuthUser(ctx.Request.Context())
	if err != nil {
		errMessage := err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	req := &requests.SwipeRequest{}
	if err := ctx.Bind(req); err != nil {
		errMessage := utils.FilterError(err, "Failed to parse request body").Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	if err := utils.ValidateStruct(c.validate, req); err != nil {
		errMessage := err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	if err := c.service.Swipe(ctx.Request.Context(), user, req); err != nil {
		errMessage := err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	ctx.JSON(http.StatusOK, response.FormatRequest(nil, nil))
}
