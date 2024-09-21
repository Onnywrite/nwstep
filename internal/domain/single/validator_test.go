package single_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Onnywrite/nwstep/internal/domain/single"
	"github.com/go-playground/validator/v10"
)

func TestTest(t *testing.T) {
	type RegisterUser struct {
		Nickname string  `json:"nickname" validate:"required,min=3,max=16"`
		Login    string  `json:"login" validate:"required,min=3,max=16"`
		Password string  `json:"password" validate:"required,min=8,max=32"`
		Birthday *string `json:"birthday" validate:"omitempty"`
	}

	u := RegisterUser{
		Nickname: "n",
		Login:    "l",
		Password: "p",
		Birthday: nil,
	}

	v := validator.New()
	err := v.Struct(u)

	httperr := single.ValidationError(err)
	b, _ := json.Marshal(httperr)
	fmt.Println(string(b))

	t.FailNow()
}
