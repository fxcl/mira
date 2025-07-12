package constant

// Unique identifier for system users within the platform
const SYS_USER = "SYS_USER"

// Normal status
const NORMAL_STATUS = "0"

// Exception status
const EXCEPTION_STATUS = "1"

// Whether it is the system default (yes)
const IS_DEFAULT_YES = "Y"

// Whether it is the system default (no)
const IS_DEFAULT_NO = "N"

// Whether the menu is an external link (yes)
const MENU_YES_FRAME = 0

// Whether the menu is an external link (no)
const MENU_NO_FRAME = 1

// Menu type (directory)
const MENU_TYPE_DIRECTORY = "M"

// Menu type (menu)
const MENU_TYPE_MENU = "C"

// Menu type (button)
const MENU_TYPE_BUTTON = "F"

// Layout component
const LAYOUT_COMPONENT = "Layout"

// ParentView component identifier
const PARENT_VIEW_COMPONENT = "ParentView"

// InnerLink component identifier
const INNER_LINK_COMPONENT = "InnerLink"

// Request operation title key (required for operation log)
const REQUEST_TITLE = "businessTitle"

// Request operation type key (required for operation log)
const REQUEST_BUSINESS_TYPE = "businessType"

// Request operation type (specific type) (required for operation log)
const (
	REQUEST_BUSINESS_TYPE_OTHER = iota
	REQUEST_BUSINESS_TYPE_INSERT // Add
	REQUEST_BUSINESS_TYPE_UPDATE // Modify
	REQUEST_BUSINESS_TYPE_DELETE // Delete
	REQUEST_BUSINESS_TYPE_GRANT  // Grant
	REQUEST_BUSINESS_TYPE_EXPORT // Export
	REQUEST_BUSINESS_TYPE_IMPORT // Import
	REQUEST_BUSINESS_TYPE_FORCE  // Force logout
	REQUEST_BUSINESS_TYPE_GENCOD // Generate code
	REQUEST_BUSINESS_TYPE_CLEAN  // Clear data
)
