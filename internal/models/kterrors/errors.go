package kterrors

const (
	UsernameAlreadyExistsError  = "USERNAME_ALREADY_EXISTS"
	InvalidPasswordError        = "INVALID_PASSWORD_ERROR"
	NeedAtLeastLength           = "NEED_AT_LEAST_LENGTH"
	NeedAtLeastOneUppercaseChar = "NEED_AT_LEAST_ONE_UPPERCASE_CHAR"
	NeedAtLeastOneLowercaseChar = "NEED_AT_LEAST_ONE_LOWERCASE_CHAR"
	NeedAtLeastOneNumber        = "NEED_AT_LEAST_ONE_NUMBER"
	NeedAtLeastOneSpecialChar   = "NEED_AT_LEAST_ONE_SPECIAL_CHAR"
	InvalidUsernameError        = "INVALID_USERNAME_ERROR"

	UserNotFoundError           = "USER_NOT_FOUND_ERROR"
	UserCannotDeleteFilmError   = "USER_CANNOT_DELETE_FILM_ERROR"
	FilmTitleAlreadyExistsError = "FILM_TITLE_ALREADY_EXISTS_ERROR"
	UserCannotUpdateFilmError   = "USER_CANNOT_UPDATE_FILM_ERROR"
)
