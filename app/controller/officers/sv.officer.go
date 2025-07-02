package officers

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Create(ctx context.Context, req request.CreateOfficer) (*model.Officer, bool, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, false, err
	}
	m := &model.Officer{
		FirstName:     req.FirstName,	
		LastName:      req.LastName,
		Email:         req.Email,
		Phone:         req.Phone,
		Password:      string(bytes),
	}
	_, err = s.db.NewInsert().Model(m).
		Column("first_name", "last_name", "email", "phone", "password").
		Returning("*").
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {	
			return nil, true, errors.New("officer already exists")
		}
	}
	return m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateOfficer, id request.GetByIdOfficer) (*model.Officer, bool, error) {
	ex, err := s.db.NewSelect().Table("officers").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, false, err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, false, err
	}
	m := &model.Officer{
		ID:          id.ID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		Password:    string(bytes),
	}
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("first_name = ?first_name").
		Set("last_name = ?last_name").
		Set("email = ?email").
		Set("phone = ?phone").
		Set("password = ?password").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("officer already exists")
		}
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListOfficer) ([]response.ListOfficer, int, error) {
	m := []response.ListOfficer{}

	var (
		offset = (req.Page - 1) * req.Size
		limit  = req.Size
	)

	query := s.db.NewSelect().TableExpr("officers as o").
		Column("o.id", "o.first_name", "o.last_name", "o.email", "o.phone", "o.created_at", "o.updated_at")

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		query.Where("LOWER(o.first_name) LIKE ? OR LOWER(o.last_name) LIKE ?", search, search)
	}

	count, err := query.Count(ctx)
	if count == 0 {
		return m, 0, err
	}

	order := fmt.Sprintf("%s %s", req.SortBy, req.OrderBy)
	err = query.Offset(offset).Limit(limit).Order(order).Scan(ctx, &m)

	return m, count, err
}

func (s *Service) Get(ctx context.Context, id request.GetByIdOfficer) (*response.ListOfficer, error) {
	m := &response.ListOfficer{}

	err := s.db.NewSelect().TableExpr("officers as o").
		Column("o.id", "o.first_name", "o.last_name", "o.email", "o.phone", "o.created_at", "o.updated_at").
		Where("o.id = ?", id.ID).
		Scan(ctx, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) Delete(ctx context.Context, id request.GetByIdOfficer) error {
	ex, err := s.db.NewSelect().Table("officers").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New("officer not found")
	}

	_, err = s.db.NewDelete().Model((*model.Officer)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}
