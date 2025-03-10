package authentication_test

import (
	"KTOnlinePlatform/internal/services/authentication"
	"KTOnlinePlatform/internal/services/authentication/mocks"
	"KTOnlinePlatform/pkg/logger"
	"context"
	"errors"
	"gorm.io/gorm"
	"testing"

	"KTOnlinePlatform/internal/dto"
	"KTOnlinePlatform/internal/models/kterrors"
	"KTOnlinePlatform/pkg/customerror"
	"KTOnlinePlatform/pkg/database/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	logger.InitializeForTest()
	testCases := []struct {
		name          string
		username      string
		password      string
		setupMocks    func(*mocks.Repository, *mocks.TokensGeneration)
		expectedError string
		expectedToken dto.JWTTokens
	}{
		{
			name:     "Success",
			username: "validuser",
			password: "ValidPassword123!",
			setupMocks: func(repo *mocks.Repository, tokenGen *mocks.TokensGeneration) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("ValidPassword123!"), bcrypt.DefaultCost)
				user := entities.User{
					ID:       1,
					Username: "validuser",
					Password: string(hashedPassword),
				}
				tokens := dto.JWTTokens{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				}
				repo.On("FindUser", mock.Anything, "validuser").Return(user, nil)
				tokenGen.On("GenerateAuthTokens", 1, "validuser").Return(tokens, nil)
			},
			expectedError: "",
			expectedToken: dto.JWTTokens{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			},
		},
		{
			name:     "User not found",
			username: "nonexistentuser",
			password: "Password123!",
			setupMocks: func(repo *mocks.Repository, tokenGen *mocks.TokensGeneration) {
				repo.On("FindUser", mock.Anything, "nonexistentuser").Return(entities.User{}, gorm.ErrRecordNotFound)
			},
			expectedError: kterrors.UserNotFoundError,
			expectedToken: dto.JWTTokens{},
		},
		{
			name:     "Wrong password",
			username: "validuser",
			password: "WrongPassword123!",
			setupMocks: func(repo *mocks.Repository, tokenGen *mocks.TokensGeneration) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("CorrectPassword123!"), bcrypt.DefaultCost)
				user := entities.User{
					ID:       1,
					Username: "validuser",
					Password: string(hashedPassword),
				}
				repo.On("FindUser", mock.Anything, "validuser").Return(user, nil)
			},
			expectedError: kterrors.WrongLoginCredentialsError,
			expectedToken: dto.JWTTokens{},
		},
		{
			name:     "Token generation error",
			username: "validuser",
			password: "ValidPassword123!",
			setupMocks: func(repo *mocks.Repository, tokenGen *mocks.TokensGeneration) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("ValidPassword123!"), bcrypt.DefaultCost)
				user := entities.User{
					ID:       1,
					Username: "validuser",
					Password: string(hashedPassword),
				}
				repo.On("FindUser", mock.Anything, "validuser").Return(user, nil)
				tokenGen.On("GenerateAuthTokens", 1, "validuser").Return(dto.JWTTokens{}, errors.New("token generation failed"))
			},
			expectedError: "token generation failed",
			expectedToken: dto.JWTTokens{},
		},
		{
			name:     "Repository error",
			username: "validuser",
			password: "ValidPassword123!",
			setupMocks: func(repo *mocks.Repository, tokenGen *mocks.TokensGeneration) {
				repo.On("FindUser", mock.Anything, "validuser").Return(entities.User{}, errors.New("database error"))
			},
			expectedError: "database error",
			expectedToken: dto.JWTTokens{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			ctx := context.Background()
			mockRepo := new(mocks.Repository)
			mockTokenGen := new(mocks.TokensGeneration)
			service := authentication.NewService(mockRepo, mockTokenGen)

			// Setup mocks
			tc.setupMocks(mockRepo, mockTokenGen)

			// Create login request
			loginRequest := dto.Login{
				Username: tc.username,
				Password: tc.password,
			}

			// Act
			tokens, err := service.Login(ctx, loginRequest)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				if customErr, ok := err.(*customerror.CustomError); ok {
					assert.Equal(t, tc.expectedError, customErr.Code)
				} else {
					assert.Contains(t, err.Error(), tc.expectedError)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedToken, tokens)

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockTokenGen.AssertExpectations(t)
		})
	}
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		username      string
		password      string
		setupMocks    func(*mocks.Repository)
		expectedError string
	}{
		{
			name:     "Success",
			username: "validuser",
			password: "ValidPassword123!",
			setupMocks: func(repo *mocks.Repository) {
				repo.On("FindUser", mock.Anything, "validuser").Return(entities.User{}, errors.New("user not found"))
				repo.On("CreateUser", mock.Anything, "validuser", mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: "",
		},
		{
			name:          "Invalid username - starts with number",
			username:      "1invaliduser",
			password:      "ValidPassword123!",
			expectedError: kterrors.InvalidUsernameError,
		},
		{
			name:          "Invalid username - special characters",
			username:      "invalid@user",
			password:      "ValidPassword123!",
			expectedError: kterrors.InvalidUsernameError,
		},
		{
			name:          "Password too short",
			username:      "validuser",
			password:      "Short1!",
			expectedError: kterrors.InvalidPasswordError,
		},
		{
			name:          "Password without uppercase",
			username:      "validuser",
			password:      "password123!",
			expectedError: kterrors.InvalidPasswordError,
		},
		{
			name:          "Password without lowercase",
			username:      "validuser",
			password:      "PASSWORD123!",
			expectedError: kterrors.InvalidPasswordError,
		},
		{
			name:          "Password without number",
			username:      "validuser",
			password:      "ValidPassword!",
			expectedError: kterrors.InvalidPasswordError,
		},
		{
			name:          "Password without special character",
			username:      "validuser",
			password:      "ValidPassword123",
			expectedError: kterrors.InvalidPasswordError,
		},
		{
			name:     "Username already exists",
			username: "existinguser",
			password: "ValidPassword123!",
			setupMocks: func(repo *mocks.Repository) {
				existingUser := entities.User{
					ID:       1,
					Username: "existinguser",
					Password: "hashedpassword",
				}
				repo.On("FindUser", mock.Anything, "existinguser").Return(existingUser, nil)
			},
			expectedError: kterrors.UsernameAlreadyExistsError,
		},
		{
			name:     "Repository error on create",
			username: "validuser",
			password: "ValidPassword123!",
			setupMocks: func(repo *mocks.Repository) {
				repo.On("FindUser", mock.Anything, "validuser").Return(entities.User{}, errors.New("user not found"))
				repo.On("CreateUser", mock.Anything, "validuser", mock.AnythingOfType("string")).Return(errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockRepo := new(mocks.Repository)
			mockTokenGen := new(mocks.TokensGeneration)
			service := authentication.NewService(mockRepo, mockTokenGen)

			// Setup mocks
			if tc.setupMocks != nil {
				tc.setupMocks(mockRepo)
			}

			// Create request
			request := dto.CreateUserRequest{dto.Login{
				Username: tc.username,
				Password: tc.password,
			}}

			// Act
			err := service.CreateUser(ctx, request)

			// Assert
			if tc.expectedError != "" {
				assert.Error(t, err)
				if customErr, ok := err.(*customerror.CustomError); ok && tc.expectedError != "database error" {
					assert.Equal(t, tc.expectedError, customErr.Code)
				} else {
					assert.Contains(t, err.Error(), tc.expectedError)
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}
