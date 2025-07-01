package login

import (
	"app/app/model"
	"context"
	"errors"
)

func (s *Service) Login(ctx context.Context, student_number, password string) (*model.Student, error) {
	student := new(model.Student) // สร้าง user instance ก่อน

	err := s.db.NewSelect().Model(student).
		Where("student_number = ?", student_number).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !student.CheckPassword(password) {
		return nil, errors.New("invalid password")
	}

	return student, nil
}

