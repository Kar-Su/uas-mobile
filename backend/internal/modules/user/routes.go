package user

import (
	"github.com/Kar-Su/uas-mobile.git/internal/middlewares"
	authService "github.com/Kar-Su/uas-mobile.git/internal/modules/auth/service"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/user/controller"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	userController := do.MustInvoke[controller.UserController](injector)
	jwtService := do.MustInvokeNamed[authService.JwtService](injector, constants.JWTService)
	auth := middlewares.AuthMiddleware(jwtService)
	superAdmin := middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN)

	router.GET("/api/me", auth, userController.Me)

	userRoutes := router.Group("/api/users", auth, superAdmin)
	// Testing non auth
	// userRoutes := router.Group("/api/users")
	{
		userRoutes.GET("", userController.GetUsers)
		userRoutes.GET("/:id", userController.GetUser)
		userRoutes.POST("", userController.CreateUser)
		userRoutes.PUT("/:id", userController.UpdateUser)
		userRoutes.DELETE("/:id", userController.DeleteUser)
	}
}
