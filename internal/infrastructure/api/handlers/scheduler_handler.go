package handlers

import (
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/dtos"
	"devfest/internal/infrastructure/api/response"
	"devfest/internal/infrastructure/api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SchedulerHandler struct {
	usecase domain.SchedulerUsecase
}

func NewSchedulerHandler(usecase domain.SchedulerUsecase) *SchedulerHandler {
	return &SchedulerHandler{
		usecase: usecase,
	}
}

// GetAllByTrack
func (sh *SchedulerHandler) GetAllByTrack(ctx *gin.Context) {
	trackId, err := uuid.Parse(ctx.Param("track-id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Track ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	schedulers, err := sh.usecase.GetAllByTrack(ctx, trackId)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	schedulersResponse := make([]dtos.SchedulerResponse, len(schedulers))
	for i, scheduler := range schedulers {
		schedulersResponse[i] = *mapToSchedulerResponse(&scheduler)
	}

	ctx.JSON(http.StatusOK, schedulersResponse)
}

// GetById
func (sh *SchedulerHandler) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	scheduler, err := sh.usecase.GetByID(ctx, id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	schedulerResponse := mapToSchedulerResponse(scheduler)

	ctx.JSON(http.StatusOK, schedulerResponse)
}

// Create
func (sh *SchedulerHandler) Create(ctx *gin.Context) {
	var req dtos.CreateSchedulerDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Request", err)
		response.HandleError(ctx, newErr)
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		newErr := domain.NewAppError(domain.TypeInternal, "Error parsing startDate", err)
		response.HandleError(ctx, newErr)
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		newErr := domain.NewAppError(domain.TypeInternal, "Error parsing endDate", err)
		response.HandleError(ctx, newErr)
		return
	}
	scheduler := &domain.Scheduler{
		Track:     domain.Track{ID: req.TrackID},
		Talk:      domain.Talk{ID: req.TalkID},
		StartTime: startTime,
		EndTime:   endTime,
		Room:      *req.Room,
	}

	scheduler, err = sh.usecase.Create(ctx, scheduler)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, mapToSchedulerResponse(scheduler))
}

// Update
func (sh *SchedulerHandler) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	var req dtos.UpdateSchedulerDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Request", err)
		response.HandleError(ctx, newErr)
		return
	}

	uid, err := utils.GetUserId(ctx)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	upScheduler := &domain.UpdateScheduler{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Room:      req.Room,
		UpdatedBy: uid,
	}

	scheduler, err := sh.usecase.Update(ctx, id, upScheduler)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, mapToSchedulerResponse(scheduler))
}

// Delete
func (sh *SchedulerHandler) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	err = sh.usecase.Delete(ctx, id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// mapToSchedulerResponse
func mapToSchedulerResponse(scheduler *domain.Scheduler) *dtos.SchedulerResponse {
	return &dtos.SchedulerResponse{
		ScheduleID: scheduler.ID,
		TrackID:    scheduler.Track.ID,
		TalkID:     scheduler.Talk.ID,
		StartTime:  scheduler.StartTime,
		EndTime:    scheduler.EndTime,
		Room:       scheduler.Room,
		UpdatedAt:  scheduler.UpdatedAt,
	}
}
