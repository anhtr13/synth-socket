package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	"github.com/anhtr13/synth-socket/api-service/internal/conf"
	"github.com/anhtr13/synth-socket/api-service/internal/util"
	"github.com/anhtr13/synth-socket/api-service/pkgs/database"
)

func HandleGetPersonalInfo(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	user, err := conf.DB_Queries.FindUserById(r.Context(), pgtype.UUID{Bytes: user_uuid, Valid: true})
	if err != nil {
		util.WriteError(w, 404, "User not found", err.Error())
		return
	}
	user.Password = ""
	util.WriteJson(w, 200, user)
}

func HandleUpdatePersonalInfo(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	payload := UserUpdateInfoPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.WriteError(w, 400, "Cannot decode payload", err.Error())
		return
	}
	err = payload.Validate()
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	user, err := conf.DB_Queries.FindUserById(r.Context(), pgtype.UUID{Bytes: user_uuid, Valid: true})
	if err != nil {
		util.WriteError(w, 404, "User not found", err.Error())
		return
	}
	params := database.UpdateUserInfoParams{
		UserID:       user.UserID,
		UserName:     user.UserName,
		ProfileImage: user.ProfileImage,
		Password:     user.Password,
	}
	if len(payload.UserName) > 0 {
		params.UserName = payload.UserName
	}
	if len(payload.ProfileImage) > 0 {
		params.ProfileImage = pgtype.Text{String: payload.ProfileImage, Valid: true}
	}
	if len(payload.Password) > 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword(
			[]byte(payload.Password),
			bcrypt.DefaultCost,
		)
		params.Password = string(hashedPassword)
	}
	result, err := conf.DB_Queries.UpdateUserInfo(r.Context(), params)
	if err != nil {
		util.WriteError(w, 400, "Cannot update user", err.Error())
		return
	}
	util.WriteJson(w, 200, result)
}
