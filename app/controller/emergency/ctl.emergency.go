package emergency

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
)
type CreateEmergencyRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Location    string `json:"location" binding:"required"`
	MapLink     string `json:"map_link"`
}


func (ctl *Controller) Create(ctx *gin.Context) {
    var req CreateEmergencyRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        response.BadRequest(ctx, "กรุณากรอกข้อมูลให้ครบถ้วน")
        return
    }

    emergency, err := ctl.Service.Create(ctx, request.CreateEmergency{
        Type:        req.Type,
        Title:       req.Title,
        Description: req.Description,
        Location:    req.Location,
        MapLink:     req.MapLink,
    }, req.UserID)

    if err != nil {
        logger.Err(err.Error())
        response.InternalError(ctx, "ไม่สามารถสร้างรายงานฉุกเฉินได้")
        return
    }

    response.Success(ctx, emergency)
}
func (ctl *Controller) Update(ctx *gin.Context) {
	var req request.UpdateEmergency
	var idReq request.GetByIDEmergency

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid input")
		return
	}

	if err := ctx.ShouldBindUri(&idReq); err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	// Get user ID from JWT token
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "User not authenticated")
		return
	}

	userIDStr := userID.(string)

	emergency, notFound, err := ctl.Service.Update(ctx, req, idReq, userIDStr)
	if err != nil {
		if notFound {
			response.NotFound(ctx, err.Error())
			return
		}
		response.InternalError(ctx, err.Error())
		return
	}

	response.Success(ctx, emergency)
}

// UpdateByOfficer - สำหรับเจ้าหน้าที่อัปเดต
func (ctl *Controller) UpdateByOfficer(ctx *gin.Context) {
    emergencyID := ctx.Param("id")
    if emergencyID == "" {
        response.BadRequest(ctx, "กรุณาระบุ emergency ID")
        return
    }

    var req request.UpdateEmergencyByOfficer
    if err := ctx.ShouldBindJSON(&req); err != nil {
        response.BadRequest(ctx, "ข้อมูล JSON ไม่ถูกต้อง: "+err.Error())
        return
    }

    if req.OfficerID == nil || *req.OfficerID == "" {
        response.BadRequest(ctx, "กรุณาระบุ officer ID")
        return
    }

    logger.Infof("UpdateByOfficer - Emergency ID: %s, Officer ID: %s, Status: %v, ActionNote: %v",
        emergencyID, *req.OfficerID, req.Status, req.ActionNote)

    idReq := request.GetByIDEmergency{ID: emergencyID}

    emergency, notFound, err := ctl.Service.UpdateByOfficer(ctx, req, idReq, *req.OfficerID)
    if err != nil {
        if notFound {
            response.NotFound(ctx, "ไม่พบรายงานฉุกเฉิน")
            return
        }
        response.InternalError(ctx, "ไม่สามารถอัปเดตรายงานฉุกเฉินได้: "+err.Error())
        return
    }

    response.Success(ctx, emergency)
}


// List - แสดง emergency ทั้งหมด
func (ctl *Controller) List(ctx *gin.Context) {
	emergencies := []model.Emergency{}
	err := ctl.Service.db.NewSelect().Model(&emergencies).Where("deleted_at IS NULL").Scan(ctx)
	if err != nil {
		logger.Errf("Failed to list emergencies: %v", err)
		response.InternalError(ctx, "ไม่สามารถดึงข้อมูลได้")
		return
	}

	response.Success(ctx, emergencies)
}


func (ctl *Controller) GetByUserIDEmergency(ctx *gin.Context) {
	var req request.GetByUserIDEmergency
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.BadRequest(ctx, "user_id ไม่ถูกต้อง")
		return
	}

	emergencies, err := ctl.Service.GetByUserIDEmergency(ctx, req)
	if err != nil {
		logger.Errf("Failed to get emergencies by user_id: %v", err)
		response.InternalError(ctx, "ไม่สามารถดึงข้อมูลได้")
		return
	}

	response.Success(ctx, emergencies)
}
// func (ctl *Controller) Get(ctx *gin.Context) {
// 	var req request.GetByIDEmergency
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		response.BadRequest(ctx, "Invalid ID")
// 		return
// 	}

// 	emergency, err := ctl.Service.Get(ctx, req)
// 	if err != nil {
// 		response.NotFound(ctx, "Emergency not found")
// 		return
// 	}

// 	response.Success(ctx, emergency)
// }

// func (ctl *Controller) Delete(ctx *gin.Context) {
// 	var req request.GetByIDEmergency
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		response.BadRequest(ctx, "Invalid ID")
// 		return
// 	}

// 	// Get user ID from JWT token
// 	userID, exists := ctx.Get(jwt.ContextUser)
// 	if !exists {
// 		response.Unauthorized(ctx, "User not authenticated")
// 		return
// 	}

// 	user := userID.(map[string]interface{})
// 	userIDStr := user["user_id"].(string)

// 	err := ctl.Service.Delete(ctx, req, userIDStr)
// 	if err != nil {
// 		if err.Error() == "emergency not found" || err.Error() == "unauthorized to delete this emergency" {
// 			response.NotFound(ctx, err.Error())
// 			return
// 		}
// 		response.InternalError(ctx, err.Error())
// 		return
// 	}

// 	response.Success(ctx, gin.H{"message": "Emergency deleted successfully"})
// }
