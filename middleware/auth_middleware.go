package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	hAppconstant "github.com/michaelyusak/go-helper/appconstant"
	hDto "github.com/michaelyusak/go-helper/dto"
	hHelper "github.com/michaelyusak/go-helper/helper"
	"github.com/michaelyusak/xyz-kredit-plus/appconstant"
	"github.com/michaelyusak/xyz-kredit-plus/entity"
)

func AuthMiddleware(jwtHelper hHelper.JWTHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get(hAppconstant.Authorization)
		t := strings.Split(authHeader, " ")

		if len(t) != 2 || t[0] != hAppconstant.Bearer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, hDto.ErrorResponse{Message: hAppconstant.MsgUnauthorized})
			return
		}

		authToken := t[1]

		claimsBytes, err := jwtHelper.ParseAndVerify(authToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, hDto.ErrorResponse{Message: hAppconstant.MsgUnauthorized})
			return
		}

		var claims entity.JwtClaims

		err = json.Unmarshal(claimsBytes, &claims)
		if err != nil {
			// [TODO] Make sure request stop here
			// appErr := hApperror.InternalServerError(hApperror.AppErrorOpt{
			// 	Message: fmt.Sprintf("[auth_middleware][AuthMiddleware][json.Unmarshal] Error: %s", err.Error()),
			// })

			// c.Error(appErr)
			// return

			c.AbortWithStatusJSON(http.StatusInternalServerError, hDto.ErrorResponse{Message: hAppconstant.MsgInternalServerError})
			return
		}

		c.Set(appconstant.AccountIdCtxKey, claims.AccountId)
		c.Set(appconstant.EmailCtxKey, claims.Email)
		c.Set(appconstant.IsKycCompletedCtxKey, claims.AccountId)

		c.Next()
	}
}

func KycFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		isKycCompleted := c.Value(appconstant.IsKycCompletedCtxKey).(bool)

		if !isKycCompleted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, hDto.ErrorResponse{Message: hAppconstant.MsgUnauthorized})
			return
		}

		c.Next()
	}
}
