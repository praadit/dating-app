package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/praadit/dating-apps/requests"
	"github.com/praadit/dating-apps/response"
	"github.com/praadit/dating-apps/utils"
)

func (c *Controller) Packages(ctx *gin.Context) {
	req := &requests.Pagination{
		Page:    1,
		PerPage: 10,
		Order:   "asc",
	}
	if err := ctx.BindQuery(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: utils.FilterError(err, "Failed to parse request body").Error(),
		})
		return
	}

	res, err := c.service.Packages(ctx.Request.Context(), req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
func (c *Controller) Package(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Id cannot be null",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid id",
		})
		return
	}

	res, err := c.service.Package(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
func (c *Controller) Buy(ctx *gin.Context) {
	user, err := c.service.AuthUser(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	req := &requests.BuyPackage{}
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

	if err := c.service.Buy(ctx.Request.Context(), user, req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status: true,
	})
}
