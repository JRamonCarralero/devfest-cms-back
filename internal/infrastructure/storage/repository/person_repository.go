package repository

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/infrastructure/storage/dbgen"

	"github.com/google/uuid"
)

type PersonRepository struct {
	queries *dbgen.Queries
}

// NewPersonRepository returns a new PersonRepository
func NewPersonRepository(queries *dbgen.Queries) *PersonRepository {
	return &PersonRepository{queries: queries}
}

// --- READERS ---

// GetAll returns all Persons
func (r *PersonRepository) GetAll(ctx context.Context) ([]domain.Person, error) {
	rows, err := r.queries.ListPersons(ctx)
	if err != nil {
		return nil, ParseDBError(err, "Person")
	}

	persons := make([]domain.Person, len(rows))
	for i, row := range rows {
		persons[i] = *mapToDomainPerson(row)
	}

	return persons, nil
}

// GetById returns a Person by its ID
func (r *PersonRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Person, error) {
	row, err := r.queries.GetPersonByID(ctx, id)
	if err != nil {
		return nil, ParseDBError(err, "Person")
	}

	return mapToDomainPerson(row), nil
}

// GetByEmail returns a Person by its email
func (r *PersonRepository) GetByEmail(ctx context.Context, email *string) (*domain.Person, error) {
	row, err := r.queries.GetPersonByEmail(ctx, PtrToText(email))
	if err != nil {
		return nil, ParseDBError(err, "Person")
	}

	return mapToDomainPerson(row), nil
}

// ListPaged returns a page of Persons
func (r *PersonRepository) ListPaged(ctx context.Context, search string, page, pageSize int32) ([]domain.Person, int64, error) {
	offset := (page - 1) * pageSize

	total, err := r.queries.CountPersons(ctx, search)
	if err != nil {
		return nil, 0, ParseDBError(err, "Person")
	}

	rows, err := r.queries.ListPersonsPaged(ctx, dbgen.ListPersonsPagedParams{
		Search: search,
		Limit:  pageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, ParseDBError(err, "Person")
	}

	persons := make([]domain.Person, len(rows))
	for i, row := range rows {
		persons[i] = *mapToDomainPerson(row)
	}

	return persons, total, nil
}

// --- WRITERS ---

// Create creates a new Person
func (r *PersonRepository) Create(ctx context.Context, person *domain.Person) (*domain.Person, error) {
	row, err := r.queries.CreatePerson(ctx, dbgen.CreatePersonParams{
		FirstName:   person.FirstName,
		LastName:    person.LastName,
		Email:       PtrToText(person.Email),
		AvatarUrl:   PtrToText(person.AvatarURL),
		GithubUser:  PtrToText(person.GithubUser),
		LinkedinUrl: PtrToText(person.LinkedinURL),
		TwitterUrl:  PtrToText(person.TwitterURL),
		WebsiteUrl:  PtrToText(person.WebsiteURL),
		CreatedBy:   person.Audit.CreatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Person")
	}

	return mapToDomainPerson(row), nil
}

// Update updates a Person
func (r *PersonRepository) Update(ctx context.Context, person *domain.Person) (*domain.Person, error) {
	row, err := r.queries.UpdatePerson(ctx, dbgen.UpdatePersonParams{
		ID:          person.ID,
		FirstName:   person.FirstName,
		LastName:    person.LastName,
		Email:       PtrToText(person.Email),
		AvatarUrl:   PtrToText(person.AvatarURL),
		GithubUser:  PtrToText(person.GithubUser),
		LinkedinUrl: PtrToText(person.LinkedinURL),
		TwitterUrl:  PtrToText(person.TwitterURL),
		WebsiteUrl:  PtrToText(person.WebsiteURL),
		UpdatedBy:   person.Audit.UpdatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Person")
	}

	return mapToDomainPerson(row), nil
}

// Delete deletes a Person
func (r *PersonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeletePerson(ctx, id)
	if err != nil {
		return ParseDBError(err, "Person")
	}

	return nil
}

// --- Mappers ---

// mapToDomain maps a dbgen.Person to a domain.Person
func mapToDomainPerson(dbPerson dbgen.Person) *domain.Person {
	return &domain.Person{
		ID:          dbPerson.ID,
		FirstName:   dbPerson.FirstName,
		LastName:    dbPerson.LastName,
		Email:       TextToPtr(dbPerson.Email),
		AvatarURL:   TextToPtr(dbPerson.AvatarUrl),
		GithubUser:  TextToPtr(dbPerson.GithubUser),
		LinkedinURL: TextToPtr(dbPerson.LinkedinUrl),
		TwitterURL:  TextToPtr(dbPerson.TwitterUrl),
		WebsiteURL:  TextToPtr(dbPerson.WebsiteUrl),

		Audit: domain.Audit{
			CreatedAt: dbPerson.CreatedAt.Time,
			UpdatedAt: dbPerson.UpdatedAt.Time,
			CreatedBy: dbPerson.CreatedBy,
			UpdatedBy: dbPerson.UpdatedBy,
		},
	}
}
