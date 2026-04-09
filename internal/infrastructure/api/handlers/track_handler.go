package handlers

import (
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/dtos"
	"devfest/internal/infrastructure/api/response"
	"devfest/internal/infrastructure/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TrackHandler struct {
	trackService domain.TrackUsecase
}

func NewTrackHandler(trackService domain.TrackUsecase) *TrackHandler {
	return &TrackHandler{
		trackService: trackService,
	}
}

// GetAll
func (th *TrackHandler) GetAll(ctx *gin.Context) {
	eventID, err := uuid.Parse(ctx.Param("event-id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Event ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	tracks, err := th.trackService.GetAll(ctx.Request.Context(), eventID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	trackResponse := make([]dtos.TrackResponse, len(tracks))
	for i, track := range tracks {
		trackResponse[i] = mapToTrackResponse(&track)
	}

	ctx.JSON(http.StatusOK, trackResponse)
}

// GetByID
func (th *TrackHandler) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	track, err := th.trackService.GetById(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	trackResponse := mapToTrackResponse(track)

	ctx.JSON(http.StatusOK, trackResponse)
}

// GetFullEventSchedule
func (th *TrackHandler) GetFullEventSchedule(ctx *gin.Context) {
	eventID, err := uuid.Parse(ctx.Param("event-id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Event ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	tracks, err := th.trackService.GetFullEventSchedule(ctx.Request.Context(), eventID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	fullEventSchedule := make([]dtos.FullTrackScheduleResponse, len(tracks))
	for i, track := range tracks {
		fullEventSchedule[i] = mapToFullTrackScheduleResponse(&track)
	}

	ctx.JSON(http.StatusOK, fullEventSchedule)
}

// Create
func (th *TrackHandler) Create(ctx *gin.Context) {
	var dto dtos.CreateTrackDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid request body", err)
		response.HandleError(ctx, newErr)
		return
	}

	uid, err := utils.GetUserId(ctx)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	track := &domain.Track{
		EventID:   dto.EventID,
		Name:      dto.Name,
		EventDate: dto.EventDate,
		Audit:     domain.Audit{CreatedBy: uid},
	}

	track, err = th.trackService.Create(ctx.Request.Context(), track)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	trackResponse := mapToTrackResponse(track)

	ctx.JSON(http.StatusCreated, trackResponse)
}

// Update
func (th *TrackHandler) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	var req dtos.UpdateTrackDTO
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

	updTrack := &domain.UpdateTrack{
		Name:      req.Name,
		EventDate: req.EventDate,
		UpdatedBy: uid,
	}

	track, err := th.trackService.Update(ctx.Request.Context(), id, updTrack)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	trackResponse := mapToTrackResponse(track)

	ctx.JSON(http.StatusOK, trackResponse)
}

// Delete
func (th *TrackHandler) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	err = th.trackService.Delete(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// mapToTrackResponse
func mapToTrackResponse(track *domain.Track) dtos.TrackResponse {
	return dtos.TrackResponse{
		ID:        track.ID,
		EventID:   track.EventID,
		Name:      track.Name,
		EventDate: track.EventDate,
		CreatedBy: track.CreatedBy,
		UpdatedBy: track.UpdatedBy,
		CreatedAt: track.CreatedAt,
		UpdatedAt: track.UpdatedAt,
	}
}

// mapToFullTrackScheduleResponse
func mapToFullTrackScheduleResponse(track *domain.FullTrackSchedule) dtos.FullTrackScheduleResponse {
	entries := make([]dtos.ScheduleEntryTrackResponse, len(track.Entries))

	for i, entry := range track.Entries {
		speakers := make([]dtos.SpeakerTrackResponse, len(entry.Talk.Speakers))
		for j, s := range entry.Talk.Speakers {
			speakers[j] = dtos.SpeakerTrackResponse(s)
		}

		entries[i] = dtos.ScheduleEntryTrackResponse{
			ScheduleID: entry.ScheduleID,
			StartTime:  entry.StartTime,
			EndTime:    entry.EndTime,
			Room:       entry.Room,
			Talk: dtos.TalkTrackResponse{
				ID:          entry.Talk.ID,
				Title:       entry.Talk.Title,
				Description: entry.Talk.Description,
				Speakers:    speakers,
			},
		}
	}

	return dtos.FullTrackScheduleResponse{
		TrackID:   track.TrackID,
		TrackName: track.TrackName,
		EventDate: track.EventDate,
		Entries:   entries,
	}
}
