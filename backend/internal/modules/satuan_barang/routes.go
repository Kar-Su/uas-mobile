package satuan_barang

import (
	"github.com/Kar-Su/uas-mobile.git/internal/middlewares"
	authService "github.com/Kar-Su/uas-mobile.git/internal/modules/auth/service"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/satuan_barang/controller"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func RegisterRoutes(router *gin.Engine, injector do.Injector) {
	satuanController := do.MustInvoke[controller.SatuanBarangController](injector)
	jwtService := do.MustInvokeNamed[authService.JwtService](injector, constants.JWTService)
	auth := middlewares.AuthMiddleware(jwtService)
	// superAdmin := middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN)
	adminGudang := middlewares.RoleMiddleware(constants.ROLE_SUPER_ADMIN, constants.ROLE_ADMIN_GUDANG)

	routes := router.Group("/api/satuan-barang", auth)
	{
		routes.GET("", adminGudang, satuanController.GetAll)
		routes.GET("/:id", adminGudang, satuanController.GetByID)
		routes.POST("", adminGudang, satuanController.Create)
		routes.PUT("/:id", adminGudang, satuanController.Update)
		routes.DELETE("/:id", adminGudang, satuanController.Delete)
	}
}
