package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/anhtr13/synth-socket/api-service/internal/conf"
	"github.com/anhtr13/synth-socket/api-service/internal/util"
	"github.com/anhtr13/synth-socket/api-service/pkgs/cache"
	"github.com/anhtr13/synth-socket/api-service/pkgs/database"
)

func HandleGetFriends(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	search := r.URL.Query().Get("search")
	limit_str := r.URL.Query().Get("limit")
	offset_str := r.URL.Query().Get("offset")
	limit, err := strconv.Atoi(limit_str)
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(offset_str)
	if err != nil {
		offset = 0
	}
	switch search {
	case "":
		friends, err := conf.DB_Queries.GetAllFriendInfoByUserId(
			r.Context(),
			database.GetAllFriendInfoByUserIdParams{
				UserID: pgtype.UUID{Bytes: user_uuid, Valid: true},
				Limit:  int32(limit),
				Offset: int32(offset),
			})
		if err != nil {
			util.WriteError(w, 404, "Cannot get friends info", err.Error())
			return
		}
		util.WriteJson(w, 200, friends)
	default:
		friends, err := conf.DB_Queries.GetFriendInfoByUserAndFriendName(r.Context(), database.GetFriendInfoByUserAndFriendNameParams{
			UserID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			Name:   pgtype.Text{String: search, Valid: true},
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			util.WriteError(w, 404, "Cannot get friends info", err.Error())
			return
		}
		util.WriteJson(w, 200, friends)
	}
}

func HandleDeleteFriend(w http.ResponseWriter, r *http.Request) {}

func HandleGetFriendRequests(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	limit_str := r.URL.Query().Get("limit")
	offset_str := r.URL.Query().Get("offset")
	limit, err := strconv.Atoi(limit_str)
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(offset_str)
	if err != nil {
		offset = 0
	}
	fr_req, err := conf.DB_Queries.GetFriendRequestsByReceiver(r.Context(), database.GetFriendRequestsByReceiverParams{
		ReceiverID: pgtype.UUID{Bytes: user_uuid, Valid: true},
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		util.WriteError(w, 404, "Cannot get friend-requests", err.Error())
		return
	}
	util.WriteJson(w, 200, fr_req)
}

func HandleAcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	request_id := r.PathValue("request_id")
	request_uuid, err := uuid.Parse(request_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	req, err := conf.DB_Queries.AcceptFriendRequest(r.Context(), database.AcceptFriendRequestParams{
		RequestID:  pgtype.UUID{Bytes: request_uuid, Valid: true},
		ReceiverID: pgtype.UUID{Bytes: user_uuid, Valid: true},
	})
	if err != nil {
		util.WriteError(w, 400, "Cannot accept friend-requests", err.Error())
		return
	}
	friendship, err := conf.DB_Queries.CreateFriendship(
		r.Context(),
		database.CreateFriendshipParams{
			User1ID: req.SenderID,
			User2ID: req.ReceiverID,
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot create friendship", err.Error())
		return
	}

	err = conf.RD_Client.SAdd(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", cache.KEY_USER, friendship.User1ID.String(), cache.USER_FRIENDS),
		friendship.User2ID.String(),
	).Err()
	err = conf.RD_Client.SAdd(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", cache.KEY_USER, friendship.User2ID.String(), cache.USER_FRIENDS),
		friendship.User1ID.String(),
	).Err()
	if err != nil {
		util.WriteError(w, 500, "Cannot update cache", err.Error())
		return
	}
	util.WriteJson(w, 200, friendship)
}

func HandleRejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	request_id := r.PathValue("request_id")
	request_uuid, err := uuid.Parse(request_id)
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	_, err = conf.DB_Queries.RejectFriendRequest(r.Context(), database.RejectFriendRequestParams{
		RequestID:  pgtype.UUID{Bytes: request_uuid, Valid: true},
		ReceiverID: pgtype.UUID{Bytes: user_uuid, Valid: true},
	})
	if err != nil {
		util.WriteError(w, 400, "Cannot reject friend-requests", err.Error())
		return
	}
	util.WriteMessage(w, 200, "Success")
}
