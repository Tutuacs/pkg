package types

import (
	"time"

	"github.com/Tutuacs/pkg/enums"
)

type User struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      enums.Role `json:"role"`
	CreatedAt time.Time  `json:"createdAt"`
}

type NewUserDto struct {
	Name     string     `json:"name" validate:"required"`
	Email    string     `json:"email" validate:"required,email"`
	Role     enums.Role `json:"role" validate:"required,number,min=0,max=1"`
	Password string     `json:"password" validate:"required"`
}

type UpdateUserDto struct {
	Name  string     `json:"name"`
	Email string     `json:"email" validate:"required,email"`
	Role  enums.Role `json:"role" validate:"number,min=0,max=1"`
}
