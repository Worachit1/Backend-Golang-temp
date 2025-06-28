package student

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"fmt"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.CreateStudent) (*model.Student, bool, error) {

	m := &model.Student{
		First_name:     req.First_name,
		Last_name:      req.Last_name,
		Student_number: req.Student_number,
	}

	_, err := s.db.NewInsert().Model(m).Exec(ctx)

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
