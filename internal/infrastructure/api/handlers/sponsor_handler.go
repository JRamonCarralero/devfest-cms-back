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

type SponsorHandler struct {
	usecase domain.SponsorUsecase
}

func NewSponsorHandler(usecase domain.SponsorUsecase) *SponsorHandler {
	return &SponsorHandler{usecase: usecase}
}

// GetAll
func (h *SponsorHandler) GetAll(ctx *gin.Context) {
	eventID, err := uuid.Parse(ctx.Param("event-id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid Event ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	sponsors, err := h.usecase.GetAll(ctx.Request.Context(), eventID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	sponsorResponse := make([]dtos.SponsorResponse, len(sponsors))
	for i, sponsor := range sponsors {
		sponsorResponse[i] = mapToSponsorResponse(&sponsor)
	}

	ctx.JSON(http.StatusOK, sponsorResponse)
}

// GetById
func (h *SponsorHandler) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	sponsor, err := h.usecase.GetById(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	sponsorResponse := mapToSponsorResponse(sponsor)

	ctx.JSON(http.StatusOK, sponsorResponse)
}

// ListPaged
func (h *SponsorHandler) ListPaged(ctx *gin.Context) {
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

	sponsors, total, err := h.usecase.ListPaged(ctx.Request.Context(), eventID, search, int32(page), int32(pageSize))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	sponsorResponse := make([]dtos.SponsorResponse, len(sponsors))
	for i, sponsor := range sponsors {
		sponsorResponse[i] = mapToSponsorResponse(&sponsor)
	}

	res := dtos.PagedResponse[dtos.SponsorResponse]{
		Data: sponsorResponse,
		Meta: dtos.PagedMeta{
			Total:    total,
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
	}

	ctx.JSON(http.StatusOK, res)
}

// Create
func (h *SponsorHandler) Create(ctx *gin.Context) {
	var req dtos.CreateSponsorDTO
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

	sponsor := &domain.Sponsor{
		EventID:       req.EventID,
		Name:          req.Name,
		LogoURL:       req.LogoURL,
		WebsiteURL:    req.WebsiteURL,
		Tier:          req.Tier,
		OrderPriority: req.OrderPriority,
		Audit:         domain.Audit{CreatedBy: uid},
	}

	sponsor, err = h.usecase.Create(ctx.Request.Context(), sponsor)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	sponsorResponse := mapToSponsorResponse(sponsor)

	ctx.JSON(http.StatusCreated, sponsorResponse)
}

// Update
func (h *SponsorHandler) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	var req dtos.UpdateSponsorDTO
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

	updSponsor := &domain.UpdateSponsor{
		Name:          req.Name,
		LogoURL:       req.LogoURL,
		WebsiteURL:    req.WebsiteURL,
		Tier:          req.Tier,
		OrderPriority: req.OrderPriority,
		UpdatedBy:     uid,
	}

	sponsor, err := h.usecase.Update(ctx.Request.Context(), id, updSponsor)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	sponsorResponse := mapToSponsorResponse(sponsor)

	ctx.JSON(http.StatusOK, sponsorResponse)
}

// Delete
func (h *SponsorHandler) Delete(ctx *gin.Context) {
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

// mapToSponsorResponse
func mapToSponsorResponse(sponsor *domain.Sponsor) dtos.SponsorResponse {
	return dtos.SponsorResponse{
		ID:            sponsor.ID,
		EventID:       sponsor.EventID,
		Name:          sponsor.Name,
		LogoURL:       sponsor.LogoURL,
		WebsiteURL:    sponsor.WebsiteURL,
		Tier:          sponsor.Tier,
		OrderPriority: utils.SafeInt(sponsor.OrderPriority),
		AuditDTO: dtos.AuditDTO{
			CreatedAt: sponsor.Audit.CreatedAt,
			UpdatedAt: sponsor.Audit.UpdatedAt,
			CreatedBy: sponsor.Audit.CreatedBy,
			UpdatedBy: sponsor.Audit.UpdatedBy,
		},
	}
}
