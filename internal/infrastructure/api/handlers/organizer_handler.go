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

type OrganizerHandler struct {
	usecase domain.OrganizerUsecase
}

func NewOrganizerHandler(organizerUsecase domain.OrganizerUsecase) *OrganizerHandler {
	return &OrganizerHandler{usecase: organizerUsecase}
}

// GetAll
func (o *OrganizerHandler) GetAll(ctx *gin.Context) {
	eventID, err := uuid.Parse(ctx.Param("event-id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Event ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	organizers, err := o.usecase.GetAll(ctx.Request.Context(), eventID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	orgResponse := make([]dtos.OrganizerDetailResponse, len(organizers))
	for i, organizer := range organizers {
		orgResponse[i] = mapToOrganizerDetailResponse(&organizer)
	}

	ctx.JSON(http.StatusOK, orgResponse)
}

// GetById
func (o *OrganizerHandler) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	organizer, err := o.usecase.GetById(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	orgResponse := mapToOrganizerDetailResponse(organizer)

	ctx.JSON(http.StatusOK, orgResponse)
}

// ListPaged
func (o *OrganizerHandler) ListPaged(ctx *gin.Context) {
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

	organizers, total, err := o.usecase.ListPaged(ctx.Request.Context(), eventID, search, int32(page), int32(pageSize))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	orgResponse := make([]dtos.OrganizerDetailResponse, len(organizers))
	for i, organizer := range organizers {
		orgResponse[i] = mapToOrganizerDetailResponse(&organizer)
	}

	res := dtos.PagedResponse[dtos.OrganizerDetailResponse]{
		Data: orgResponse,
		Meta: dtos.PagedMeta{
			Total:    total,
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
	}

	ctx.JSON(http.StatusOK, res)
}

// Create
func (o *OrganizerHandler) Create(ctx *gin.Context) {
	var dto dtos.CreateOrganizerDTO
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

	org := &domain.Organizer{
		EventID:         dto.EventID,
		Company:         dto.Company,
		RoleDescription: dto.RoleDescription,
		Person:          domain.Person{ID: dto.PersonID},
		Audit:           domain.Audit{CreatedBy: uid},
	}

	organizer, err := o.usecase.Create(ctx.Request.Context(), org)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	orgResponse := mapToOrganizerResponse(organizer)

	ctx.JSON(http.StatusCreated, orgResponse)
}

// Update
func (o *OrganizerHandler) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	var dto dtos.UpdateOrganizerDTO
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

	updOrganizer := domain.UpdateOrganizer{
		Company:         dto.Company,
		RoleDescription: dto.RoleDescription,
		UpdatedBy:       uid,
	}

	organizer, err := o.usecase.Update(ctx.Request.Context(), id, &updOrganizer)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	orgResponse := mapToOrganizerResponse(organizer)

	ctx.JSON(http.StatusOK, orgResponse)
}

// Delete
func (o *OrganizerHandler) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	err = o.usecase.Delete(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func mapToOrganizerResponse(organizer *domain.Organizer) dtos.OrganizerResponse {
	return dtos.OrganizerResponse{
		ID:              organizer.ID,
		EventID:         organizer.EventID,
		PersonID:        organizer.Person.ID,
		Company:         utils.SafeString(organizer.Company),
		RoleDescription: utils.SafeString(organizer.RoleDescription),
		CreatedAt:       organizer.CreatedAt,
		UpdatedAt:       organizer.UpdatedAt,
		CreatedBy:       organizer.CreatedBy,
		UpdatedBy:       organizer.UpdatedBy,
	}
}

func mapToOrganizerDetailResponse(organizer *domain.Organizer) dtos.OrganizerDetailResponse {
	return dtos.OrganizerDetailResponse{
		ID:              organizer.ID,
		EventID:         organizer.EventID,
		PersonID:        organizer.Person.ID,
		Company:         utils.SafeString(organizer.Company),
		RoleDescription: utils.SafeString(organizer.RoleDescription),
		CreatedAt:       organizer.CreatedAt,
		UpdatedAt:       organizer.UpdatedAt,
		CreatedBy:       organizer.CreatedBy,
		UpdatedBy:       organizer.UpdatedBy,
		FirstName:       organizer.Person.FirstName,
		LastName:        organizer.Person.LastName,
		Email:           utils.SafeString(organizer.Person.Email),
		AvatarUrl:       utils.SafeString(organizer.Person.AvatarURL),
		GithubUser:      utils.SafeString(organizer.Person.GithubUser),
		TwitterUrl:      utils.SafeString(organizer.Person.TwitterURL),
		LinkedinUrl:     utils.SafeString(organizer.Person.LinkedinURL),
		WebsiteUrl:      utils.SafeString(organizer.Person.WebsiteURL),
	}
}
