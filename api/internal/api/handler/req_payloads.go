package handler

import (
	"errors"
	"regexp"
	"strings"
)

// ==============================
type UserSignUpPayload struct {
	UserEmail string `json:"user_email"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
}

func (payload *UserSignUpPayload) Validate() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	errArr := []string{}
	if !emailRegex.MatchString(payload.UserEmail) {
		errArr = append(errArr, "invalid email")
	}
	if len(payload.UserName) < 3 {
		errArr = append(errArr, "name too short")
	}
	if len(payload.Password) < 6 {
		errArr = append(errArr, "password too short")
	}
	if len(errArr) > 0 {
		errMsg := strings.Join(errArr, ", ")
		return errors.New(errMsg)
	}
	return nil
}

// ==============================
type UserSignInPayload struct {
	UserEmail string `json:"user_email"`
	Password  string `json:"password"`
}

func (payload *UserSignInPayload) Validate() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	errArr := []string{}
	if !emailRegex.MatchString(payload.UserEmail) {
		errArr = append(errArr, "invalid email")
	}
	if len(payload.Password) < 6 {
		errArr = append(errArr, "password too short")
	}
	if len(errArr) > 0 {
		errMsg := strings.Join(errArr, ", ")
		return errors.New(errMsg)
	}
	return nil
}

// ==============================
type UserUpdateInfoPayload struct {
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
}

func (user *UserUpdateInfoPayload) Validate() error {
	errArr := []string{}
	if len(user.UserName) > 0 && len(user.UserName) < 3 {
		errArr = append(errArr, "name too short")
	}
	if len(user.Password) > 0 && len(user.Password) < 6 {
		errArr = append(errArr, "password too short")
	}
	if len(user.UserName) == 0 && len(user.ProfileImage) == 0 {
		errArr = append(errArr, "nothing has changed")
	}
	if len(errArr) > 0 {
		errMsg := strings.Join(errArr, ", ")
		return errors.New(errMsg)
	}
	return nil
}

// ==============================
type CreateRoomPayload struct {
	RoomName    string `json:"room_name"`
	RoomPicture string `json:"room_picture"`
}
