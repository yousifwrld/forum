Retrieve and Validate Cookies:

Create a function GetSessionCookie(r *http.Request) (string, error) that:
Reads the "cookie" from the request.
Returns the session ID if the cookie is present and valid.
Returns an error if the cookie is missing or invalid.
Verify Session ID:

Create a function ValidateSession(sessionID string) (bool, error) that:
Queries the database to check if the session ID is valid.
Returns true if the session is valid, otherwise false.
Middleware for Authentication:

Create a middleware AuthMiddleware(next http.Handler) http.Handler that:
Uses GetSessionCookie to retrieve the session ID.
Uses ValidateSession to check if the session is valid.
Proceeds to the next handler if the session is valid.
Redirects to the login page or returns an error if the session is invalid.
Logout Functionality:

Create a function LogoutHandler(w http.ResponseWriter, r *http.Request) that:
Reads the session ID from the cookie.
Deletes the session from the database.
Deletes the cookie from the client's browser.
Apply Middleware:

Apply AuthMiddleware to routes that require authentication.
Ensure that public routes like login and registration are accessible without authentication.
