package lang

const (
	SuccessMsg = "success"
)

const (
	ErrorUserIdNotFound      = "User id not found"
	ErrorUserIdNotFoundf     = "User id %d not found"
	ErrorConflictPhoneNumber = "Conflict phone number already exist"
	ErrorConflictEmail       = "Conflict email already exist"
	ErrorPhoneNumberNotExist = "Phone number not found, please register"
	ErrorUnsuccessfulLogin   = "Unsuccessful Login, please insert correct password"
	ErrorGenerateToken       = "Error generating token"

	ErrorSecurityScheme = "security scheme %s != 'BearerAuth'"

	SuccessUpdateUserById = "Success update user id %d"
	SuccesLogin           = "Success Login"
)
