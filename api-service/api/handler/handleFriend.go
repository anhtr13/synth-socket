package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/anhtr13/synth-socket/api-service/internal/conf"
	"github.com/anhtr13/synth-socket/api-service/internal/database"
	"github.com/anhtr13/synth-socket/api-service/internal/util"
)

func HandleGetMyFriends(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	search := r.URL.Query().Get("search")
	limit_str := r.URL.Query().Get("limit")
	offset_str := r.URL.Query().Get("offset")
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
		friends, err := conf.DB_Queries.GetAllUserFriendInfo(r.Context(), database.GetAllUserFriendInfoParams{
			User1ID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			Limit:   int32(limit),
			Offset:  int32(offset),
		})
		if err != nil {
			util.WriteError(w, 404, "Cannot select from database", err.Error())
			return
		}
		util.WriteJson(w, 200, friends)
	default:
		friends, err := conf.DB_Queries.GetUserFriendByName(r.Context(), database.GetUserFriendByNameParams{
			User1ID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			Column2: pgtype.Text{String: search, Valid: true},
			Limit:   int32(limit),
			Offset:  int32(offset),
		})
		if err != nil {
			util.WriteError(w, 404, "Cannot select from database", err.Error())
			return
		}
		util.WriteJson(w, 200, friends)
	}
}

func HandleDeleteFriend(w http.ResponseWriter, r *http.Request) {}

func HandleGetMyFriendRequests(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	limit_str := r.URL.Query().Get("limit")
	offset_str := r.URL.Query().Get("offset")
	limit, err := strconv.Atoi(limit_str)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offset_str)
	if err != nil {
		offset = 0
	}
	fr_req, err := conf.DB_Queries.GetUserFriendRequests(r.Context(), database.GetUserFriendRequestsParams{
		ReceiverID: pgtype.UUID{Bytes: user_uuid, Valid: true},
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		util.WriteError(w, 404, "Cannot select from database", err.Error())
		return
	}
	util.WriteJson(w, 200, fr_req)
}

func HandleCreateFriendRequest(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	sender_name := user_session["name"].(string)
	sender_id := user_session["id"].(string)
	sender_uuid, _ := uuid.Parse(sender_id)
	receiver_id := r.PathValue("user_id")
	receiver_uuid, err := uuid.Parse(receiver_id)
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	cur_time := time.Now()
	req, err := conf.DB_Queries.CreateFriendRequest(r.Context(), database.CreateFriendRequestParams{
		SenderID:   pgtype.UUID{Bytes: sender_uuid, Valid: true},
		ReceiverID: pgtype.UUID{Bytes: receiver_uuid, Valid: true},
		CreatedAt:  pgtype.Timestamp{Time: cur_time, Valid: true, InfinityModifier: pgtype.Finite},
	})
	if err != nil {
		util.WriteError(w, 404, "Cannot insert into database", err.Error())
		return
	}
	noti, err := conf.DB_Queries.CreateNotification(r.Context(), database.CreateNotificationParams{
		UserID:    pgtype.UUID{Bytes: receiver_uuid, Valid: true},
		Message:   fmt.Sprintf("%s wanna be friend", sender_name),
		IDRef:     req.RequestID,
		Type:      database.NotificationTypeFriendRequest,
		CreatedAt: pgtype.Timestamp{Time: cur_time, Valid: true, InfinityModifier: pgtype.Finite},
	})
	if err != nil {
		util.WriteError(w, 404, "Cannot insert into database", err.Error())
		return
	}
	notification_msg, _ := json.Marshal(
		conf.QueueMsg_Notification{
			NotificationId: noti.NotificationID.String(),
			UserId:         noti.UserID.String(),
			Message:        noti.Message,
			Type:           conf.REQ_FRIEND,
			IdRef:          noti.IDRef.String(),
			Seen:           false,
			CreatedAt:      noti.CreatedAt.Time.String(),
		},
	)
	err = conf.RBMQ_Channel.Publish(
		conf.RBMQ_EXCHANGE_API_SOCKET,
		conf.RBMQ_KEY_NOTIFICATION,
		false,
		false,
		amqp.Publishing{
			Body: notification_msg,
		},
	)
	if err != nil {
		util.WriteError(w, 500, "Cannot publish message to queue", err.Error())
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
		util.WriteError(w, 400, "Cannot delete row from database", err.Error())
		return
	}
	util.WriteMessage(w, 200, "Success")
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
		util.WriteError(w, 400, "Cannot update row from database", err.Error())
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
		util.WriteError(w, 400, "Cannot insert into database", err.Error())
		return
	}

	err = conf.RD_Client.SAdd(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", conf.RDKEY_PREFIX_USER, friendship.User1ID.String(), conf.RDKEY_SUFFIX_FRIENDSHIPS),
		friendship.FriendshipID.String(),
	).Err()
	err = conf.RD_Client.SAdd(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", conf.RDKEY_PREFIX_USER, friendship.User2ID.String(), conf.RDKEY_SUFFIX_FRIENDSHIPS),
		friendship.FriendshipID.String(),
	).Err()
	if err != nil {
		util.WriteError(w, 500, "Cannot update cache", err.Error())
		return
	}

	room_in_msg1, _ := json.Marshal(
		conf.QueueMsg_RoomIo{
			UserId: friendship.User1ID.String(),
			RoomId: friendship.FriendshipID.String(),
			Type:   conf.ROOM_IN,
		},
	)
	room_in_msg2, _ := json.Marshal(
		conf.QueueMsg_RoomIo{
			UserId: friendship.User2ID.String(),
			RoomId: friendship.FriendshipID.String(),
			Type:   conf.ROOM_IN,
		},
	)
	err = conf.RBMQ_Channel.Publish(
		conf.RBMQ_EXCHANGE_API_SOCKET,
		conf.RBMQ_KEY_ROOM_IO,
		false,
		false,
		amqp.Publishing{
			Body: room_in_msg1,
		},
	)
	if err != nil {
		util.WriteError(w, 500, "Cannot publish message to queue", err.Error())
		return
	}
	err = conf.RBMQ_Channel.Publish(
		conf.RBMQ_EXCHANGE_API_SOCKET,
		conf.RBMQ_KEY_ROOM_IO,
		false,
		false,
		amqp.Publishing{
			Body: room_in_msg2,
		},
	)
	if err != nil {
		util.WriteError(w, 500, "Cannot publish message to queue", err.Error())
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
		util.WriteError(w, 400, "Cannot update row from database", err.Error())
		return
	}
	util.WriteMessage(w, 200, "Success")
}
