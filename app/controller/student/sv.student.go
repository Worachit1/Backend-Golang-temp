package student

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Create(ctx context.Context, req request.CreateStudent) (*model.Student, bool, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, false, err
	}

	m := &model.Student{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		StudentNumber: req.StudentNumber,
		Password:      string(bytes),
	}
	m.SetCreatedNow() // Set created_at and updated_at if your model supports it

	_, err = s.db.NewInsert().Model(m).Returning("*").Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("student number already exists")
		}
		return nil, false, err
	}
	return m, false, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateStudent, id request.GetByIDStudent) (*model.Student, bool, error) {
	ex, err := s.db.NewSelect().Table("students").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, false, err
	}

	m := &model.Student{
		ID:            id.ID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		StudentNumber: req.StudentNumber,
	}
	logger.Info(m)
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("first_name = ?first_name").
		Set("last_name = ?last_name").
		Set("student_number = ?student_number").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("student number already exists")
		}
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListStudent) ([]response.ListStudent, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []response.ListStudent{}
	query := s.db.NewSelect().
		TableExpr("students AS s").
		Column("s.id", "s.first_name", "s.last_name", "s.student_number", "s.created_at", "s.updated_at").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprintf("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			search := strings.ToLower(req.Search)
			query.Where(fmt.Sprintf("LOWER(s.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(first_name) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("s.%s %s", req.SortBy, req.OrderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err
}

func (s *Service) Get(ctx context.Context, id request.GetByIDStudent) (*response.ListStudent, error) {
	m := response.ListStudent{}
	err := s.db.NewSelect().
		TableExpr("students AS s").
		Column("s.id", "s.first_name", "s.last_name", "s.student_number", "s.created_at", "s.updated_at").
		Where("id = ?", id.ID).Where("deleted_at IS NULL").Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDStudent) error {
	ex, err := s.db.NewSelect().Table("students").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New("user not found")
	}

	// data, err := s.db.NewDelete().Table("users").Where("id = ?", id.ID).Exec(ctx)
	_, err = s.db.NewDelete().Model((*model.Student)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}
