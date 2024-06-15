package logger

type Resource string

const (
	ResourceDevice  Resource = "DEVICE"
	ResourcePlugin  Resource = "PLUGIN"
	ResourceUser    Resource = "USER"
	ResourceConfig  Resource = "CONFIG"
	ResourceEnv     Resource = "ENV"
)

// Resource messages (User, Device, etc.)
var (
	MsgResourceRegisterSuccess     = NewMessageWrapper("MsgResourceRegisterSuccess", "resource %s (%s) successfully registered", InfoLevel)
	MsgResourceFetchSuccess        = NewMessageWrapper("MsgResourceFetchSuccess", "resource %s (%s) successfully fetched", InfoLevel)
	MsgResourceUpdateSuccess       = NewMessageWrapper("MsgResourceUpdateSuccess", "resource %s (%s) successfully updated", InfoLevel)
	MsgResourceRemoveSuccess       = NewMessageWrapper("MsgResourceRemoveSuccess", "resource %s (%s) successfully removed", InfoLevel)
	MsgResourceAuthenticateSuccess = NewMessageWrapper("MsgResourceAuthenticateSuccess", "resource %s (%s) successfully authenticated", InfoLevel)
)

// Uncategorizated messages
var (
	MsgResourceLoaded = NewMessageWrapper("MsgResourceLoaded", "resource %s (%s) loaded", InfoLevel)
	MsgInitializing   = NewMessageWrapper("MsgInitializing", "initializing", InfoLevel)
	MsgCleaningUp     = NewMessageWrapper("MsgCleaningUp", "cleaning up", InfoLevel)
	MsgInitialized    = NewMessageWrapper("MsgInitialized", "initialized", InfoLevel)
	MsgCleanedUp      = NewMessageWrapper("MsgCleanedUp", "cleaned up", InfoLevel)
)

// Database errors
var (
	ErrConnectingToDB      = NewMessageWrapper("ErrConnectingToDB", "error connecting to database", ErrorLevel)
	ErrGettingDBConnection = NewMessageWrapper("ErrGettingDBConnection", "error getting database connection", ErrorLevel)
	ErrClosingDBConnection = NewMessageWrapper("ErrClosingDBConnection", "error closing database connection", ErrorLevel)

	ErrResourceAlreadyExists = NewMessageWrapper("ErrResourceAlreadyExists", "resource %s (%s) already exists", ErrorLevel)
	ErrResourceNotFound      = NewMessageWrapper("ErrResourceNotFound", "resource %s (%s) not found", ErrorLevel)

	ErrRegisteringResourceInDB = NewMessageWrapper("ErrRegisteringResourceInDB", "error registering resource %s (%s) in database", ErrorLevel)
	ErrFetchingResourceFromDB  = NewMessageWrapper("ErrFetchingResourceFromDB", "error fetching resource %s (%s) from database", ErrorLevel)
	ErrUpdatingResourceInDB    = NewMessageWrapper("ErrUpdatingResourceInDB", "error updating resource %s (%s) in database", ErrorLevel)
	ErrRemovingResourceFromDB  = NewMessageWrapper("ErrRemovingResourceFromDB", "error removing resource %s (%s) from database", ErrorLevel)

	ErrInitializingResource    = NewMessageWrapper("ErrInitializingResource", "error initializing resource %s (%s)", ErrorLevel)
)

// Resource errors
var (
	ErrReadingResource    = NewMessageWrapper("ErrReadingResource", "error reading %s (%s)", FatalLevel)
	ErrFormattingResource = NewMessageWrapper("ErrFormattingResource", "error formatting %s (%s)", FatalLevel)
)


// General errors
var (
	ErrInvalidRequestPayload  = NewMessageWrapper("ErrInvalidRequestPayload", "invalid request payload", ErrorLevel)
	ErrOperationNotPermitted  = NewMessageWrapper("ErrOperationNotPermitted", "operation not permitted", ErrorLevel)
	ErrUnexpected 		      = NewMessageWrapper("ErrUnexpected", "unexpected error: %s", ErrorLevel)
)

// Uncategorized errors
var (
	ErrInitializing      = NewMessageWrapper("ErrInitializing", "error initializing: %s", ErrorLevel)
	ErrCleaningUp        = NewMessageWrapper("ErrCleaningUp", "error cleaning up: %s", ErrorLevel)
)
