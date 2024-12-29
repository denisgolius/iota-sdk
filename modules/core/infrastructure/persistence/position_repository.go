package persistence

import (
	"context"
	position2 "github.com/iota-uz/iota-sdk/modules/core/domain/entities/position"
	"github.com/iota-uz/iota-sdk/pkg/composables"
)

type GormPositionRepository struct{}

func NewPositionRepository() position2.Repository {
	return &GormPositionRepository{}
}

func (g *GormPositionRepository) GetPaginated(
	ctx context.Context,
	limit, offset int,
	sortBy []string,
) ([]*position2.Position, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return nil, composables.ErrNoTx
	}
	var uploads []*position2.Position
	q := tx.Limit(limit).Offset(offset)
	for _, s := range sortBy {
		q = q.Order(s)
	}
	if err := q.Find(&uploads).Error; err != nil {
		return nil, err
	}
	return uploads, nil
}

func (g *GormPositionRepository) Count(ctx context.Context) (int64, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return 0, composables.ErrNoTx
	}
	var count int64
	if err := tx.Model(&position2.Position{}).Count(&count).Error; err != nil { //nolint:exhaustruct
		return 0, err
	}
	return count, nil
}

func (g *GormPositionRepository) GetAll(ctx context.Context) ([]*position2.Position, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return nil, composables.ErrNoTx
	}
	var entities []*position2.Position
	if err := tx.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (g *GormPositionRepository) GetByID(ctx context.Context, id int64) (*position2.Position, error) {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return nil, composables.ErrNoTx
	}
	var entity position2.Position
	if err := tx.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (g *GormPositionRepository) Create(ctx context.Context, data *position2.Position) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	return tx.Create(data).Error
}

func (g *GormPositionRepository) Update(ctx context.Context, data *position2.Position) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	return tx.Save(data).Error
}

func (g *GormPositionRepository) Delete(ctx context.Context, id int64) error {
	tx, ok := composables.UseTx(ctx)
	if !ok {
		return composables.ErrNoTx
	}
	return tx.Delete(&position2.Position{}, id).Error //nolint:exhaustruct
}
