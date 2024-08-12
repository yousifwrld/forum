package forum

import (
	"net/http"
)

type Error struct {
	Code string
	Msg  string
}

// function renders the error html page based on the error type
func ErrorPages(w http.ResponseWriter, r *http.Request, errMessage string, code int, path ...string) {
	w.WriteHeader(code)

	var data Error
	var templatePath string

	defaultTemplatePath := "templates/errors.html"

	switch errMessage {
	// general errors
	case "404":
		data = Error{Code: "404", Msg: "The page you are looking for could not be found."}
		templatePath = defaultTemplatePath

	case "400":
		data = Error{Code: "400", Msg: "The server cannot process your request due to invalid or malformed data."}
		templatePath = defaultTemplatePath

	case "405":
		data = Error{Code: "405", Msg: "The requested method is not allowed."}
		templatePath = defaultTemplatePath

	case "500":
		data = Error{Code: "500", Msg: "The server encountered an error from our end and could not complete your request."}
		templatePath = defaultTemplatePath

	// register related errors
	case "invalid email":
		data = Error{Code: "400", Msg: "Invalid email address."}
		templatePath = path[0]

	case "invalid username":
		data = Error{Code: "400", Msg: "Username must be alphanumeric and between 3 and 30 characters."}
		templatePath = path[0]

	case "password too short":
		data = Error{Code: "400", Msg: "Password must be at least 8 characters."}
		templatePath = path[0]

	case "exists":
		data = Error{Code: "409", Msg: "User already exists."}
		templatePath = path[0]

	// login related errors
	case "user not found":
		data = Error{Code: "400", Msg: "Invalid username or password."}
		templatePath = path[0]

	case "invalid password":
		data = Error{Code: "400", Msg: "Invalid password."}
		templatePath = path[0]

	case "not logged in":
		data = Error{Code: "401", Msg: "Must be a logged in user."}
		templatePath = path[0]

	case "oauth account":
		data = Error{Code: "403", Msg: "account created through OAuth, cannot log in through form"}
		templatePath = path[0]
	default:
		// Handle any unspecified error codes if necessary
		data = Error{Code: "500", Msg: "An unexpected error occurred."}
		templatePath = defaultTemplatePath
	}

	RenderTemplate(w, templatePath, data)
}
