package guards

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
