package emergency

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"fmt"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.CreateEmergency, userID string) (*model.Emergency, error) {
	emergency := &model.Emergency{
		UserID:      userID, 
		Type:        req.Type,
		Status:      "รอการตอบสนอง", // ค่า default
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		MapLink:     req.MapLink,
		// FileURL และ OfficerID จะเป็น nil ตอนสร้างใหม่
	}
	emergency.SetCreatedNow()
	_, err := s.db.NewInsert().Model(emergency).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}

	return emergency, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateEmergency, id request.GetByIDEmergency, officerID string) (*model.Emergency, bool, error) {
	// ตรวจสอบว่า emergency มีอยู่หรือไม่
	currentEmergency := &model.Emergency{}
	err := s.db.NewSelect().Model(currentEmergency).
		Where("id = ? AND deleted_at IS NULL", id.ID).
		Scan(ctx)
	if err != nil {
		return nil, true, err // Not found
	}

	// เตรียม query สำหรับ update
	query := s.db.NewUpdate().Model((*model.Emergency)(nil)).Where("id = ?", id.ID)
	hasUpdates := false

	// อัปเดต officer_id เสมอ (เจ้าหน้าที่ผู้รับเรื่อง)
	query = query.Set("officer_id = ?", officerID)
	currentEmergency.OfficerID = officerID
	hasUpdates = true

	// อัปเดตฟิลด์อื่นๆ ถ้ามีการส่งมา
	if req.Type != nil && *req.Type != "" {
		query = query.Set("type = ?", *req.Type)
		currentEmergency.Type = *req.Type
		hasUpdates = true
	}

	if req.Status != nil && *req.Status != "" {
		query = query.Set("status = ?", *req.Status)
		currentEmergency.Status = *req.Status
		hasUpdates = true
	}

	if req.Title != nil && *req.Title != "" {
		query = query.Set("title = ?", *req.Title)
		currentEmergency.Title = *req.Title
		hasUpdates = true
	}

	if req.Description != nil && *req.Description != "" {
		query = query.Set("description = ?", *req.Description)
		currentEmergency.Description = *req.Description
		hasUpdates = true
	}

	if req.Location != nil && *req.Location != "" {
		query = query.Set("location = ?", *req.Location)
		currentEmergency.Location = *req.Location
		hasUpdates = true
	}

	if req.MapLink != nil {
		query = query.Set("map_link = ?", *req.MapLink)
		currentEmergency.MapLink = *req.MapLink
		hasUpdates = true
	}

	if req.ActionNote != nil {
		query = query.Set("action_note = ?", *req.ActionNote)
		currentEmergency.ActionNote = *req.ActionNote
		hasUpdates = true
	}

	if !hasUpdates {
		return currentEmergency, false, nil
	}

	// Set updated_at
	currentEmergency.SetUpdateNow()
	query = query.Set("updated_at = ?", currentEmergency.UpdatedAt)

	_, err = query.Returning("*").Exec(ctx)
	if err != nil {
		return nil, false, err
	}

	return currentEmergency, false, nil
}

// UpdateByOfficer - สำหรับเจ้าหน้าที่อัปเดตเฉพาะฟิลด์ที่เกี่ยวข้อง
func (s *Service) UpdateByOfficer(ctx context.Context, req request.UpdateEmergencyByOfficer, id request.GetByIDEmergency, officerID string) (*model.Emergency, bool, error) {
	logger.Infof("UpdateByOfficer - Looking for emergency ID: %s", id.ID)

	// ตรวจสอบว่า emergency มีอยู่หรือไม่
	currentEmergency := &model.Emergency{}
	err := s.db.NewSelect().Model(currentEmergency).
		Where("id = ? AND deleted_at IS NULL", id.ID).
		Scan(ctx)
	if err != nil {
		logger.Errf("Emergency not found in database: %s, error: %v", id.ID, err)
		return nil, true, err // Not found
	}

	logger.Infof("Found emergency: %+v", currentEmergency)

	// เตรียม query สำหรับ update
	query := s.db.NewUpdate().Model(currentEmergency).Where("id = ?", id.ID)

	hasUpdates := false

	// อัปเดต officer_id เสมอ (เจ้าหน้าที่ผู้รับเรื่อง)
	query = query.Set("officer_id = ?", officerID)
	currentEmergency.OfficerID = officerID
	hasUpdates = true

	// อัปเดตสถานะถ้ามีการส่งมา
	if req.Status != nil && *req.Status != "" {
		logger.Infof("Updating status to: %s", *req.Status)
		query = query.Set("status = ?", *req.Status)
		currentEmergency.Status = *req.Status
		hasUpdates = true
	}

	// อัปเดตบันทึกการดำเนินการถ้ามีการส่งมา
	if req.ActionNote != nil {
		logger.Infof("Updating action_note to: %s", *req.ActionNote)
		query = query.Set("action_note = ?", *req.ActionNote)
		currentEmergency.ActionNote = *req.ActionNote
		hasUpdates = true
	}

	if !hasUpdates {
		return currentEmergency, false, nil
	}

	// Set updated_at
	currentEmergency.SetUpdateNow()
	query = query.Set("updated_at = ?", currentEmergency.UpdatedAt)

	logger.Infof("Executing update query for emergency ID: %s", id.ID)
	_, err = query.Returning("*").Exec(ctx)
	if err != nil {
		logger.Errf("Database update failed: %v", err)
		return nil, false, err
	}

	logger.Infof("Update successful for emergency ID: %s", id.ID)
	return currentEmergency, false, nil
}

func (s *Service) ListEmergencies(ctx context.Context, req request.ListEmergency) ([]response.ListEmergency, int, error) {
	offset := (req.Page - 1) * req.Size
	if offset < 0 {
		offset = 0
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	m := []response.ListEmergency{}

	query := s.db.NewSelect().
		TableExpr("emergencies as e").
		Column("e.id", "e.user_id", "e.officer_id", "e.type", "e.status", "e.title", "e.description", "e.location", "e.map_link", "e.action_note", "e.created_at", "e.updated_at").
		Where("e.deleted_at IS NULL")

	// Filtering by search on allowed fields
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"

		allowedSearchBy := map[string]bool{
			"title":       true,
			"description": true,
			"type":        true,
			"status":      true,
			"location":    true,
		}

		searchBy := "title"
		if allowedSearchBy[strings.ToLower(req.SearchBy)] {
			searchBy = strings.ToLower(req.SearchBy)
		}

		query = query.Where(fmt.Sprintf("LOWER(e.%s) LIKE ?", searchBy), search)
	}

	// Count total before pagination
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Sorting whitelist
	allowedSortBy := map[string]bool{
		"id": true, "type": true, "status": true, "title": true, "created_at": true, "updated_at": true,
	}
	allowedOrderBy := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	sortBy := "id"
	if allowedSortBy[strings.ToLower(req.SortBy)] {
		sortBy = strings.ToLower(req.SortBy)
	}

	orderBy := "asc"
	if allowedOrderBy[strings.ToLower(req.OrderBy)] {
		orderBy = strings.ToLower(req.OrderBy)
	}

	order := fmt.Sprintf("e.%s %s", sortBy, orderBy)

	// Final query with order + pagination
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	return m, count, nil
}


func (s *Service) GetByUserIDEmergency(ctx context.Context, req request.GetByUserIDEmergency) ([]model.Emergency, error) {
	var emergencies []model.Emergency
	err := s.db.NewSelect().
		Model(&emergencies).
		Where("user_id = ? AND deleted_at IS NULL", req.UserID).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		logger.Errf("Database query failed: %v", err)  // เพิ่มบรรทัดนี้
		return nil, err
	}
	return emergencies, nil
}
