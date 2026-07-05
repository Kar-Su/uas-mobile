package auth

import (
	"github.com/Kar-Su/uas-mobile.git/internal/middlewares"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/auth/controller"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/auth/service"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	authController := do.MustInvoke[controller.AuthController](injector)
	jwtService := do.MustInvokeNamed[service.JwtService](injector, constants.JWTService)

	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/logout", middlewares.AuthMiddleware(jwtService), authController.Logout)
		authRoutes.POST("/refresh-token", authController.RefreshToken)
		authRoutes.GET("/refresh-token/:refresh_token", middlewares.AuthMiddleware(jwtService), authController.FindRefreshToken)
		authRoutes.POST("/reset-password", middlewares.AuthMiddleware(jwtService), authController.ResetPassword)
	}
}
