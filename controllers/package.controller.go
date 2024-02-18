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
		errMessage := utils.FilterError(err, "Failed to parse request body").Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	res, err := c.service.Packages(ctx.Request.Context(), req)
	if err != nil {
		errMessage := err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	ctx.JSON(http.StatusOK, response.FormatRequest(res, nil))
}
func (c *Controller) Package(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		errMessage := "Id cannot be null"
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errMessage := "Invalid id"
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	res, err := c.service.Package(ctx.Request.Context(), id)
	if err != nil {
		errMessage := err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	ctx.JSON(http.StatusOK, response.FormatRequest(res, nil))
}
func (c *Controller) Buy(ctx *gin.Context) {
	user, err := c.service.AuthUser(ctx.Request.Context())
	if err != nil {
		errMessage := err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	req := &requests.BuyPackage{}
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

	if err := c.service.Buy(ctx.Request.Context(), user, req); err != nil {
		errMessage := err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	ctx.JSON(http.StatusOK, response.FormatRequest(nil, nil))
}
