package seeds

import (
	"context"
	"time"

	"github.com/Kar-Su/uas-mobile.git/internal/database/entities"
	"github.com/Kar-Su/uas-mobile.git/internal/package/helpers"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedTipeBarang(ctx context.Context, db *gorm.DB) error {
	tipe := []entities.TipeBarang{
		{Name: "Elektronik"},
		{Name: "Pakaian"},
		{Name: "Makanan & Minuman"},
		{Name: "Alat Tulis Kantor"},
		{Name: "Perlengkapan Rumah Tangga"},
	}
	return db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&tipe).Error
}

func SeedSatuanBarang(ctx context.Context, db *gorm.DB) error {
	satuan := []entities.SatuanBarang{
		{Satuan: strPtr("pcs"), Keterangan: strPtr("pieces")},
		{Satuan: strPtr("kg"), Keterangan: strPtr("kilogram")},
		{Satuan: strPtr("liter"), Keterangan: strPtr("liter")},
		{Satuan: strPtr("meter"), Keterangan: strPtr("meter")},
		{Satuan: strPtr("box"), Keterangan: strPtr("box / dus")},
		{Satuan: strPtr("pack"), Keterangan: strPtr("pack / bungkus")},
		{Satuan: strPtr("unit"), Keterangan: strPtr("unit")},
	}
	return db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&satuan).Error
}

func SeedBarang(ctx context.Context, db *gorm.DB) error {
	now := time.Now()

	type barangSeed struct {
		kode     string
		name     string
		tipeID   int
		satuanID int
		qty      int
	}

	items := []barangSeed{
		// Elektronik (tipe_id=1)
		{"BRG-001", "Kabel HDMI 2 meter", 1, 1, 50},
		{"BRG-002", "Mouse Wireless Logitech", 1, 1, 30},
		{"BRG-003", "Adaptor USB-C 65W", 1, 1, 20},
		// Pakaian (tipe_id=2)
		{"BRG-004", "Kaos Polos Hitam", 2, 1, 100},
		{"BRG-005", "Kemeja Flanel", 2, 1, 40},
		{"BRG-006", "Jaket Hoodie", 2, 1, 25},
		// Makanan & Minuman (tipe_id=3)
		{"BRG-007", "Kopi Arabica 250gr", 3, 6, 60},
		{"BRG-008", "Minyak Goreng 1L", 3, 3, 80},
		{"BRG-009", "Gula Pasir 1kg", 3, 2, 90},
		// Alat Tulis Kantor (tipe_id=4)
		{"BRG-010", "Buku Tulis Sidu A5", 4, 1, 200},
		{"BRG-011", "Pulpen Standard AE7", 4, 5, 150},
		{"BRG-012", "Stapler Small", 4, 1, 45},
		// Perlengkapan Rumah Tangga (tipe_id=5)
		{"BRG-013", "Sapu Lantai", 5, 1, 35},
		{"BRG-014", "Ember Plastik 10L", 5, 1, 40},
		{"BRG-015", "Lap Microfiber", 5, 1, 70},
	}

	barang := make([]entities.Barang, len(items))
	for i, it := range items {
		barang[i] = entities.Barang{
			Kode:      it.kode,
			Name:      it.name,
			TipeID:    it.tipeID,
			SatuanID:  it.satuanID,
			Quantity:  it.qty,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&barang).Error
}

func strPtr(s string) *string { return &s }

func SeedRoles(ctx context.Context, db *gorm.DB) error {
	roles := []entities.Role{
		{ID: 1, Name: "super-admin"},
		{ID: 2, Name: "admin-gudang"},
		{ID: 3, Name: "user"},
	}
	return db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error
}

func SeedSuperAdmin(ctx context.Context, db *gorm.DB) error {
	password, err := helpers.HashPassword("superadmin")
	if err != nil {
		return err
	}

	users := make([]entities.User, 0, 3)
	users = append(users, entities.User{
		Email:    "super@akun.com",
		Name:     "Super Admin",
		Password: password,
		RoleID:   1,
	})

	password, err = helpers.HashPassword("gudangadmin")
	if err != nil {
		return err
	}

	users = append(users, entities.User{
		Email:    "admin@akun.com",
		Name:     "admin gudang",
		Password: password,
		RoleID:   2,
	})

	password, err = helpers.HashPassword("userbiasa")
	if err != nil {
		return err
	}

	users = append(users, entities.User{
		Email:    "user@akun.com",
		Name:     "Helmi",
		Password: password,
		RoleID:   3,
	})

	password, err = helpers.HashPassword("userbiasa")
	if err != nil {
		return err
	}

	users = append(users, entities.User{
		Email:    "dummy1@akun.com",
		Name:     "dummy1",
		Password: password,
		RoleID:   3,
	})
	password, err = helpers.HashPassword("userbiasa")
	if err != nil {
		return err
	}

	users = append(users, entities.User{
		Email:    "dummy2@akun.com",
		Name:     "dummy2",
		Password: password,
		RoleID:   3,
	})
	password, err = helpers.HashPassword("userbiasa")
	if err != nil {
		return err
	}

	users = append(users, entities.User{
		Email:    "dummy3@akun.com",
		Name:     "dummy3",
		Password: password,
		RoleID:   3,
	})
	password, err = helpers.HashPassword("userbiasa")
	if err != nil {
		return err
	}

	users = append(users, entities.User{
		Email:    "dummy4@akun.com",
		Name:     "dummy4",
		Password: password,
		RoleID:   3,
	})
	password, err = helpers.HashPassword("userbiasa")
	if err != nil {
		return err
	}

	users = append(users, entities.User{
		Email:    "dummy5@akun.com",
		Name:     "dummy5",
		Password: password,
		RoleID:   3,
	})
	password, err = helpers.HashPassword("userbiasa")
	if err != nil {
		return err
	}

	users = append(users, entities.User{
		Email:    "dummy6@akun.com",
		Name:     "dummy6",
		Password: password,
		RoleID:   3,
	})
	password, err = helpers.HashPassword("userbiasa")
	if err != nil {
		return err
	}

	users = append(users, entities.User{
		Email:    "dummy7@akun.com",
		Name:     "dummy7",
		Password: password,
		RoleID:   3,
	})
	password, err = helpers.HashPassword("userbiasa")
	if err != nil {
		return err
	}

	users = append(users, entities.User{
		Email:    "dummy8@akun.com",
		Name:     "dummy8",
		Password: password,
		RoleID:   3,
	})
	password, err = helpers.HashPassword("userbiasa")
	if err != nil {
		return err
	}

	users = append(users, entities.User{
		Email:    "dummy9@akun.com",
		Name:     "dummy9",
		Password: password,
		RoleID:   3,
	})

	return db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error
}
