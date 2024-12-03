package services

import (
	"context"
	"github.com/iota-agency/iota-sdk/pkg/domain/entities/tab"
)

type TabService struct {
	repo tab.Repository
}

func NewTabService(repo tab.Repository) *TabService {
	return &TabService{repo}
}

func (s *TabService) GetByID(ctx context.Context, id uint) (*tab.Tab, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TabService) GetAll(ctx context.Context, params *tab.FindParams) ([]*tab.Tab, error) {
	return s.repo.GetAll(ctx, params)
}

func (s *TabService) GetUserTabs(ctx context.Context, userID uint) ([]*tab.Tab, error) {
	return s.repo.GetUserTabs(ctx, userID)
}

func (s *TabService) Create(ctx context.Context, data *tab.CreateDTO) (*tab.Tab, error) {
	entity, err := data.ToEntity()
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return entity, err
	}
	return entity, nil
}

func (s *TabService) CreateMany(ctx context.Context, data []*tab.CreateDTO) ([]*tab.Tab, error) {
	entities := make([]*tab.Tab, 0, len(data))
	for _, d := range data {
		entity, err := s.Create(ctx, d)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

func (s *TabService) Update(ctx context.Context, id uint, data *tab.UpdateDTO) error {
	entity, err := data.ToEntity(id)
	if err != nil {
		return err
	}
	if err := s.repo.Update(ctx, entity); err != nil {
		return err
	}
	return nil
}

func (s *TabService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *TabService) DeleteUserTabs(ctx context.Context, userID uint) error {
	return s.repo.DeleteUserTabs(ctx, userID)
}