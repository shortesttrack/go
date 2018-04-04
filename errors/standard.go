package errors

var (
	NotFound = New("Not found")
	BadRequest = New("Bad request").SetStatus(400)
	Unauthorized = New("Unauthorized").SetStatus(401)
	Forbidden = New("Forbidden").SetStatus(403)
	InternalServerError = New("Internal server error").SetStatus(500)
	AccessDenied = New("Access denied")
)