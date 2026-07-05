package constants

import (
	"errors"
	"time"

	"github.com/Kar-Su/uas-mobile.git/internal/package/env"
)

const (
	DB      = "db"
	DB_TEST = "db_test"

	DB_LOG_DIR = "internal/configs/logs/query_log"

	JWTService      = "JWTService"
	JWT_ISSUER      = "Inventory App"
	JWT_ACCESS_EXP  = 15 * time.Minute
	JWT_REFRESH_EXP = 24 * time.Hour

	ROLE_SUPER_ADMIN  = "super-admin"
	ROLE_ADMIN_GUDANG = "admin-gudang"
	ROLE_USER         = "user"
)

var (
	JWT_SECRET_KEY = env.GetWithDefault[string]("JWT_SECRET", "")

	ErrInternalErr = errors.New("Internal Error")

	MESAGE_FAILED_GET_DATA_FROM_BODY = "failed to get data from body"
)
