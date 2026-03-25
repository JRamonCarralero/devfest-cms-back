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

type PersonHandler struct {
	usecase domain.PersonUsecase
}

func NewPersonHandler(usecase domain.PersonUsecase) *PersonHandler {
	return &PersonHandler{usecase: usecase}
}

// GetAll
func (h *PersonHandler) GetAll(ctx *gin.Context) {
	people, err := h.usecase.GetAll(ctx.Request.Context())
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	if len(people) == 0 {
		ctx.Status(http.StatusNoContent)
		return
	}

	peopleResponse := make([]dtos.PersonResponse, len(people))
	for i, person := range people {
		peopleResponse[i] = mapToDTOPersonResponse(person)
	}

	ctx.JSON(http.StatusOK, peopleResponse)
}

// GetByID
func (h *PersonHandler) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	person, err := h.usecase.GetByID(ctx.Request.Context(), id)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	personResponse := mapToDTOPersonResponse(*person)

	ctx.JSON(http.StatusOK, personResponse)
}

// GetByEmail
func (h *PersonHandler) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	person, err := h.usecase.GetByEmail(ctx.Request.Context(), &email)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	personResponse := mapToDTOPersonResponse(*person)

	ctx.JSON(http.StatusOK, personResponse)
}

// GetPaged
func (h *PersonHandler) GetPaged(ctx *gin.Context) {
	search := ctx.DefaultQuery("search", "")
	page, pageSize := utils.GetPaginationParams(ctx)
	if pageSize > 100 {
		pageSize = 100
	}

	people, total, err := h.usecase.ListPaged(ctx.Request.Context(), search, int32(page), int32(pageSize))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	peopleResponse := make([]dtos.PersonResponse, len(people))
	for i, person := range people {
		peopleResponse[i] = mapToDTOPersonResponse(person)
	}

	res := dtos.PagedResponse[dtos.PersonResponse]{
		Data: peopleResponse,
		Meta: dtos.PagedMeta{
			Total:    total,
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
	}

	ctx.JSON(http.StatusOK, res)
}

// Create
func (h *PersonHandler) Create(ctx *gin.Context) {
	var dto dtos.CreatePersonDTO
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

	person := domain.Person{
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Email:       dto.Email,
		AvatarURL:   dto.AvatarURL,
		GithubUser:  dto.GithubUser,
		LinkedinURL: dto.LinkedinURL,
		TwitterURL:  dto.TwitterURL,
		WebsiteURL:  dto.WebsiteURL,
		Audit: domain.Audit{
			CreatedBy: uid,
			UpdatedBy: uid,
		},
	}

	newPerson, err := h.usecase.Create(ctx.Request.Context(), &person)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	personResponse := mapToDTOPersonResponse(*newPerson)

	ctx.JSON(http.StatusCreated, personResponse)
}

// Update
func (h *PersonHandler) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(ctx, newErr)
		return
	}

	var dto dtos.UpdatePersonDTO
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

	upPerson := domain.UpdatePerson{
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Email:       dto.Email,
		AvatarURL:   dto.AvatarURL,
		GithubUser:  dto.GithubUser,
		LinkedinURL: dto.LinkedinURL,
		TwitterURL:  dto.TwitterURL,
		WebsiteURL:  dto.WebsiteURL,
		UpdatedBy:   uid,
	}

	person, err := h.usecase.Update(ctx.Request.Context(), id, &upPerson)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	personResponse := mapToDTOPersonResponse(*person)

	ctx.JSON(http.StatusOK, personResponse)
}

// Delete
func (h *PersonHandler) Delete(ctx *gin.Context) {
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

func mapToDTOPersonResponse(person domain.Person) dtos.PersonResponse {
	return dtos.PersonResponse{
		ID: person.ID,
		PersonFieldsDTO: dtos.PersonFieldsDTO{
			FirstName:   person.FirstName,
			LastName:    person.LastName,
			Email:       utils.SafeString(person.Email),
			AvatarURL:   utils.SafeString(person.AvatarURL),
			GithubUser:  utils.SafeString(person.GithubUser),
			LinkedinURL: utils.SafeString(person.LinkedinURL),
			TwitterURL:  utils.SafeString(person.TwitterURL),
			WebsiteURL:  utils.SafeString(person.WebsiteURL),
		},
		AuditDTO: dtos.AuditDTO{
			CreatedAt: person.Audit.CreatedAt,
			CreatedBy: person.Audit.CreatedBy,
			UpdatedAt: person.Audit.UpdatedAt,
			UpdatedBy: person.Audit.UpdatedBy,
		},
	}
}
