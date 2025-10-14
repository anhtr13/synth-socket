package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/anhtr13/synth-socket/api/internal/conf"
	"github.com/anhtr13/synth-socket/api/internal/util"
	"github.com/anhtr13/synth-socket/api/pkgs/database"
	"github.com/anhtr13/synth-socket/api/pkgs/queue"
)

func HandleGetAllUsersInfo(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
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
		limit = 20
	}
	offset, err := strconv.Atoi(offset_str)
	if err != nil {
		offset = 0
	}
	switch search {
	case "":
		users, err := conf.DB_Queries.GetAllUserInfo(
			r.Context(),
			database.GetAllUserInfoParams{
				UserID: pgtype.UUID{Bytes: user_uuid, Valid: true},
				Limit:  int32(limit),
				Offset: int32(offset),
			},
		)
		if err != nil {
			util.WriteError(w, 404, "Cannot find users", err.Error())
			return
		}
		util.WriteJson(w, 200, users)
	default:
		users, err := conf.DB_Queries.GetAllUserInfoByName(
			r.Context(),
			database.GetAllUserInfoByNameParams{
				UserID: pgtype.UUID{Bytes: user_uuid, Valid: true},
				Name:   pgtype.Text{String: search, Valid: true},
				Limit:  int32(limit),
				Offset: int32(offset),
			})
		if err != nil {
			util.WriteError(w, 404, "Cannot find users", err.Error())
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

func HandleCreateFriendRequest(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	sender_name := user_session["name"].(string)
	sender_id := user_session["id"].(string)
	sender_uuid, _ := uuid.Parse(sender_id)
	receiver_id := r.PathValue("user_id")
	receiver_uuid, err := uuid.Parse(receiver_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	req, err := conf.DB_Queries.CreateFriendRequest(
		r.Context(),
		database.CreateFriendRequestParams{
			SenderID:   pgtype.UUID{Bytes: sender_uuid, Valid: true},
			ReceiverID: pgtype.UUID{Bytes: receiver_uuid, Valid: true},
		},
	)
	if err != nil {
		util.WriteError(w, 404, "Cannot create request", err.Error())
		return
	}
	noti, err := conf.DB_Queries.CreateNotification(
		r.Context(),
		database.CreateNotificationParams{
			UserID:  pgtype.UUID{Bytes: receiver_uuid, Valid: true},
			Message: fmt.Sprintf("%s wanna be friend", sender_name),
			IDRef:   req.RequestID,
			Type:    database.NotificationTypeFriendRequest,
		},
	)
	if err != nil {
		util.WriteError(w, 404, "Cannot create notification", err.Error())
		return
	}
	notification_msg, _ := json.Marshal(
		queue.Notification{
			NotificationId: noti.NotificationID.String(),
			UserId:         noti.UserID.String(),
			Message:        noti.Message,
			Type:           queue.REQ_FRIEND,
			IdRef:          noti.IDRef.String(),
			Seen:           false,
			CreatedAt:      noti.CreatedAt.Time.String(),
		},
	)
	err = conf.RBMQ_Channel.Publish(
		queue.EXCHANGE_API_TO_SOCKET,
		queue.ROUTE_NOTIFICATION,
		false,
		false,
		amqp.Publishing{
			Body: notification_msg,
		},
	)
	if err != nil {
		util.WriteError(w, 500, "Cannot publish notification", err.Error())
		return
	}
	util.WriteJson(w, 200, req)
}

func HandleDeleteFriendRequest(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	sender_id := user_session["id"].(string)
	sender_uuid, _ := uuid.Parse(sender_id)
	receiver_id := r.PathValue("user_id")
	receiver_uuid, err := uuid.Parse(receiver_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	err = conf.DB_Queries.DeleteFriendRequest(r.Context(), database.DeleteFriendRequestParams{
		SenderID:   pgtype.UUID{Bytes: sender_uuid, Valid: true},
		ReceiverID: pgtype.UUID{Bytes: receiver_uuid, Valid: true},
	})
	if err != nil {
		util.WriteError(w, 400, "Cannot delete request", err.Error())
		return
	}
	util.WriteMessage(w, 200, "Success")
}
