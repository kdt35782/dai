package handler

import (
	"github.com/gin-gonic/gin"
	"sm-medical/internal/service"
	"sm-medical/pkg/utils"
	"strconv"
)

type RecordHandler struct {
	service *service.RecordService
}

func NewRecordHandler() *RecordHandler {
	return &RecordHandler{
		service: service.NewRecordService(),
	}
}

// GetList 获取病历列表
func (h *RecordHandler) GetList(c *gin.Context) {
	userID := c.GetInt64("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	list, total, err := h.service.GetList(userID, page, pageSize, startDate, endDate)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetDetail 获取病历详情
func (h *RecordHandler) GetDetail(c *gin.Context) {
	userID := c.GetInt64("userID")
	recordIDStr := c.Query("recordId")

	recordID, err := strconv.ParseInt(recordIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "病历ID格式错误")
		return
	}

	detail, err := h.service.GetDetail(userID, recordID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, detail)
}
