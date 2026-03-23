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

type SpeakerHandler struct {
	usecase domain.SpeakerUsecase
}

func NewSpeakerHandler(usecase domain.SpeakerUsecase) *SpeakerHandler {
	return &SpeakerHandler{usecase: usecase}
}

// GetAll
func (s *SpeakerHandler) GetAll(ctx *gin.Context) {
	eventID, err := uuid.Parse(ctx.Param("event-id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Event ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	speakers, err := s.usecase.GetAll(ctx.Request.Context(), eventID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	speakerResponse := make([]dtos.SpeakerDetailResponse, len(speakers))
	for i, speaker := range speakers {
		speakerResponse[i] = mapToSpeakerDetailResponse(&speaker)
	}

	ctx.JSON(http.StatusOK, speakerResponse)
}

// GetById
func (s *SpeakerHandler) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	speaker, err := s.usecase.GetById(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	speakerResponse := mapToSpeakerDetailResponse(speaker)

	ctx.JSON(http.StatusOK, speakerResponse)
}

// ListPaged
func (s *SpeakerHandler) ListPaged(ctx *gin.Context) {
	eventID, err := uuid.Parse(ctx.Param("event-id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Event ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	search := ctx.DefaultQuery("search", "")
	page, pageSize := utils.GetPaginationParams(ctx)
	if pageSize > 100 {
		pageSize = 100
	}

	speakers, total, err := s.usecase.ListPaged(ctx.Request.Context(), eventID, search, int32(page), int32(pageSize))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	speakerResponse := make([]dtos.SpeakerDetailResponse, len(speakers))
	for i, speaker := range speakers {
		speakerResponse[i] = mapToSpeakerDetailResponse(&speaker)
	}

	res := dtos.PagedResponse[dtos.SpeakerDetailResponse]{
		Data: speakerResponse,
		Meta: dtos.PagedMeta{
			Total:    total,
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
	}

	ctx.JSON(http.StatusOK, res)
}

// Create
func (s *SpeakerHandler) Create(ctx *gin.Context) {
	var dto dtos.CreateSpeakerDTO
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

	speaker := &domain.Speaker{
		EventID: dto.EventID,
		Company: dto.Company,
		Bio:     dto.Bio,
		Person:  domain.Person{ID: dto.PersonID},
		Audit:   domain.Audit{CreatedBy: uid},
	}

	speaker, err = s.usecase.Create(ctx.Request.Context(), speaker)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	speakerResponse := mapToSpeakerResponse(speaker)

	ctx.JSON(http.StatusCreated, speakerResponse)
}

// Update
func (s *SpeakerHandler) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	var dto dtos.UpdateSpeakerDTO
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

	updSpeaker := domain.UpdateSpeaker{
		Company:   dto.Company,
		Bio:       dto.Bio,
		UpdatedBy: uid,
	}

	speaker, err := s.usecase.Update(ctx.Request.Context(), id, &updSpeaker)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	speakerResponse := mapToSpeakerResponse(speaker)

	ctx.JSON(http.StatusOK, speakerResponse)
}

// Delete
func (s *SpeakerHandler) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	err = s.usecase.Delete(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func mapToSpeakerResponse(speaker *domain.Speaker) dtos.SpeakerResponse {
	return dtos.SpeakerResponse{
		ID:        speaker.ID,
		EventID:   speaker.EventID,
		PersonID:  speaker.Person.ID,
		Company:   utils.SafeString(speaker.Company),
		Bio:       utils.SafeString(speaker.Bio),
		CreatedAt: speaker.Audit.CreatedAt,
		UpdatedAt: speaker.Audit.UpdatedAt,
		CreatedBy: speaker.Audit.CreatedBy,
		UpdatedBy: speaker.Audit.UpdatedBy,
	}
}

func mapToSpeakerDetailResponse(speaker *domain.Speaker) dtos.SpeakerDetailResponse {
	return dtos.SpeakerDetailResponse{
		ID:          speaker.ID,
		EventID:     speaker.EventID,
		PersonID:    speaker.Person.ID,
		Company:     utils.SafeString(speaker.Company),
		Bio:         utils.SafeString(speaker.Bio),
		CreatedAt:   speaker.Audit.CreatedAt,
		UpdatedAt:   speaker.Audit.UpdatedAt,
		CreatedBy:   speaker.Audit.CreatedBy,
		UpdatedBy:   speaker.Audit.UpdatedBy,
		FirstName:   speaker.Person.FirstName,
		LastName:    speaker.Person.LastName,
		Email:       utils.SafeString(speaker.Person.Email),
		AvatarUrl:   utils.SafeString(speaker.Person.AvatarURL),
		GithubUser:  utils.SafeString(speaker.Person.GithubUser),
		TwitterUrl:  utils.SafeString(speaker.Person.TwitterURL),
		LinkedinUrl: utils.SafeString(speaker.Person.LinkedinURL),
		WebsiteUrl:  utils.SafeString(speaker.Person.WebsiteURL),
	}
}
