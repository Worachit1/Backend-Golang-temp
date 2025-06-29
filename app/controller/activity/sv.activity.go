package activity

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"errors"
	"fmt"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.CreateActivity) (*model.Activity, bool, error) {

	m := &model.Activity{
		Name:        req.Name,
		Description: req.Description,
		ReleaseDate: req.ReleaseDate,
	}

	_, err := s.db.NewInsert().Model(m).Exec(ctx)

	return m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateActivity, id request.GetByIDActivity) (*model.Activity, bool, error) {
	ex, err := s.db.NewSelect().Table("activities").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, false, err
	}

	m := &model.Activity{
		ID:          id.ID,
		Name:        req.Name,
		Description: req.Description,
		ReleaseDate: req.ReleaseDate,
	}
	logger.Info(m)
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("name = ?name").
		Set("description = ?description").
		Set("release_date = ?release_date").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("activity with this name already exists")
		}
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListActivity) ([]response.ListActivity, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []response.ListActivity{}
	query := s.db.NewSelect().
		TableExpr("activities AS a").
		Column("a.id", "a.name", "a.description", "a.release_date", "a.created_at", "a.updated_at").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprintf("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			search := strings.ToLower(req.Search)
			query.Where(fmt.Sprintf("LOWER(a.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(name) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("a.%s %s", req.SortBy, req.OrderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err
}

func (s *Service) Get(ctx context.Context, id request.GetByIDActivity) (*response.ListActivity, error) {
	m := response.ListActivity{}
	err := s.db.NewSelect().
		TableExpr("activities AS a").
		Column("a.id", "a.name", "a.description", "a.release_date", "a.created_at", "a.updated_at").
		Where("id = ?", id.ID).Where("deleted_at IS NULL").Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDActivity) error {
	ex, err := s.db.NewSelect().Table("activities").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New("user not found")
	}

	// data, err := s.db.NewDelete().Table("users").Where("id = ?", id.ID).Exec(ctx)
	_, err = s.db.NewDelete().Model((*model.Activity)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}
