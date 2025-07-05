package dto

// Dropdown option tree
type SeleteTree struct {
	Id       int          `json:"id"`
	Label    string       `json:"label"`
	Children []SeleteTree `json:"children" gorm:"-"`
	ParentId int          `json:"-"`
}
