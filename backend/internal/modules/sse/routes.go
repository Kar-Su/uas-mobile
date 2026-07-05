package sse

import (
	"github.com/Kar-Su/uas-mobile.git/internal/middlewares"
	authService "github.com/Kar-Su/uas-mobile.git/internal/modules/auth/service"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/sse/controller"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	sseController := do.MustInvoke[controller.SSEController](injector)
	jwtService := do.MustInvokeNamed[authService.JwtService](injector, constants.JWTService)
	auth := middlewares.AuthMiddleware(jwtService)
	allRoles := middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_GUDANG, constants.ROLE_USER)

	router.GET("/api/sse", auth, allRoles, sseController.Stream)
}
