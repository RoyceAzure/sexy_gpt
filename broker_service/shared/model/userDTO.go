package model

type userDTOKey int

const MyStructKey userDTOKey = iota

type UserDTO struct {
	UserId   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Email    string `json:"email,omitempty"`
	RoleName string `json:"role_name,omitempty"`
	RoleId   string `json:"role_id,omitempty"`
}
