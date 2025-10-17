package constants

const (
	// User roles
	RoleAdmin   = "admin"
	RoleDoctor  = "doctor"
	RolePatient = "patient"

	// Error messages
	ErrUnuserorized = "unuserorized access"
	ErrNotFound     = "resource not found"
	ErrBadRequest   = "bad request"

	// Config keys
	ConfigDBHost     = "DB_HOST"
	ConfigDBUser     = "DB_USER"
	ConfigDBPassword = "DB_PASSWORD"
)
