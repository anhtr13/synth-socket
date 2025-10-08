package handler

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	"github.com/anhtr13/synth-socket/api-service/internal/conf"
	"github.com/anhtr13/synth-socket/api-service/internal/database"
	"github.com/anhtr13/synth-socket/api-service/internal/util"
)

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	payload := UserSignUpPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	err = payload.Validate()
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	hashed_password, err := bcrypt.GenerateFromPassword(
		[]byte(payload.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		util.WriteError(w, 500, err.Error())
		return
	}
	user, err := conf.DB_Queries.CreateUser(
		r.Context(),
		database.CreateUserParams{
			UserEmail: payload.UserEmail,
			UserName:  payload.UserName,
			Password:  string(hashed_password),
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot create user", err.Error())
		return
	}

	current_time := time.Now()
	refresh_token, _ := util.SecureRandomUtf8String(64)
	_, err = conf.DB_Queries.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			UserID:    user.UserID,
			UserEmail: user.UserEmail,
			UserName:  user.UserName,
			Token:     refresh_token,
			ExpiredAt: pgtype.Timestamp{
				Time:             current_time.Add(conf.AGE_RF_TOKEN * time.Second),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
		},
	)
	if err != nil {
		util.WriteError(w, 500, "Cannot store refresh_token", err.Error())
		return
	}

	access_token, err := util.SignJWT(
		util.JwtClaims{
			Id:    user.UserID.String(),
			Email: user.UserEmail,
			Name:  user.UserName,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt: jwt.NewNumericDate(current_time),
				ExpiresAt: jwt.NewNumericDate(
					current_time.Add(time.Duration(conf.AGE_AC_TOKEN) * time.Second),
				),
				ID:      user.UserID.String(),
				Subject: "access_token",
			},
		},
	)
	if err != nil {
		util.WriteError(w, 500, "Cannot create access_token", err.Error())
		return
	}

	util.SetCookie(w, "access_token", access_token)
	util.SetCookie(w, "refresh_token", refresh_token)

	util.WriteJson(w, 201, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	payload := UserSignInPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	if err := payload.Validate(); err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}

	user, err := conf.DB_Queries.FindUserByEmail(r.Context(), payload.UserEmail)
	if err != nil {
		util.WriteError(w, 404, "User not found", err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		util.WriteError(w, 401, "Wrong password")
		return
	}

	current_time := time.Now()
	rf_token, err := conf.DB_Queries.UpdateTokenExpiratedTimeByUserId(
		r.Context(),
		database.UpdateTokenExpiratedTimeByUserIdParams{
			UserID: user.UserID,
			ExpiredAt: pgtype.Timestamp{
				Time:             current_time.Add(conf.AGE_RF_TOKEN * time.Second),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
		},
	)
	if err != nil {
		refresh_token, _ := util.SecureRandomUtf8String(64)
		rf_token, err = conf.DB_Queries.CreateRefreshToken(
			r.Context(),
			database.CreateRefreshTokenParams{
				UserID:    user.UserID,
				UserEmail: user.UserEmail,
				UserName:  user.UserName,
				Token:     refresh_token,
				ExpiredAt: pgtype.Timestamp{
					Time:             current_time.Add(conf.AGE_RF_TOKEN * time.Second),
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
			},
		)
	}
	if err != nil {
		util.WriteError(w, 500, "Cannot update or create refresh_token", err.Error())
		return
	}
	access_token, err := util.SignJWT(util.JwtClaims{
		Id:    user.UserID.String(),
		Email: user.UserEmail,
		Name:  user.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(current_time),
			ExpiresAt: jwt.NewNumericDate(
				current_time.Add(time.Duration(conf.AGE_AC_TOKEN) * time.Second),
			),
			ID:      user.UserID.String(),
			Subject: "access_token",
		},
	})
	if err != nil {
		util.WriteError(w, 500, err.Error())
		return
	}

	util.SetCookie(w, "access_token", access_token)
	util.SetCookie(w, "refresh_token", rf_token.Token)

	util.WriteJson(w, 200, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: access_token, RefreshToken: rf_token.Token})
}

func HandleLogOut(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	conf.DB_Queries.UpdateTokenExpiratedTimeByUserId(
		r.Context(),
		database.UpdateTokenExpiratedTimeByUserIdParams{
			UserID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			ExpiredAt: pgtype.Timestamp{
				Time:             time.Now(),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
		},
	)
	util.DeleteCookie(w, "refresh_token")
	util.DeleteCookie(w, "access_token")
	util.WriteMessage(w, 200, "Success!")
}

func HandleCreateAccessToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("refresh_token")
	if token == "" {
		util.WriteError(w, 400, "refresh_token not found")
		return
	}
	refresh_token, err := conf.DB_Queries.FindRefreshTokenByToken(r.Context(), token)
	if err != nil {
		util.WriteError(w, 404, "refresh_token not found", err.Error())
		return
	}

	current_time := time.Now()
	if refresh_token.ExpiredAt.Time.Before(current_time) {
		util.WriteError(w, 403, "refresh_token expired")
		return
	}

	access_token, err := util.SignJWT(util.JwtClaims{
		Id:    refresh_token.UserID.String(),
		Email: refresh_token.UserEmail,
		Name:  refresh_token.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(current_time),
			ExpiresAt: jwt.NewNumericDate(
				current_time.Add(time.Duration(conf.AGE_AC_TOKEN) * time.Second),
			),
			ID:      refresh_token.UserID.String(),
			Subject: "access_token",
		},
	})
	if err != nil {
		util.WriteError(w, 500, "Cannot create access_token", err.Error())
		return
	}
	util.SetCookie(w, "access_token", access_token)
	util.WriteJson(w, 201, struct {
		AccessToken string `json:"access_token"`
	}{AccessToken: access_token})
}
