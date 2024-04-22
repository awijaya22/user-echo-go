package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/util"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (POST /register)
func (s *Server) Register(ctx echo.Context) error {
	req := new(generated.RegisterRequest)
	if err := ctx.Bind(req); err != nil {
		return err
	}
	if err := validateRequest(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}
	// hash password
	hashed, err := util.Generate(req.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	// insert into user table
	ID, err := s.Repository.CreateUser(ctx.Request().Context(), repository.User{
		FullName:    req.FullName,
		Password:    hashed,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}
	//return response
	var resp generated.RegisterResponse
	resp.Id = ID
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Login(ctx echo.Context) error {
	req := new(generated.LoginJSONBody)
	if err := ctx.Bind(req); err != nil {
		return err
	}
	// get user by phone number
	user, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), req.PhoneNumber)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "invalid phone number or password"})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}
	// compare password
	if !util.Check(req.Password, user.Password) {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "invalid phone number or password"})
	}
	// generate jwt RS256
	token, err := util.GenerateToken(user.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}
	// update successful login count
	err = s.Repository.UpdateSuccessfulLoginCount(ctx.Request().Context(), user.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	//return response
	var resp generated.LoginResponse
	resp.Token = token
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetMyProfile(ctx echo.Context) error {
	// parse jwt
	userID, err := util.ParseToken(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{Message: err.Error()})
	}

	// get user with id
	user, err := s.Repository.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{Message: "user not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}
	// return response
	var resp generated.MyProfileResponse
	resp.FullName = user.FullName
	resp.PhoneNumber = user.PhoneNumber
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateMyProfile(ctx echo.Context) error {
	// parse jwt
	userID, err := util.ParseToken(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{Message: err.Error()})
	}

	req := new(generated.UpdateMyProfileJSONBody)
	if err := ctx.Bind(req); err != nil {
		return err
	}

	if req.PhoneNumber != nil {
		// check if new phone number exists
		existUser, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), *req.PhoneNumber)
		if err != nil && err != repository.ErrUserNotFound {
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
		}
		if existUser.ID != "" {
			return ctx.JSON(http.StatusConflict, generated.ErrorResponse{Message: "phone number already exists"})
		}
	}

	// all check done, update user
	err = s.Repository.UpdateUser(ctx.Request().Context(), repository.UpdateUser{
		ID:          userID,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		if err == repository.ErrPhoneNumberAlreadyExists {
			return ctx.JSON(http.StatusConflict, generated.ErrorResponse{Message: err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	// get new user
	user, err := s.Repository.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,generated.ErrorResponse{Message: err.Error()})
	}

	var resp generated.MyProfileResponse
	resp.FullName = user.FullName
	resp.PhoneNumber = user.PhoneNumber
	return ctx.JSON(http.StatusOK, resp)
}

func validateRequest(req *generated.RegisterRequest) error {
	// phone number length check
	if !(len(req.PhoneNumber) >= 10 && len(req.PhoneNumber) <= 13) {
		return fmt.Errorf("invalid phone number, phone number must be between 10 and 13 digits")
	}
	// has prefix +62
	if !strings.HasPrefix(req.PhoneNumber, "+62") {
		return fmt.Errorf("invalid phone number, phone number must start with +62")
	}
	// full name length check
	if !(len(req.FullName) >= 3 && len(req.FullName) <= 60) {
		return fmt.Errorf("invalid full name, full name must be between 3 and 60 characters")
	}
	// password check
	if !(len(req.Password) >= 6 && len(req.Password) <= 64) {
		return fmt.Errorf("invalid password, password must be between 6 and 64 characters")
	}
	haveCapital, _ := regexp.MatchString(`[A-Z]`, req.Password)
	haveNumber, _ := regexp.MatchString(`[0-9]`, req.Password)
	haveSpecialChar, _ := regexp.MatchString(`[!@#$%^&*()_+{}":;'?/>.<,]`, req.Password)
	if !(haveCapital && haveNumber && haveSpecialChar) {
		return fmt.Errorf("invalid password, password must contain at least 1 capital letter, 1 number, and 1 special character")
	}
	return nil
}
