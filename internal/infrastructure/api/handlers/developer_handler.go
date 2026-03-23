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

type DeveloperHandler struct {
	usecase domain.DeveloperUsecase
}

func NewDeveloperHandler(developerUsecase domain.DeveloperUsecase) *DeveloperHandler {
	return &DeveloperHandler{usecase: developerUsecase}
}

// GetAll
func (h *DeveloperHandler) GetAll(ctx *gin.Context) {
	eventID, err := uuid.Parse(ctx.Param("event-id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Event ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	developers, err := h.usecase.GetAll(ctx.Request.Context(), eventID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	devResponse := make([]dtos.DeveloperDetailResponse, len(developers))
	for i, developer := range developers {
		devResponse[i] = mapToDeveloperDetailResponse(&developer)
	}

	ctx.JSON(http.StatusOK, devResponse)
}

// GetByID
func (h *DeveloperHandler) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	developer, err := h.usecase.GetById(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	devResponse := mapToDeveloperDetailResponse(developer)

	ctx.JSON(http.StatusOK, devResponse)
}

// ListPaged
func (h *DeveloperHandler) ListPaged(ctx *gin.Context) {
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

	developers, total, err := h.usecase.ListPaged(ctx.Request.Context(), eventID, search, int32(page), int32(pageSize))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	devResponse := make([]dtos.DeveloperDetailResponse, len(developers))
	for i, developer := range developers {
		devResponse[i] = mapToDeveloperDetailResponse(&developer)
	}

	res := dtos.PagedResponse[dtos.DeveloperDetailResponse]{
		Data: devResponse,
		Meta: dtos.PagedMeta{
			Total:    total,
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
	}

	ctx.JSON(http.StatusOK, res)
}

// Create
func (h *DeveloperHandler) Create(ctx *gin.Context) {
	var dto dtos.CreateDeveloperDTO
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

	dev := &domain.Developer{
		EventID:         dto.EventID,
		RoleDescription: dto.RoleDescription,
		Person:          domain.Person{ID: dto.PersonID},
		Audit:           domain.Audit{CreatedBy: uid},
	}

	developer, err := h.usecase.Create(ctx.Request.Context(), dev)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	devResponse := mapToDeveloperResponse(developer)

	ctx.JSON(http.StatusCreated, devResponse)
}

// Update
func (h *DeveloperHandler) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	var dto dtos.UpdateDeveloperDTO
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

	updDeveloper := domain.UpdateDeveloper{
		RoleDescription: dto.RoleDescription,
		UpdatedBy:       uid,
	}

	developer, err := h.usecase.Update(ctx.Request.Context(), id, &updDeveloper)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	devResponse := mapToDeveloperResponse(developer)

	ctx.JSON(http.StatusOK, devResponse)
}

// Delete
func (h *DeveloperHandler) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	err = h.usecase.Delete(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func mapToDeveloperResponse(d *domain.Developer) dtos.DeveloperResponse {
	return dtos.DeveloperResponse{
		ID:              d.ID,
		EventID:         d.EventID,
		PersonID:        d.Person.ID,
		RoleDescription: utils.SafeString(d.RoleDescription),
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
		CreatedBy:       d.CreatedBy,
		UpdatedBy:       d.UpdatedBy,
	}
}

func mapToDeveloperDetailResponse(d *domain.Developer) dtos.DeveloperDetailResponse {
	return dtos.DeveloperDetailResponse{
		ID:              d.ID,
		EventID:         d.EventID,
		RoleDescription: utils.SafeString(d.RoleDescription),
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
		CreatedBy:       d.CreatedBy,
		UpdatedBy:       d.UpdatedBy,
		FirstName:       d.Person.FirstName,
		LastName:        d.Person.LastName,
		Email:           utils.SafeString(d.Person.Email),
		AvatarUrl:       utils.SafeString(d.Person.AvatarURL),
		GithubUser:      utils.SafeString(d.Person.GithubUser),
		TwitterUrl:      utils.SafeString(d.Person.TwitterURL),
		LinkedinUrl:     utils.SafeString(d.Person.LinkedinURL),
		WebsiteUrl:      utils.SafeString(d.Person.WebsiteURL),
	}
}
