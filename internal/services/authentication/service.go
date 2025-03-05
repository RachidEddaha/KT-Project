package authentication

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"task/internal/dto"
	"task/internal/models/kterrors"
	"task/pkg/customerror"
	"task/pkg/database/entities"
)

const (
	pwMinLength = 8
)

type Repository interface {
	FindUser(ctx context.Context, username string) (entities.User, error)
	CreateUser(ctx context.Context, username, password string) error
}

type TokensGeneration interface {
	GenerateAuthTokens(userID int, username string) (dto.JWTTokens, error)
}

type Service struct {
	repo Repository
	tg   TokensGeneration
}

func NewService(repo Repository, tg TokensGeneration) *Service {
	return &Service{
		repo: repo,
		tg:   tg,
	}
}

func (s *Service) Login(ctx context.Context, request dto.Login) (dto.JWTTokens, error) {
	user, err := s.repo.FindUser(ctx, request.Username)
	if err != nil {
		if customerror.IsNotFoundError(err) {
			return dto.JWTTokens{}, customerror.NewCustomError(kterrors.UserNotFoundError)
		}
		return dto.JWTTokens{}, err
	}
	jwtTokens, err := s.tg.GenerateAuthTokens(user.ID, user.Username)
	if err != nil {
		return dto.JWTTokens{}, err
	}
	return jwtTokens, nil
}

func (s *Service) CreateUser(ctx context.Context, request dto.CreateUserRequest) error {
	if err := validateUsername(request.Username); err != nil {
		return err
	}
	passHashed, err := s.validateAndHashPassword(request.Password)
	if err != nil {
		return err
	}
	if err := s.checkUniqueUsername(ctx, request.Username); err != nil {
		return err
	}
	if err := s.repo.CreateUser(ctx, request.Username, passHashed); err != nil {
		return err
	}
	return nil
}

func (s *Service) checkUniqueUsername(ctx context.Context, username string) error {
	if _, err := s.repo.FindUser(ctx, username); err == nil {
		return customerror.NewI18nErrorWithParams(
			kterrors.UsernameAlreadyExistsError,
			map[string]interface{}{"username": username})
	}
	return nil
}

func (s *Service) validateAndHashPassword(password string) (string, error) {
	err := validate(password)
	if err != nil {
		return "", err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}

func validate(password string) error {
	var errorMap = make(map[string]interface{})
	if len(password) < pwMinLength {
		errorMap[kterrors.NeedAtLeastLength] = pwMinLength
	}
	if !matchRegexp("[A-Z]", password) {
		errorMap[kterrors.NeedAtLeastOneUppercaseChar] = true
	}
	if !matchRegexp("[a-z]", password) {
		errorMap[kterrors.NeedAtLeastOneLowercaseChar] = true
	}
	if !matchRegexp("[0-9]", password) {
		errorMap[kterrors.NeedAtLeastOneNumber] = true
	}
	if !matchRegexp("[\\p{P}\\p{S}]", password) {
		errorMap[kterrors.NeedAtLeastOneSpecialChar] = true
	}

	if len(errorMap) > 0 {
		return customerror.NewI18nErrorWithParams(kterrors.InvalidPasswordError, errorMap)
	}
	return nil
}

func matchRegexp(regex, src string) bool {
	matched, err := regexp.MatchString(regex, src)
	if err != nil {
		return false
	}
	return matched
}

func validateUsername(username string) error {
	if !matchRegexp(`^[A-Za-z][A-Za-z0-9]*$`, username) {
		return customerror.NewI18nErrorWithParams(
			kterrors.InvalidUsernameError,
			map[string]interface{}{"username": username})
	}
	return nil
}
