package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	current_time := time.Now()
	user, err := conf.DB_Queries.CreateUser(
		r.Context(),
		database.CreateUserParams{
			UserEmail: payload.UserEmail,
			UserName:  payload.UserName,
			Password:  string(hashed_password),
			CreatedAt: pgtype.Timestamp{Time: current_time, InfinityModifier: pgtype.Finite, Valid: true},
		},
	)
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	refresh_token, _ := util.SecureRandomUtf8String(64)
	_, err = conf.DB_Queries.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:     refresh_token,
			UserID:    user.UserID,
			UserEmail: user.UserEmail,
			UserName:  user.UserName,
			ExpiredAt: pgtype.Timestamp{
				Time:             current_time.Add(time.Duration(conf.AGE_RF_TOKEN) * time.Second),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
			CreatedAt: pgtype.Timestamp{
				Time:             current_time,
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
		},
	)
	if err != nil {
		util.WriteError(w, 400, err.Error())
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
		util.WriteError(w, 500, err.Error())
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

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
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
		util.WriteError(w, 404, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		util.WriteError(w, 401, "Wrong password")
		return
	}

	current_time := time.Now()
	access_token, err := util.SignJWT(util.JwtClaims{
		Id:    user.UserID.String(),
		Email: payload.UserEmail,
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
	user_refresh_token, err := conf.DB_Queries.FindRefreshTokenByUserId(r.Context(), user.UserID)
	if err != nil {
		refresh_token, _ := util.SecureRandomUtf8String(64)
		user_refresh_token, err = conf.DB_Queries.CreateRefreshToken(
			r.Context(),
			database.CreateRefreshTokenParams{
				Token:     refresh_token,
				UserID:    user.UserID,
				UserEmail: user.UserEmail,
				UserName:  user.UserName,
				ExpiredAt: pgtype.Timestamp{
					Time: current_time.Add(
						time.Duration(conf.AGE_RF_TOKEN) * time.Second,
					),
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
				CreatedAt: pgtype.Timestamp{
					Time:             current_time,
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
			},
		)
		if err != nil {
			util.WriteError(w, 400, err.Error())
			return
		}
	} else {
		_, err = conf.DB_Queries.UpdateRefreshTokenExpiratedTime(r.Context(), database.UpdateRefreshTokenExpiratedTimeParams{
			Token: user_refresh_token.Token,
			ExpiredAt: pgtype.Timestamp{
				Time: current_time.Add(
					time.Duration(conf.AGE_RF_TOKEN) * time.Second,
				),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
		})
		if err != nil {
			util.WriteError(w, 400, err.Error())
			return
		}
	}

	util.SetCookie(w, "access_token", access_token)
	util.SetCookie(w, "refresh_token", user_refresh_token.Token)

	util.WriteJson(w, 200, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: access_token, RefreshToken: user_refresh_token.Token})
}

func HandleSignOut(w http.ResponseWriter, r *http.Request) {
	util.DeleteCookie(w, "refresh_token")
	util.WriteMessage(w, 200, "Success!")
}

func HandleGetMyInfo(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	user, err := conf.DB_Queries.FindUserById(r.Context(), pgtype.UUID{Bytes: user_uuid, Valid: true})
	if err != nil {
		util.WriteError(w, 404, err.Error())
		return
	}
	user.Password = ""
	util.WriteJson(w, 200, user)
}

func HandleGetAccessToken(w http.ResponseWriter, r *http.Request) {
	refresh_token := r.URL.Query().Get("refresh_token")
	if refresh_token == "" {
		cookie, err := r.Cookie("refresh_token")
		if err == nil {
			refresh_token = cookie.Value
		}
	}
	if refresh_token == "" {
		util.WriteError(w, 400, "missing refresh_token")
		return
	}
	rft, err := conf.DB_Queries.FindRefreshTokenByToken(r.Context(), refresh_token)
	if err != nil {
		util.WriteError(w, 404, err.Error())
		return
	}
	current_time := time.Now()
	if rft.ExpiredAt.Time.Compare(current_time) < 0 {
		util.WriteError(w, 403, "refresh_token expired")
		return
	}
	access_token, err := util.SignJWT(util.JwtClaims{
		Id:    rft.UserID.String(),
		Email: rft.UserEmail,
		Name:  rft.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(current_time),
			ExpiresAt: jwt.NewNumericDate(
				current_time.Add(time.Duration(conf.AGE_AC_TOKEN) * time.Second),
			),
			ID:      rft.UserID.String(),
			Subject: "access_token",
		},
	})
	if err != nil {
		util.WriteError(w, 403, err.Error())
		return
	}
	util.SetCookie(w, "access_token", access_token)
	util.WriteJson(w, 201, struct {
		AccessToken string `json:"access_token"`
	}{AccessToken: access_token})
}

func HandleSearchUserByName(w http.ResponseWriter, r *http.Request) {
	query_params := r.URL.Query()
	limit_str, offset_str, search := query_params.Get(
		"limit",
	), query_params.Get(
		"offset",
	), query_params.Get(
		"search",
	)
	limit, err := strconv.Atoi(limit_str)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offset_str)
	if err != nil {
		offset = 0
	}
	switch search {
	case "":
		users, err := conf.DB_Queries.GetAllUserInfo(
			r.Context(),
			database.GetAllUserInfoParams{Limit: int32(limit), Offset: int32(offset)},
		)
		if err != nil {
			util.WriteError(w, 404, err.Error())
			return
		}
		util.WriteJson(w, 200, users)
	default:
		users, err := conf.DB_Queries.GetUserInfoByName(r.Context(), database.GetUserInfoByNameParams{
			Limit:   int32(limit),
			Offset:  int32(offset),
			Column1: pgtype.Text{String: search, Valid: true},
		})
		if err != nil {
			util.WriteError(w, 404, err.Error())
			return
		}
		util.WriteJson(w, 200, users)
	}
}

func HandleGetUserInfo(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("user_id")
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	user, err := conf.DB_Queries.FindUserInfoById(r.Context(), pgtype.UUID{Bytes: user_uuid, Valid: true})
	if err != nil {
		util.WriteError(w, 404, "User not found")
		return
	}
	util.WriteJson(w, 200, user)
}

func HandleUpdateMyInfo(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	payload := UserUpdateInfoPayload{}
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
	user, err := conf.DB_Queries.FindUserById(r.Context(), pgtype.UUID{Bytes: user_uuid, Valid: true})
	if err != nil {
		util.WriteError(w, 404, "User not found")
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
		util.WriteError(w, 404, err.Error())
		return
	}
	util.WriteJson(w, 200, result)
}
