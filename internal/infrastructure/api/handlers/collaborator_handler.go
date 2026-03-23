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

type CollaboratorHandler struct {
	usecase domain.CollaboratorUsecase
}

func NewCollaboratorHandler(collaboratorUsecase domain.CollaboratorUsecase) *CollaboratorHandler {
	return &CollaboratorHandler{usecase: collaboratorUsecase}
}

// GetAll
func (ch *CollaboratorHandler) GetAll(ctx *gin.Context) {
	eventID, err := uuid.Parse(ctx.Param("event-id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Event ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	collaborators, err := ch.usecase.GetAll(ctx.Request.Context(), eventID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	collResponse := make([]dtos.CollaboratorDetailResponse, len(collaborators))
	for i, collaborator := range collaborators {
		collResponse[i] = mapToCollaboratorDetailResponse(&collaborator)
	}

	ctx.JSON(http.StatusOK, collResponse)
}

// GetByID
func (ch *CollaboratorHandler) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	collaborator, err := ch.usecase.GetById(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	colResponse := mapToCollaboratorResponse(collaborator)

	ctx.JSON(http.StatusOK, colResponse)
}

// ListPaged
func (ch *CollaboratorHandler) ListPaged(ctx *gin.Context) {
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

	collaborators, total, err := ch.usecase.ListPaged(ctx.Request.Context(), eventID, search, int32(page), int32(pageSize))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	collaboratorsDetail := make([]dtos.CollaboratorDetailResponse, len(collaborators))
	for i, collaborator := range collaborators {
		collaboratorsDetail[i] = mapToCollaboratorDetailResponse(&collaborator)
	}

	res := dtos.PagedResponse[dtos.CollaboratorDetailResponse]{
		Data: collaboratorsDetail,
		Meta: dtos.PagedMeta{
			Total:    total,
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
	}

	ctx.JSON(http.StatusOK, res)
}

// Create
func (ch *CollaboratorHandler) Create(ctx *gin.Context) {
	var dto dtos.CreateCollaboratorDTO
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

	coll := &domain.Collaborator{
		EventID: dto.EventID,
		Area:    dto.Area,
		Person:  domain.Person{ID: dto.PersonID},
		Audit:   domain.Audit{CreatedBy: uid},
	}

	collaborator, err := ch.usecase.Create(ctx.Request.Context(), coll)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	collResponse := mapToCollaboratorResponse(collaborator)

	ctx.JSON(http.StatusCreated, collResponse)
}

// Update
func (ch *CollaboratorHandler) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	var dto dtos.UpdateCollaboratorDTO
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

	updCollaborator := domain.UpdateCollaborator{
		Area:      dto.Area,
		UpdatedBy: uid,
	}

	collaborator, err := ch.usecase.Update(ctx.Request.Context(), id, &updCollaborator)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	collResponse := mapToCollaboratorResponse(collaborator)

	ctx.JSON(http.StatusOK, collResponse)
}

// Delete
func (ch *CollaboratorHandler) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	err = ch.usecase.Delete(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func mapToCollaboratorResponse(c *domain.Collaborator) dtos.CollaboratorResponse {
	return dtos.CollaboratorResponse{
		ID:        c.ID,
		PersonID:  c.Person.ID,
		EventID:   c.EventID,
		Area:      utils.SafeString(c.Area),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		CreatedBy: c.CreatedBy,
		UpdatedBy: c.UpdatedBy,
	}
}

func mapToCollaboratorDetailResponse(c *domain.Collaborator) dtos.CollaboratorDetailResponse {
	return dtos.CollaboratorDetailResponse{
		ID:          c.ID,
		PersonID:    c.Person.ID,
		EventID:     c.EventID,
		Area:        utils.SafeString(c.Area),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		CreatedBy:   c.CreatedBy,
		UpdatedBy:   c.UpdatedBy,
		FirstName:   c.Person.FirstName,
		LastName:    c.Person.LastName,
		Email:       utils.SafeString(c.Person.Email),
		AvatarUrl:   utils.SafeString(c.Person.AvatarURL),
		GithubUser:  utils.SafeString(c.Person.GithubUser),
		TwitterUrl:  utils.SafeString(c.Person.TwitterURL),
		LinkedinUrl: utils.SafeString(c.Person.LinkedinURL),
		WebsiteUrl:  utils.SafeString(c.Person.WebsiteURL),
	}
}
