package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"neatly/internal/handlers/middleware"
	"neatly/internal/mapper"
	"neatly/internal/model/account"
	"neatly/internal/model/report"
	"neatly/internal/service"
	"neatly/pkg/e"
	"neatly/pkg/logging"
	"net/http"
	"strconv"
)

const (
	labelsURLGroup  = "/labels"
	reportsURLGroup = "/reports"
	apiURLGroup     = "/api"
	apiVersion      = "1"
	searchURL       = "/search"
	labelSearchKey  = "label"
)

type Handler struct {
	logger  logging.Logger
	service service.Report
	mapper  mapper.Report
}

func NewHandler(logger logging.Logger, service service.Report, mapper mapper.Report) *Handler {
	return &Handler{logger: logger, service: service, mapper: mapper}
}

func (h *Handler) Register(router *gin.Engine) {
	groupName := fmt.Sprintf("%v/v%v%v", apiURLGroup, apiVersion, reportsURLGroup)

	h.logger.Tracef("Register route: %v", groupName)

	group := router.Group(groupName, middleware.Authenticate)
	{
		group.GET("", h.getAllReports)       // /api/v1/reports
		group.POST("", h.createReport)       // /api/v1/reports
		group.GET("/:id", h.getOneReport)    // /api/v1/reports/:id
		group.PATCH("/:id", h.updateReport)  // /api/v1/reports/:id
		group.DELETE("/:id", h.deleteReport) // /api/v1/reports/:id
	}
}

// @Summary Create report
// @Security ApiKeyAuth
// @Labels reports
// @Description create report
// @Accept  json
// @Produce  json
// @Param dto body report.CreateReportDTO true "report content"
// @Success 201 {string} string 1
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/reports [post]
func (h *Handler) createReport(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		h.logger.Info(err)
		return
	}

	var dto report.CreateReportDTO
	if err := ctx.BindJSON(&dto); err != nil {
		h.logger.Info(err)
		if errors.Is(err, &report.ReportNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	n := h.mapper.MapCreateReportDTO(dto)
	err = h.service.Create(userID, &n)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, fmt.Sprintf(
		"%s%s/%v", apiURLGroup, reportsURLGroup, n.ID))
}

// @Summary Get all reports from user filter by label
// @Security ApiKeyAuth
// @Labels reports
// @Description create report
// @Accept  json
// @Produce  json
// @Param   label query  string  false  "reports search by label"
// @Success 200 {object} report.GetAllReportsDTO
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router  /api/v1/reports [get]
func (h *Handler) getAllReports(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		if errors.Is(err, &account.AccountNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	var ns []report.Report

	keys := ctx.Request.URL.Query()
	values := keys[labelSearchKey]
	if values == nil {
		ns, err = h.service.GetAll(userID)
		if err != nil {
			if errors.Is(err, &report.ReportNotFoundErr{}) {
				e.NewErrorResponse(ctx, http.StatusNotFound, err)
				return
			}
			e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		ns, err = h.service.FindByLabels(userID, values)
		if err != nil {
			h.logger.Info(err)
			e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	dto := h.mapper.MapGetAllReportsDTO(ns)

	ctx.JSON(http.StatusOK, dto)
}

// @Summary Get Report By Id
// @Security ApiKeyAuth
// @Labels reports
// @Description get report by id
// @ID get-report-by-id
// @Accept  json
// @Produce json
// @Param   id  path  string  true  "id"
// @Success 200 {object} report.Report
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/reports/{id} [get]
func (h *Handler) getOneReport(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	reportID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	n, err := h.service.GetOne(userID, reportID)
	if err != nil {
		if errors.Is(err, &report.ReportNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, n)
}

// @Summary Update Report
// @Security ApiKeyAuth
// @Labels reports
// @Description update report
// @ID update-report
// @Accept  json
// @Produce json
// @Param   id   path  string  true  "id"
// @Param dto body report.UpdateReportDTO true "report content"
// @Success 204
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/reports/{id} [patch]
func (h *Handler) updateReport(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	reportID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	h.logger.Debug("unmarshal body bytes")
	var (
		dto            report.UpdateReportDTO
		data           map[string]interface{}
		needBodyUpdate bool
	)
	h.logger.Infof("NOTE ID: %v", reportID)
	dto.ID = reportID
	if err := json.Unmarshal(bodyBytes, &dto); err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	_, needBodyUpdate = data["body"]
	h.logger.Infof("Need body update: %v", needBodyUpdate)

	n := h.mapper.MapUpdateReportDTO(dto)
	err = h.service.Update(userID, n, needBodyUpdate)

	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

// @Summary Delete Report
// @Security ApiKeyAuth
// @Labels reports
// @Description delete report
// @ID delete-report
// @Accept  json
// @Produce json
// @Param   id   path string  true  "id"
// @Success 204
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/reports/{id} [delete]
func (h *Handler) deleteReport(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	reportID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.service.Delete(userID, reportID)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}
