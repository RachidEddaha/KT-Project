// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"
	entities "task/pkg/database/entities"

	mock "github.com/stretchr/testify/mock"

	models "task/internal/models"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// CreateFilm provides a mock function with given fields: ctx, film
func (_m *Repository) CreateFilm(ctx context.Context, film entities.Film) error {
	ret := _m.Called(ctx, film)

	if len(ret) == 0 {
		panic("no return value specified for CreateFilm")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entities.Film) error); ok {
		r0 = rf(ctx, film)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteFilm provides a mock function with given fields: ctx, id
func (_m *Repository) DeleteFilm(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteFilm")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetFilm provides a mock function with given fields: ctx, ID
func (_m *Repository) GetFilm(ctx context.Context, ID int) (entities.Film, error) {
	ret := _m.Called(ctx, ID)

	if len(ret) == 0 {
		panic("no return value specified for GetFilm")
	}

	var r0 entities.Film
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (entities.Film, error)); ok {
		return rf(ctx, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) entities.Film); ok {
		r0 = rf(ctx, ID)
	} else {
		r0 = ret.Get(0).(entities.Film)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFilmsPaginated provides a mock function with given fields: ctx, pageSize, offset
func (_m *Repository) GetFilmsPaginated(ctx context.Context, pageSize int, offset int) ([]models.FilmPaginated, error) {
	ret := _m.Called(ctx, pageSize, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetFilmsPaginated")
	}

	var r0 []models.FilmPaginated
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) ([]models.FilmPaginated, error)); ok {
		return rf(ctx, pageSize, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int) []models.FilmPaginated); ok {
		r0 = rf(ctx, pageSize, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.FilmPaginated)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, pageSize, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateFilm provides a mock function with given fields: ctx, film
func (_m *Repository) UpdateFilm(ctx context.Context, film entities.Film) error {
	ret := _m.Called(ctx, film)

	if len(ret) == 0 {
		panic("no return value specified for UpdateFilm")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entities.Film) error); ok {
		r0 = rf(ctx, film)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
