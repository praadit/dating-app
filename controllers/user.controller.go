package controllers

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

	req.Email = strings.ToLower(req.Email)

	res, err := c.service.Login(ctx.Request.Context(), req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *Controller) Signup(ctx *gin.Context) {
	req := &requests.SignupRequest{}
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

	if req.Password != req.ConfirmPassword {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Password and confirm password should be same",
		})
	}
	if len(req.Picture) > 1000 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Picture lenght should be less than 1000 characters",
		})
	}

	req.Email = strings.ToLower(req.Email)

	res, err := c.service.SignupUser(ctx.Request.Context(), req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
