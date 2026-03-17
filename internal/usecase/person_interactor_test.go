package usecase_test

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/domain/mocks"
	"devfest/internal/infrastructure/api/dtos"
	"devfest/internal/usecase"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPersonInteractor_Create(t *testing.T) {
	mockRepo := new(mocks.PersonRepository)
	interactor := usecase.NewPersonInteractor(mockRepo)
	ctx := context.Background()

	t.Run("Should create person successfully", func(t *testing.T) {
		email := "test@test.com"
		dto := dtos.CreatePersonDTO{
			FirstName: "John",
			LastName:  "Doe",
			Email:     &email,
			CreatedBy: uuid.New(),
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Person")).
			Return(&domain.Person{
				ID:        uuid.New(),
				FirstName: dto.FirstName,
				LastName:  dto.LastName,
				Email:     dto.Email,
			}, nil).Once()

		result, err := interactor.Create(ctx, dto)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, dto.FirstName, result.FirstName)
		mockRepo.AssertExpectations(t)
	})
}

func TestPersonInteractor_GetByEmail(t *testing.T) {
	mockRepo := new(mocks.PersonRepository)
	interactor := usecase.NewPersonInteractor(mockRepo)
	ctx := context.Background()

	t.Run("Should return error if email is nil", func(t *testing.T) {
		result, err := interactor.GetByEmail(ctx, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "GetByEmail", mock.Anything, mock.Anything)
	})
}

func TestPersonInteractor_Update(t *testing.T) {
	mockRepo := new(mocks.PersonRepository)
	interactor := usecase.NewPersonInteractor(mockRepo)
	ctx := context.Background()
	personID := uuid.New()

	t.Run("Should update only provided fields", func(t *testing.T) {
		existing := &domain.Person{
			ID:        personID,
			FirstName: "OldName",
			LastName:  "OldLastName",
		}

		newName := "NewName"
		dto := dtos.UpdatePersonDTO{
			FirstName: &newName,
			UpdatedBy: uuid.New(),
		}

		mockRepo.On("GetById", ctx, personID).Return(existing, nil).Once()
		mockRepo.On("Update", ctx, mock.MatchedBy(func(p *domain.Person) bool {
			return p.FirstName == newName && p.LastName == "OldLastName"
		})).Return(existing, nil).Once()

		_, err := interactor.Update(ctx, personID, dto)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPersonInteractor_Read(t *testing.T) {
	mockRepo := new(mocks.PersonRepository)
	interactor := usecase.NewPersonInteractor(mockRepo)
	ctx := context.Background()

	t.Run("GetAll - Should return list of persons", func(t *testing.T) {
		expected := []domain.Person{{ID: uuid.New(), FirstName: "Test"}}
		mockRepo.On("GetAll", ctx).Return(expected, nil).Once()

		result, err := interactor.GetAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByID - Should return person if exists", func(t *testing.T) {
		id := uuid.New()
		expected := &domain.Person{ID: id, FirstName: "Test"}
		mockRepo.On("GetById", ctx, id).Return(expected, nil).Once()

		result, err := interactor.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, id, result.ID)
	})
}

func TestPersonInteractor_ListPaged(t *testing.T) {
	mockRepo := new(mocks.PersonRepository)
	interactor := usecase.NewPersonInteractor(mockRepo)
	ctx := context.Background()

	t.Run("Should apply default pagination values", func(t *testing.T) {
		mockRepo.On("ListPaged", ctx, "", int32(1), int32(10)).
			Return([]domain.Person{}, int64(0), nil).Once()

		_, _, err := interactor.ListPaged(ctx, "", 0, 0)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPersonInteractor_Delete(t *testing.T) {
	mockRepo := new(mocks.PersonRepository)
	interactor := usecase.NewPersonInteractor(mockRepo)
	ctx := context.Background()
	id := uuid.New()

	t.Run("Should call delete on repository", func(t *testing.T) {
		mockRepo.On("Delete", ctx, id).Return(nil).Once()

		err := interactor.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPersonInteractor_Update_Errors(t *testing.T) {
	mockRepo := new(mocks.PersonRepository)
	interactor := usecase.NewPersonInteractor(mockRepo)
	ctx := context.Background()
	id := uuid.New()

	t.Run("Should return error if person not found", func(t *testing.T) {
		dto := dtos.UpdatePersonDTO{}
		mockRepo.On("GetById", ctx, id).
			Return(nil, domain.NewAppError(domain.TypeNotFound, "not found", nil)).Once()

		result, err := interactor.Update(ctx, id, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return error if repository update fails", func(t *testing.T) {
		existing := &domain.Person{ID: id}
		dto := dtos.UpdatePersonDTO{}

		mockRepo.On("GetById", ctx, id).Return(existing, nil).Once()
		mockRepo.On("Update", ctx, mock.Anything).
			Return(nil, domain.NewAppError(domain.TypeInternal, "db error", nil)).Once()

		_, err := interactor.Update(ctx, id, dto)

		assert.Error(t, err)
	})
}
