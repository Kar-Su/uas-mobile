package utils

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

func TranslateValidationError(err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return err.Error()
	}

	fieldTags := map[string]string{
		"Password":   "password",
		"Email":      "email",
		"Name":       "nama",
		"Satuan":     "satuan",
		"Keterangan": "keterangan",
		"Kode":       "kode barang",
		"RoleName":   "role",
	}

	tagMessages := map[string]string{
		"required": "wajib diisi",
		"min":      "minimal %s karakter",
		"email":    "format email tidak valid",
		"max":      "maksimal %s karakter",
	}

	var msgs []string
	for _, fe := range ve {
		field := fe.Field()
		tag := fe.Tag()
		param := fe.Param()

		label := fieldTags[field]
		if label == "" {
			label = field
		}

		tmpl := tagMessages[tag]
		if tmpl == "" {
			tmpl = "tidak valid"
		}

		if param != "" && strings.Contains(tmpl, "%s") {
			msgs = append(msgs, label+" "+strings.ReplaceAll(tmpl, "%s", param))
		} else {
			msgs = append(msgs, label+" "+tmpl)
		}
	}
	return strings.Join(msgs, ", ")
}

func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return true
	}
	return strings.Contains(err.Error(), "foreign key constraint")
}

func IsNilStruct[T any](s *T) bool {
	if s == nil {
		return true
	}

	v := reflect.ValueOf(s).Elem()

	count := 0
	for _, fieldVal := range v.Fields() {
		if fieldVal.IsZero() {
			count++
		}
	}
	return count == v.NumField()
}
