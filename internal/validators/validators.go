package validators

import (
	"github.com/asaskevich/govalidator"
)

// IsEmailValid checks if given email valid or not
func IsEmailValid(email string) bool {
	return govalidator.IsEmail(email)
}
