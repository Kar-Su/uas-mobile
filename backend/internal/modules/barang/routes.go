package barang

import (
	"github.com/Kar-Su/uas-mobile.git/internal/middlewares"
	authService "github.com/Kar-Su/uas-mobile.git/internal/modules/auth/service"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/barang/controller"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	barangController := do.MustInvoke[controller.BarangController](injector)
	jwtService := do.MustInvokeNamed[authService.JwtService](injector, constants.JWTService)
	auth := middlewares.AuthMiddleware(jwtService)
	adminGudang := middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_GUDANG)
	allRoles := middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_GUDANG, constants.ROLE_USER)

	routes := router.Group("/api/barang", auth)
	{
		routes.GET("", allRoles, barangController.GetAll)
		routes.GET("/:kode", allRoles, barangController.GetByKode)
		routes.POST("", adminGudang, barangController.Create)
		routes.PUT("/:kode", adminGudang, barangController.Update)
		routes.DELETE("/:kode", adminGudang, barangController.Delete)
	}
}
