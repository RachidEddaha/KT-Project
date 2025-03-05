package films

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"task/internal/services/films/mocks"
	"task/pkg/database/entities/entitiescustom"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"task/internal/dto"
	"task/internal/models"
	"task/internal/models/kterrors"
	"task/pkg/customerror"
	"task/pkg/database/entities"
)

// setupMockRepository creates a mock repository for testing
func setupMockRepository(t *testing.T) *mocks.Repository {
	return mocks.NewRepository(t)
}

func TestGetFilmPaginated(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		mockBehavior   func(*mocks.Repository)
		inputRequest   dto.FilmSearchRequest
		expectedResult dto.FilmsPaginated
		expectedError  error
	}{
		{
			name: "Successful pagination with results",
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilmsPaginated", mock.Anything, 10, 0).Return(
					[]models.FilmPaginated{
						{ID: 1, Title: "Film 1", Qty: 2},
						{ID: 2, Title: "Film 2", Qty: 2},
					}, nil)
			},
			inputRequest: dto.FilmSearchRequest{
				Page:     1,
				PageSize: 10,
			},
			expectedResult: dto.FilmsPaginated{
				Films: []dto.Film{
					{ID: 1, Title: "Film 1"},
					{ID: 2, Title: "Film 2"},
				},
				Count:    2,
				Page:     1,
				PageSize: 10,
			},
		},
		{
			name: "Empty result set",
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilmsPaginated", mock.Anything, 10, 0).Return(
					[]models.FilmPaginated{}, nil)
			},
			inputRequest: dto.FilmSearchRequest{
				Page:     1,
				PageSize: 10,
			},
			expectedResult: dto.FilmsPaginated{
				Page:     1,
				PageSize: 10,
			},
		},
		{
			name: "Repository error",
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilmsPaginated", mock.Anything, 10, 0).Return(
					nil, errors.New("database error"))
			},
			inputRequest: dto.FilmSearchRequest{
				Page:     1,
				PageSize: 10,
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := setupMockRepository(t)
			tc.mockBehavior(mockRepo)

			// Create service with mock repository
			service := NewService(mockRepo)

			// Execute method
			result, err := service.GetFilmPaginated(context.Background(), tc.inputRequest)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetFilmDetail(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		filmID         int
		mockBehavior   func(*mocks.Repository)
		expectedResult dto.FilmDetail
		expectedError  error
	}{
		{
			name:   "Successful film detail retrieval",
			filmID: 1,
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilm", mock.Anything, 1).Return(
					entities.Film{
						ID:          1,
						Title:       "Test Film",
						Director:    "Test Director",
						ReleaseDate: entitiescustom.ReleaseDate{Time: time.Now()},
						Synopsis:    "Test Synopsis",
					}, nil)
			},
			expectedResult: dto.FilmDetail{
				ID:          1,
				Title:       "Test Film",
				Director:    "Test Director",
				ReleaseDate: entitiescustom.ReleaseDate{Time: time.Now()},
				Synopsis:    "Test Synopsis",
			},
		},
		{
			name:   "Film not found",
			filmID: 999,
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilm", mock.Anything, 999).Return(
					entities.Film{}, errors.New("film not found"))
			},
			expectedError: errors.New("film not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := setupMockRepository(t)
			tc.mockBehavior(mockRepo)

			service := NewService(mockRepo)

			result, err := service.GetFilmDetail(context.Background(), tc.filmID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult.ID, result.ID)
			assert.Equal(t, tc.expectedResult.Title, result.Title)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteFilm(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		filmID        int
		userID        int
		mockBehavior  func(*mocks.Repository)
		expectedError error
	}{
		{
			name:   "Successful film deletion",
			filmID: 1,
			userID: 100,
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilm", mock.Anything, 1).Return(
					entities.Film{
						ID:     1,
						UserID: 100,
					}, nil)
				mr.On("DeleteFilm", mock.Anything, 1).Return(nil)
			},
		},
		{
			name:   "Unauthorized film deletion",
			filmID: 1,
			userID: 200,
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilm", mock.Anything, 1).Return(
					entities.Film{
						ID:     1,
						UserID: 100,
					}, nil)
			},
			expectedError: customerror.NewCustomError(kterrors.UserCannotDeleteFilmError),
		},
		{
			name:   "Film not found",
			filmID: 999,
			userID: 100,
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilm", mock.Anything, 999).Return(
					entities.Film{}, errors.New("film not found"))
			},
			expectedError: errors.New("film not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := setupMockRepository(t)
			tc.mockBehavior(mockRepo)

			service := NewService(mockRepo)

			err := service.DeleteFilm(context.Background(), tc.filmID, tc.userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateFilm(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		request       dto.FilmCreateRequest
		mockBehavior  func(*mocks.Repository)
		expectedError error
	}{
		{
			name: "Successful film creation",
			request: dto.FilmCreateRequest{
				Title:       "New Film",
				Director:    "New Director",
				ReleaseDate: entitiescustom.ReleaseDate{Time: time.Now()},
				Synopsis:    "New Synopsis",
				UserID:      100,
			},
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("CreateFilm", mock.Anything, mock.AnythingOfType("entities.Film")).Return(nil)
			},
		},
		{
			name: "Duplicate film title",
			request: dto.FilmCreateRequest{
				Title:       "Existing Film",
				Director:    "Director",
				ReleaseDate: entitiescustom.ReleaseDate{Time: time.Now()},
				Synopsis:    "Synopsis",
				UserID:      100,
			},
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("CreateFilm", mock.Anything, mock.AnythingOfType("entities.Film")).Return(
					gorm.ErrDuplicatedKey)
			},
			expectedError: customerror.NewCustomError(kterrors.FilmTitleAlreadyExistsError),
		},
		{
			name: "Other database error",
			request: dto.FilmCreateRequest{
				Title:       "Another Film",
				Director:    "Another Director",
				ReleaseDate: entitiescustom.ReleaseDate{Time: time.Now()},
				Synopsis:    "Another Synopsis",
				UserID:      100,
			},
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("CreateFilm", mock.Anything, mock.AnythingOfType("entities.Film")).Return(
					errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := setupMockRepository(t)
			tc.mockBehavior(mockRepo)

			service := NewService(mockRepo)

			err := service.CreateFilm(context.Background(), tc.request)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateFilm(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		request       dto.FilmUpdateRequest
		mockBehavior  func(*mocks.Repository)
		expectedError error
	}{
		{
			name: "Successful film update",
			request: dto.FilmUpdateRequest{
				ID:          1,
				Title:       "Updated Film",
				Director:    "Updated Director",
				ReleaseDate: entitiescustom.ReleaseDate{Time: time.Now()},
				Synopsis:    "Updated Synopsis",
				UserID:      100,
			},
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilm", mock.Anything, 1).Return(
					entities.Film{
						ID:     1,
						UserID: 100,
					}, nil)
				mr.On("UpdateFilm", mock.Anything, mock.AnythingOfType("entities.Film")).Return(nil)
			},
		},
		{
			name: "Unauthorized film update",
			request: dto.FilmUpdateRequest{
				ID:          1,
				Title:       "Updated Film",
				Director:    "Updated Director",
				ReleaseDate: entitiescustom.ReleaseDate{Time: time.Now()},
				Synopsis:    "Updated Synopsis",
				UserID:      200,
			},
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilm", mock.Anything, 1).Return(
					entities.Film{
						ID:     1,
						UserID: 100,
					}, nil)
			},
			expectedError: customerror.NewCustomError(kterrors.UserCannotUpdateFilmError),
		},
		{
			name: "Duplicate film title",
			request: dto.FilmUpdateRequest{
				ID:          1,
				Title:       "Existing Film",
				Director:    "Updated Director",
				ReleaseDate: entitiescustom.ReleaseDate{Time: time.Now()},
				Synopsis:    "Updated Synopsis",
				UserID:      100,
			},
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilm", mock.Anything, 1).Return(
					entities.Film{
						ID:     1,
						UserID: 100,
					}, nil)
				mr.On("UpdateFilm", mock.Anything, mock.AnythingOfType("entities.Film")).Return(
					gorm.ErrDuplicatedKey)
			},
			expectedError: customerror.NewCustomError(kterrors.FilmTitleAlreadyExistsError),
		},
		{
			name: "Film not found",
			request: dto.FilmUpdateRequest{
				ID:          999,
				Title:       "Updated Film",
				Director:    "Updated Director",
				ReleaseDate: entitiescustom.ReleaseDate{Time: time.Now()},
				Synopsis:    "Updated Synopsis",
				UserID:      100,
			},
			mockBehavior: func(mr *mocks.Repository) {
				mr.On("GetFilm", mock.Anything, 999).Return(
					entities.Film{}, errors.New("film not found"))
			},
			expectedError: errors.New("film not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := setupMockRepository(t)
			tc.mockBehavior(mockRepo)

			service := NewService(mockRepo)

			err := service.UpdateFilm(context.Background(), tc.request)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
