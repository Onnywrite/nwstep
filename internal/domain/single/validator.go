package single

import (
	"errors"
	"net/http"

	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	rutranslations "github.com/go-playground/validator/v10/translations/ru"
	"github.com/labstack/echo/v4"
)

var (
	V   = validator.New()
	uni *ut.UniversalTranslator
)

func ValidationError(from error) error {
	rus := ru.New()
	uni = ut.New(rus, rus)

	trans, _ := uni.GetTranslator("ru")

	V = validator.New()
	rutranslations.RegisterDefaultTranslations(V, trans)

	var ve validator.ValidationErrors

	if errors.As(from, &ve) {
		// out := make([]echo.Map, len(ve))
		// for i, fe := range ve {
		// 	out[i] = echo.Map{"field": fe.Field(), "message": fe.Error()}
		// }

		return echo.NewHTTPError(http.StatusBadRequest, ve.Translate(trans))
	}

	return echo.NewHTTPError(http.StatusBadRequest, from.Error())
}
