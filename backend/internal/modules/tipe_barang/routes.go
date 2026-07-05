package tipe_barang

import (
	"github.com/Kar-Su/uas-mobile.git/internal/middlewares"
	authService "github.com/Kar-Su/uas-mobile.git/internal/modules/auth/service"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/tipe_barang/controller"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	tipeController := do.MustInvoke[controller.TipeBarangController](injector)
	jwtService := do.MustInvokeNamed[authService.JwtService](injector, constants.JWTService)
	auth := middlewares.AuthMiddleware(jwtService)
	// superAdmin := middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN)
	adminGudang := middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_GUDANG)

	routes := router.Group("/api/tipe-barang", auth)
	{
		routes.GET("", adminGudang, tipeController.GetAll)
		routes.GET("/:id", adminGudang, tipeController.GetByID)
		routes.POST("", adminGudang, tipeController.Create)
		routes.PUT("/:id", adminGudang, tipeController.Update)
		routes.DELETE("/:id", adminGudang, tipeController.Delete)
	}
}
