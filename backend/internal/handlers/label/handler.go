package label

import (
	"errors"
	"fmt"
	"net/http"
	"reports_system/internal/handlers/middleware"
	"reports_system/internal/mapper"
	"reports_system/internal/model/label"
	"reports_system/internal/model/report"
	"reports_system/internal/service"
	"reports_system/pkg/e"
	"reports_system/pkg/logging"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	apiURLGroup     = "/api"
	reportsURLGroup = "/reports"
	labelsURLGroup  = "/labels"
	apiVersion      = "1"
)

type Handler struct {
	logger  logging.Logger
	service service.Label
	mapper  mapper.Label
}

func NewHandler(logger logging.Logger, service service.Label, mapper mapper.Label) *Handler {
	return &Handler{logger: logger, service: service, mapper: mapper}
}

func (h *Handler) Register(router *gin.Engine) {
	labelsGroupName := fmt.Sprintf("%v/v%v%v", apiURLGroup, apiVersion, labelsURLGroup)
	labelsOnReportGroupName := fmt.Sprintf("%v/v%v%v/:id%v", apiURLGroup, apiVersion, reportsURLGroup, labelsURLGroup)

	h.logger.Tracef("Register route: %v", labelsGroupName)
	h.logger.Tracef("Register route: %v", labelsOnReportGroupName)

	labelsGroup := router.Group(labelsGroupName, middleware.Authenticate)
	{
		labelsGroup.GET("", h.getAllLabels)
		labelsGroup.GET("/:id", h.getOneLabel)
		labelsGroup.PATCH("/:id", h.updateLabel)
		labelsGroup.DELETE("/:id", h.deleteLabel)
	}

	h.logger.Tracef("Register route: %v", labelsOnReportGroupName)
	labelsOnReportGroup := router.Group(labelsOnReportGroupName, middleware.Authenticate)
	{
		labelsOnReportGroup.GET("", h.getAllLabelsOnReport)
		labelsOnReportGroup.POST("", h.createLabel)             // /api/reports/:id/labels/
		labelsOnReportGroup.DELETE("/:label_id", h.detachLabel) // /api/reports/:id/labels/
	}
}

// @Summary Create label
// @Security ApiKeyAuth
// @Tags labels
// @Description create label
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "id"
// @Param dto body label.CreateLabelDTO true "label info"
// @Success 201 {string} string 1
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/reports/{id}/labels [post]
func (h *Handler) createLabel(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	reportID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	var (
		dto label.CreateLabelDTO
		t   label.Label
	)

	if err := ctx.BindJSON(&dto); err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	t = h.mapper.MapCreateLabelDTO(dto)
	err = h.service.Create(userID, reportID, &t)

	if err != nil {
		if errors.Is(err, &report.ReportNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, fmt.Sprintf(
		"%s%s/%v", apiURLGroup, labelsURLGroup, t.ID))
}

func (h *Handler) getAllLabelsOnReport(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	reportID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	labels, err := h.service.GetAllByReport(userID, reportID)

	if err != nil {
		h.logger.Info(err)
		if errors.Is(err, &report.ReportNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	dto := h.mapper.MapGetAllLabelsDTO(labels)

	ctx.JSON(http.StatusOK, dto)
}

// @Summary Get all labels
// @Security ApiKeyAuth
// @Tags labels
// @Parametrs
// @Description get labels from user
// @Accept  json
// @Produce  json
// @Success 200 {object} label.GetAllLabelsDTO
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/labels [get]
func (h *Handler) getAllLabels(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	labels, err := h.service.GetAll(userID)

	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	dto := h.mapper.MapGetAllLabelsDTO(labels)

	ctx.JSON(http.StatusOK, dto)
}

func (h *Handler) getOneLabel(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	labelID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	t, err := h.service.GetOne(userID, labelID)

	if err != nil {
		h.logger.Info(err)
		if errors.Is(err, &label.LabelNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, t)
}

// @Summary Update label by ID
// @Security ApiKeyAuth
// @Tags labels
// @Description update one label by ID
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "id"
// @Param dto body label.UpdateLabelDTO true "label info"
// @Success 204
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/labels/{id} [patch]
func (h *Handler) updateLabel(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return
	}

	labelID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	var (
		dto label.UpdateLabelDTO
	)
	if err := ctx.BindJSON(&dto); err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	t := h.mapper.MapUpdateLabelDTO(dto)
	err = h.service.Update(userID, labelID, t)
	if err != nil {
		h.logger.Info(err)
		if errors.Is(err, &label.LabelNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

// @Summary Delete one label by ID
// @Security ApiKeyAuth
// @Tags labels
// @Description delete one label by ID
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "id"
// @Success 200 {integer} integer 1
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/labels/{id} [delete]
func (h *Handler) deleteLabel(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	labelID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.service.Delete(userID, labelID)

	if err != nil {
		h.logger.Info(err)
		if errors.Is(err, &label.LabelNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

// @Summary Detach label by ID from report by ID
// @Security ApiKeyAuth
// @Tags reports
// @Description detach label by ID from report by ID
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "id"
// @Param   label_id  path  string  true  "label id"
// @Success 200 {integer} integer 1
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/reports/{id}/labels/{label_id} [delete]
func (h *Handler) detachLabel(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	reportID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	labelID, err := strconv.Atoi(ctx.Param("label_id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.service.Detach(userID, labelID, reportID)

	if err != nil {
		h.logger.Info(err)
		if errors.Is(err, &label.LabelNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)

}
