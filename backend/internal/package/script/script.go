package script

import (
	"log"
	"os"

	"github.com/Kar-Su/uas-mobile.git/internal/database"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func Commands(injector do.Injector) bool {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)

	if len(os.Args) < 2 {
		return false
	}

	if os.Args[1] == "--seed" {
		if err := database.Seeder(db); err != nil {
			log.Fatalf("error migration seeder: %v", err)
		}
		log.Println("seeder completed successfully")
	}

	return false
}
