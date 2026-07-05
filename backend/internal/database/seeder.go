package database

import (
	"context"

	"github.com/Kar-Su/uas-mobile.git/internal/database/seeders/seeds"

	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	ctx := context.Background()

	if err := seeds.SeedRoles(ctx, db); err != nil {
		return err
	}

	if err := seeds.SeedSuperAdmin(ctx, db); err != nil {
		return err
	}

	if err := seeds.SeedTipeBarang(ctx, db); err != nil {
		return err
	}

	if err := seeds.SeedSatuanBarang(ctx, db); err != nil {
		return err
	}

	if err := seeds.SeedBarang(ctx, db); err != nil {
		return err
	}

	return nil
}
