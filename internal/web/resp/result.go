package resp

// Code definitions
const (
	// CodeSuccess indicates that the request was processed successfully.
	CodeSuccess = 0

	// CodeInvalidParam indicates that the client provided invalid parameters (e.g., missing fields, wrong format).
	// This usually maps to a 400 Bad Request in RESTful terms.
	CodeInvalidParam = 40001

	// CodeUserExist indicates that the user registration failed because the email already exists.
	// This prevents duplicate accounts.
	CodeUserExist = 40101

	// CodeInvalidCreds indicates that the login failed due to incorrect email or password.
	CodeInvalidCreds = 40102

	// CodeUserNotFound indicates that the requested user does not exist.
	CodeUserNotFound = 40103

	// CodeServerBusy indicates an internal server error or unexpected failure.
	// This maps to a 500 Internal Server Error, telling the client to retry later.
	CodeServerBusy = 50001
)

// Result is the unified response structure for all API endpoints.
// It ensures that the frontend always receives a consistent JSON format,
// regardless of whether the request succeeded or failed.
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
