package registration

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

func (s *Service) Create(ctx context.Context, req request.CreateRegistration) (*model.Registration, bool, error) {

	m := &model.Registration{
		ActivitiesID: req.ActivitiesID,
		StudentsID:   req.StudentsID,
	}

	_, err := s.db.NewInsert().Model(m).Exec(ctx)

	return m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateRegistration, id request.GetByIDRegistration) (*model.Registration, bool, error) {
	ex, err := s.db.NewSelect().Table("registrations").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, false, err
	}

	m := &model.Registration{
		ActivitiesID: req.ActivitiesID,
		StudentsID:   req.StudentsID,
	}
	logger.Info(m)
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("activities_id = ?activities_id").
		Set("students_id = ?students_id").
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

func (s *Service) List(ctx context.Context, req request.ListRegistration) ([]response.ListRegistration, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []response.ListRegistration{}
	query := s.db.NewSelect().
		TableExpr("registrations AS r").
		Column("r.id", "r.activity_id", "r.student_id", "r.created_at", "r.updated_at").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprintf("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			search := strings.ToLower(req.Search)
			query.Where(fmt.Sprintf("LOWER(r.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(name) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("r.%s %s", req.SortBy, req.OrderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err
}

func (s *Service) Get(ctx context.Context, id request.GetByIDRegistration) (*response.ListRegistration, error) {
	m := response.ListRegistration{}
	err := s.db.NewSelect().
		TableExpr("registrations AS r").
		Column("r.id", "r.activity_id", "r.student_id", "r.created_at", "r.updated_at").
		Where("id = ?", id.ID).Where("deleted_at IS NULL").Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDRegistration) error {
	ex, err := s.db.NewSelect().Table("registrations").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
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
