package handlerauth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/domain/single"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
	"github.com/google/uuid"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type userProfile struct {
	Id        uuid.UUID `json:"id"`
	Login     string    `json:"login"`
	Nickname  string    `json:"nickname"`
	IsTeacher bool      `json:"isTeacher"`
}

type UserSaver interface {
	SaveUser(context.Context, models.User) (*models.User, error)
}

func PostRegister(saver UserSaver, secret string) echo.HandlerFunc {
	type RegisterUser struct {
		Nickname string `json:"nickname" validate:"required,min=3,max=16"`
		Login    string `json:"login" validate:"required,min=3,max=16,alphanum"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	return func(c echo.Context) error {
		var u RegisterUser
		if err := c.Bind(&u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "unprocessable entity").SetInternal(err)
		}

		err := single.V.Struct(u)
		if err != nil {
			return single.ValidationError(err)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hash)

		savedUser, err := saver.SaveUser(c.Request().Context(), models.User{
			Login:        u.Login,
			Nickname:     u.Nickname,
			PasswordHash: string(hash),
		})
		switch {
		case errors.Is(err, cuteql.ErrUnique):
			return echo.NewHTTPError(http.StatusConflict, "user with this login exists").SetInternal(err)
		case err != nil:
			return echo.NewHTTPError(http.StatusConflict, "internal error").SetInternal(err)
		default:
		}

		token, err := getToken(*savedUser, secret)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, echo.Map{
			"profile": userProfile{
				Id:        savedUser.Id,
				Login:     savedUser.Login,
				Nickname:  savedUser.Nickname,
				IsTeacher: savedUser.IsTeacher,
			},
			"accessToken": token,
		})

		return nil
	}
}

type UserProvider interface {
	UserByLogin(context.Context, string) (*models.User, error)
}

func PostSignIn(provider UserProvider, secret string) echo.HandlerFunc {
	type LoginUser struct {
		Login    string `json:"login" validate:"required,min=3,max=16,alphanum"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	return func(c echo.Context) error {
		var u LoginUser
		if err := c.Bind(&u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err := single.V.Struct(&u)
		if err != nil {
			return single.ValidationError(err)
		}

		usr, err := provider.UserByLogin(c.Request().Context(), u.Login)
		switch {
		case errors.Is(err, cuteql.ErrEmptyResult):
			return echo.NewHTTPError(http.StatusNotFound, "invalid login or password")
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		if err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(u.Password)); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid login or password").SetInternal(err)
		}

		token, err := getToken(*usr, secret)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, echo.Map{
			"profile": userProfile{
				Id:        usr.Id,
				Login:     usr.Login,
				Nickname:  usr.Nickname,
				IsTeacher: usr.IsTeacher,
			},
			"accessToken": token,
		})

		return nil
	}
}

func getToken(user models.User, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"login": user.Login,
		"exp":   time.Now().Add(time.Hour * 168).Unix(),
		"tchr":  user.IsTeacher,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return tokenString, nil
}
