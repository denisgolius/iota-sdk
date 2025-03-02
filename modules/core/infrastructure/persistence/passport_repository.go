package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/core/domain/entities/passport"
	"github.com/iota-uz/iota-sdk/modules/core/infrastructure/persistence/models"
	"github.com/iota-uz/iota-sdk/pkg/composables"
)

var (
	ErrPassportNotFound = errors.New("passport not found")
)

const (
	selectPassportQuery = `
		SELECT 
			id,
			first_name,
			last_name,
			middle_name,
			gender,
			birth_date,
			birth_place,
			nationality,
			passport_type,
			passport_number,
			series,
			issuing_country,
			issued_at,
			issued_by,
			expires_at,
			machine_readable_zone,
			biometric_data,
			signature_image,
			remarks,
			created_at,
			updated_at
		FROM passports
	`
	insertPassportQuery = `
		INSERT INTO passports (
			id,
			first_name,
			last_name,
			middle_name,
			gender,
			birth_date,
			birth_place,
			nationality,
			passport_type,
			passport_number,
			series,
			issuing_country,
			issued_at,
			issued_by,
			expires_at,
			machine_readable_zone,
			biometric_data,
			signature_image,
			remarks
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) RETURNING id
	`
	updatePassportQuery = `
		UPDATE passports 
		SET first_name = $1,
			last_name = $2, 
			middle_name = $3, 
			gender = $4,
			birth_date = $5,
			birth_place = $6,
			nationality = $7,
			passport_type = $8,
			passport_number = $9,
			series = $10,
			issuing_country = $11,
			issued_at = $12,
			issued_by = $13,
			expires_at = $14,
			machine_readable_zone = $15,
			biometric_data = $16,
			signature_image = $17,
			remarks = $18,
			updated_at = $19
		WHERE id = $20
	`
	deletePassportQuery = `DELETE FROM passports WHERE id = $1`
)

type PassportRepository struct{}

func NewPassportRepository() passport.PassportRepository {
	return &PassportRepository{}
}

func (r *PassportRepository) queryPassports(ctx context.Context, query string, args ...interface{}) ([]passport.Passport, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passports []passport.Passport
	for rows.Next() {
		var p models.Passport

		if err := rows.Scan(
			&p.ID,
			&p.FirstName,
			&p.LastName,
			&p.MiddleName,
			&p.Gender,
			&p.BirthDate,
			&p.BirthPlace,
			&p.Nationality,
			&p.PassportType,
			&p.PassportNumber,
			&p.Series,
			&p.IssuingCountry,
			&p.IssuedAt,
			&p.IssuedBy,
			&p.ExpiresAt,
			&p.MachineReadableZone,
			&p.BiometricData,
			&p.SignatureImage,
			&p.Remarks,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Convert DB model to domain model
		passportEntity, err := ToDomainPassport(&p)
		if err != nil {
			return nil, err
		}
		passports = append(passports, passportEntity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return passports, nil
}

func (r *PassportRepository) Create(ctx context.Context, data passport.Passport) (passport.Passport, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	dbRow, err := ToDBPassport(data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert passport to db model: %w", err)
	}

	var id string
	err = pool.QueryRow(
		ctx,
		insertPassportQuery,
		dbRow.ID,
		dbRow.FirstName,
		dbRow.LastName,
		dbRow.MiddleName,
		dbRow.Gender,
		dbRow.BirthDate,
		dbRow.BirthPlace,
		dbRow.Nationality,
		dbRow.PassportType,
		dbRow.PassportNumber,
		dbRow.Series,
		dbRow.IssuingCountry,
		dbRow.IssuedAt,
		dbRow.IssuedBy,
		dbRow.ExpiresAt,
		dbRow.MachineReadableZone,
		dbRow.BiometricData,
		dbRow.SignatureImage,
		dbRow.Remarks,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create passport: %w", err)
	}

	return r.GetByID(ctx, uuid.MustParse(id))
}

func (r *PassportRepository) GetByID(ctx context.Context, id uuid.UUID) (passport.Passport, error) {
	passports, err := r.queryPassports(ctx, selectPassportQuery+" WHERE id = $1", id.String())
	if err != nil {
		return nil, err
	}
	if len(passports) == 0 {
		return nil, ErrPassportNotFound
	}
	return passports[0], nil
}

func (r *PassportRepository) GetByPassportNumber(ctx context.Context, series, number string) (passport.Passport, error) {
	passports, err := r.queryPassports(ctx, selectPassportQuery+" WHERE series = $1 AND passport_number = $2", series, number)
	if err != nil {
		return nil, err
	}
	if len(passports) == 0 {
		return nil, ErrPassportNotFound
	}
	return passports[0], nil
}

func (r *PassportRepository) Update(ctx context.Context, id uuid.UUID, data passport.Passport) (passport.Passport, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	dbRow, err := ToDBPassport(data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert passport to db model: %w", err)
	}

	_, err = pool.Exec(
		ctx,
		updatePassportQuery,
		dbRow.FirstName,
		dbRow.LastName,
		dbRow.MiddleName,
		dbRow.Gender,
		dbRow.BirthDate,
		dbRow.BirthPlace,
		dbRow.Nationality,
		dbRow.PassportType,
		dbRow.PassportNumber,
		dbRow.Series,
		dbRow.IssuingCountry,
		dbRow.IssuedAt,
		dbRow.IssuedBy,
		dbRow.ExpiresAt,
		dbRow.MachineReadableZone,
		dbRow.BiometricData,
		dbRow.SignatureImage,
		dbRow.Remarks,
		time.Now(),
		id.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update passport: %w", err)
	}

	return r.GetByID(ctx, id)
}

func (r *PassportRepository) Delete(ctx context.Context, id uuid.UUID) error {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, deletePassportQuery, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete passport: %w", err)
	}

	return nil
}

