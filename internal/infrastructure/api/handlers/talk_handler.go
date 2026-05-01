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

type TalkHandler struct {
	talkService domain.TalkUsecase
}

func NewTalkHandler(talkService domain.TalkUsecase) TalkHandler {
	return TalkHandler{talkService: talkService}
}

// GetAll
func (h *TalkHandler) GetAll(ctx *gin.Context) {
	eventID, err := uuid.Parse(ctx.Param("event-id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Event ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	talks, err := h.talkService.GetAll(ctx.Request.Context(), eventID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	talkResponse := make([]dtos.TalkResponse, len(talks))
	for i, talk := range talks {
		talkResponse[i] = mapToTalkResponse(&talk)
	}

	ctx.JSON(http.StatusOK, talkResponse)
}

// GetById
func (h *TalkHandler) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	talk, err := h.talkService.GetById(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	talkResponse := mapToTalkResponse(talk)

	ctx.JSON(http.StatusOK, talkResponse)
}

// Create
func (h *TalkHandler) Create(ctx *gin.Context) {
	var req dtos.CreateTalkDTO
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

	var tags []string
	if req.Tags != nil {
		tags = *req.Tags
	}

	talk := &domain.Talk{
		EventID:     req.EventID,
		Title:       req.Title,
		Description: req.Description,
		Tags:        tags,
		Audit: domain.Audit{
			CreatedBy: uid,
		},
	}

	talk, err = h.talkService.Create(ctx.Request.Context(), talk)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, mapToTalkResponse(talk))
}

// Update
func (h *TalkHandler) Update(ctx *gin.Context) {
	var req dtos.UpdateTalkDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Request", err)
		response.HandleError(ctx, newErr)
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	uid, err := utils.GetUserId(ctx)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var tags []string
	if req.Tags != nil {
		tags = req.Tags
	}

	updTalk := &domain.UpdateTalk{
		Title:       req.Title,
		Description: req.Description,
		Tags:        tags,
		UpdatedBy:   uid,
	}

	talk, err := h.talkService.Update(ctx.Request.Context(), id, updTalk)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, mapToTalkResponse(talk))
}

// Delete
func (h *TalkHandler) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	err = h.talkService.Delete(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// AddSpeaker
func (h *TalkHandler) AddSpeaker(ctx *gin.Context) {
	var req dtos.CreateTalkSpeakerDTO
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

	talkSpeaker := &domain.TalkSpeaker{
		TalkID:    req.TalkID,
		SpeakerID: req.SpeakerID,
		CreatedBy: uid,
	}

	err = h.talkService.AddSpeaker(ctx.Request.Context(), talkSpeaker)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

// RemoveSpeaker
func (h *TalkHandler) RemoveSpeaker(ctx *gin.Context) {
	talkId, err := uuid.Parse(ctx.Query("talk_id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Talk ID", err)
		response.HandleError(ctx, newErr)
		return
	}
	speakerId, err := uuid.Parse(ctx.Query("speaker_id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Speaker ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	talkSpeaker := &domain.TalkSpeaker{
		TalkID:    talkId,
		SpeakerID: speakerId,
	}

	err = h.talkService.RemoveSpeaker(ctx.Request.Context(), talkSpeaker)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// mapToTalkResponse
func mapToTalkResponse(talk *domain.Talk) dtos.TalkResponse {
	speakers := make([]dtos.SpeakerTalkDetailResponse, len(talk.Speakers))
	for i, s := range talk.Speakers {
		speakers[i] = dtos.SpeakerTalkDetailResponse{
			ID:        s.ID,
			FirstName: s.FirstName,
			LastName:  s.LastName,
			Email:     s.Email,
			AvatarURL: s.AvatarURL,
			Company:   s.Company,
			Bio:       s.Bio,
		}
	}
	return dtos.TalkResponse{
		ID:          talk.ID,
		EventID:     talk.EventID,
		Title:       talk.Title,
		Description: talk.Description,
		Tags:        talk.Tags,
		Speakers:    speakers,
	}
}
