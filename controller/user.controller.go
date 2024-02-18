package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/praadit/dating-apps/requests"
	"github.com/praadit/dating-apps/response"
	"github.com/praadit/dating-apps/utils"
)

func (c *Controller) Login(ctx *gin.Context) {
	req := &requests.LoginRequest{}
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

	req.Email = strings.ToLower(req.Email)

	res, err := c.service.Login(ctx.Request.Context(), req)
	if err != nil {
		errMessage := err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	ctx.JSON(http.StatusOK, response.FormatRequest(res, nil))
}

func (c *Controller) Signup(ctx *gin.Context) {
	req := &requests.SignupRequest{}
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

	if req.Password != req.ConfirmPassword {
		errMessage := "Password and confirm password should be same"
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
	}
	if len(req.Picture) > 1000 {
		errMessage := "Picture lenght should be less than 1000 characters"
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
	}

	req.Email = strings.ToLower(req.Email)

	err := c.service.SignupUser(ctx.Request.Context(), req)
	if err != nil {
		errMessage := err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FormatRequest(nil, &errMessage))
		return
	}

	ctx.JSON(http.StatusOK, response.FormatRequest(nil, nil))
}
