package providers

import (
	"github.com/Kar-Su/uas-mobile.git/internal/configs"
	authController "github.com/Kar-Su/uas-mobile.git/internal/modules/auth/controller"
	authRepo "github.com/Kar-Su/uas-mobile.git/internal/modules/auth/repository"
	authService "github.com/Kar-Su/uas-mobile.git/internal/modules/auth/service"
	barangController "github.com/Kar-Su/uas-mobile.git/internal/modules/barang/controller"
	barangRepo "github.com/Kar-Su/uas-mobile.git/internal/modules/barang/repository"
	barangService "github.com/Kar-Su/uas-mobile.git/internal/modules/barang/service"
	satuanController "github.com/Kar-Su/uas-mobile.git/internal/modules/satuan_barang/controller"
	satuanRepo "github.com/Kar-Su/uas-mobile.git/internal/modules/satuan_barang/repository"
	satuanService "github.com/Kar-Su/uas-mobile.git/internal/modules/satuan_barang/service"
	sseController "github.com/Kar-Su/uas-mobile.git/internal/modules/sse/controller"
	tipeController "github.com/Kar-Su/uas-mobile.git/internal/modules/tipe_barang/controller"
	tipeRepo "github.com/Kar-Su/uas-mobile.git/internal/modules/tipe_barang/repository"
	tipeService "github.com/Kar-Su/uas-mobile.git/internal/modules/tipe_barang/service"
	userController "github.com/Kar-Su/uas-mobile.git/internal/modules/user/controller"
	userRepo "github.com/Kar-Su/uas-mobile.git/internal/modules/user/repository"
	userService "github.com/Kar-Su/uas-mobile.git/internal/modules/user/service"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func InitDatabases(injector do.Injector) {
	do.ProvideNamed[*gorm.DB](injector, constants.DB, func(i do.Injector) (*gorm.DB, error) {
		return configs.SetUpDatabaseConnection(), nil
	})
}

func RegisterProviders(injector do.Injector) {
	do.ProvideNamed[authService.JwtService](injector, constants.JWTService, func(i do.Injector) (authService.JwtService, error) {
		return authService.NewJwtService(), nil
	})

	InitDatabases(injector)
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)

	jwtService := do.MustInvokeNamed[authService.JwtService](injector, constants.JWTService)

	roleRepoInst := userRepo.NewRoleRepository(db)
	userRepoInst := userRepo.NewUserRepository(db, roleRepoInst)
	refreshTokenRepoInst := authRepo.NewRefreshTokenRepository(db)
	tipeRepoInst := tipeRepo.NewTipeBarangRepository(db)
	satuanRepoInst := satuanRepo.NewSatuanBarangRepository(db)
	barangRepoInst := barangRepo.NewBarangRepository(db)

	userServiceInst := userService.NewUserService(userRepoInst, roleRepoInst, db)
	tipeServiceInst := tipeService.NewTipeBarangService(tipeRepoInst, db)
	satuanServiceInst := satuanService.NewSatuanBarangService(satuanRepoInst, db)
	barangServiceInst := barangService.NewBarangService(barangRepoInst, tipeRepoInst, satuanRepoInst, db)

	do.Provide(injector, func(i do.Injector) (authService.AuthService, error) {
		return authService.NewAuthService(userRepoInst, refreshTokenRepoInst, jwtService, db), nil
	})

	do.Provide(injector, func(i do.Injector) (userController.UserController, error) {
		return userController.NewUserController(i, db, userServiceInst, roleRepoInst), nil
	})

	do.Provide(injector, func(i do.Injector) (authController.AuthController, error) {
		authSvc := do.MustInvoke[authService.AuthService](i)
		return authController.NewAuthController(i, authSvc, db), nil
	})

	do.Provide(injector, func(i do.Injector) (tipeController.TipeBarangController, error) {
		return tipeController.NewTipeBarangController(i, db, tipeServiceInst), nil
	})

	do.Provide(injector, func(i do.Injector) (satuanController.SatuanBarangController, error) {
		return satuanController.NewSatuanBarangController(i, db, satuanServiceInst), nil
	})

	do.Provide(injector, func(i do.Injector) (barangController.BarangController, error) {
		return barangController.NewBarangController(i, db, barangServiceInst), nil
	})

	do.Provide(injector, func(i do.Injector) (sseController.SSEController, error) {
		return sseController.NewSSEController(), nil
	})
}
