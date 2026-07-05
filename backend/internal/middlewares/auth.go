package middlewares

import (
	"net/http"
	"strings"

	"github.com/Kar-Su/uas-mobile.git/internal/middlewares/dto"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/auth/service"
	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService service.JwtService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString != "" {
			if !strings.Contains(tokenString, "Bearer ") {
				res := utils.BuildResponseFailed(dto.FAILED_AUTH, dto.ErrInvalidHeader.Error())
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
				return
			}
			tokenString = tokenString[len("Bearer "):]
		} else {
			tokenString, _ = ctx.Cookie("access_token")
		}

		if tokenString == "" {
			res := utils.BuildResponseFailed(dto.FAILED_AUTH, dto.ErrHeaderMissing.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			res := utils.BuildResponseFailed(dto.FAILED_AUTH, dto.ErrInvalidToken.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		if !token.Valid {
			res := utils.BuildResponseFailed(dto.FAILED_AUTH, dto.ErrInvalidToken.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		userId, err := jwtService.GetUserIDByToken(tokenString)
		if err != nil {
			res := utils.BuildResponseFailed(dto.FAILED_AUTH, err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		roleName, err := jwtService.GetRoleNameByToken(tokenString)
		if err != nil {
			res := utils.BuildResponseFailed(dto.FAILED_AUTH, err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		userEmail, err := jwtService.GetUserEmailByToken(tokenString)
		if err != nil {
			res := utils.BuildResponseFailed(dto.FAILED_AUTH, err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		ctx.Set("user_id", userId)
		ctx.Set("role_name", roleName)
		ctx.Set("user_email", userEmail)
		ctx.Set("token", tokenString)
		ctx.Next()
	}
}
