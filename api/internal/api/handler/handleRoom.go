package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/anhtr13/synth-socket/api/internal/cache"
	"github.com/anhtr13/synth-socket/api/internal/conf"
	"github.com/anhtr13/synth-socket/api/internal/database"
	"github.com/anhtr13/synth-socket/api/internal/queue"
	"github.com/anhtr13/synth-socket/api/internal/util"
)

func HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	payload := CreateRoomPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.WriteError(w, 400, "Cannot decode request body", err.Error())
		return
	}
	if payload.RoomName == "" {
		util.WriteError(w, 400, "Invalid room name")
		return
	}
	room, err := conf.DB_Queries.CreateRoom(
		r.Context(),
		database.CreateRoomParams{
			RoomName:    payload.RoomName,
			RoomPicture: pgtype.Text{String: payload.RoomPicture},
			CreatedBy:   pgtype.UUID{Bytes: user_uuid, Valid: true},
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot create room", err.Error())
		return
	}
	_, err = conf.DB_Queries.CreateRoomMember(
		r.Context(),
		database.CreateRoomMemberParams{
			RoomID:   room.RoomID,
			MemberID: pgtype.UUID{Bytes: user_uuid, Valid: true},
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot create member", err.Error())
		return
	}
	err = conf.RD_Client.SAdd(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user_id, cache.USER_ROOMS),
		room.RoomID.String(),
	).Err()
	if err != nil {
		util.WriteError(w, 500, "Cannot update cache", err.Error())
		return
	}
	room_in_msg, _ := json.Marshal(
		queue.RoomIo{
			UserId: user_id,
			RoomId: room.RoomID.String(),
			Type:   queue.ROOM_IN,
		},
	)
	err = conf.RBMQ_Channel.Publish(
		queue.EXCHANGE_API_TO_SOCKET,
		queue.ROUTE_ROOM_IO,
		false,
		false,
		amqp.Publishing{
			Body: room_in_msg,
		},
	)
	if err != nil {
		util.WriteError(w, 500, "Cannot publish message to queue", err.Error())
		return
	}
	util.WriteJson(w, 201, room)
}

func HandleGetRoomMembers(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	room_id := r.PathValue("room_id")
	room_uuid, err := uuid.Parse(room_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
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
	ok, err := conf.RD_Client.SIsMember(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user_id, cache.USER_ROOMS),
		room_id,
	).Result()
	if err != nil || !ok {
		util.WriteError(w, 403, "Not a room member")
		return
	}
	members, err := conf.DB_Queries.GetRoomMemberInfoByRoomId(
		r.Context(),
		database.GetRoomMemberInfoByRoomIdParams{
			RoomID: pgtype.UUID{Bytes: room_uuid, Valid: true},
			Limit:  int32(limit),
			Offset: int32(offset),
		},
	)
	if err != nil {
		util.WriteError(w, 404, "Cannot get room's members", err.Error())
		return
	}
	util.WriteJson(w, 200, members)
}

func HandleGetOwnedRooms(w http.ResponseWriter, r *http.Request) {
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
		rooms, err := conf.DB_Queries.GetRoomsByCreator(
			r.Context(),
			database.GetRoomsByCreatorParams{
				CreatedBy: pgtype.UUID{Bytes: user_uuid, Valid: true},
				Limit:     int32(limit),
				Offset:    int32(offset),
			},
		)
		if err != nil {
			util.WriteError(w, 404, "Cannot get owned rooms", err.Error())
			return
		}
		util.WriteJson(w, 200, rooms)
	default:
		rooms, err := conf.DB_Queries.GetRoomsByCreatorAndName(
			r.Context(),
			database.GetRoomsByCreatorAndNameParams{
				CreatedBy: pgtype.UUID{Bytes: user_uuid, Valid: true},
				Limit:     int32(limit),
				Offset:    int32(offset),
				RoomName:  pgtype.Text{String: search, Valid: true},
			},
		)
		if err != nil {
			util.WriteError(w, 404, "Cannot get owned rooms", err.Error())
			return
		}
		util.WriteJson(w, 200, rooms)
	}
}

func HandleGetRoomData(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	room_id := r.PathValue("room_id")
	room_uuid, err := uuid.Parse(room_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	ok, err := conf.RD_Client.SIsMember(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user_id, cache.USER_ROOMS),
		room_id,
	).Result()
	if err != nil || !ok {
		util.WriteError(w, 403, "Not a room member")
		return
	}
	room, err := conf.DB_Queries.FindRoomById(r.Context(), pgtype.UUID{Bytes: room_uuid, Valid: true})
	if err != nil || !ok {
		util.WriteError(w, 404, "Room not found", err.Error())
		return
	}
	util.WriteJson(w, 200, room)
}

func HandleGetAllRooms(w http.ResponseWriter, r *http.Request) {
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
		rooms, err := conf.DB_Queries.GetRoomsDataByMemberId(
			r.Context(),
			database.GetRoomsDataByMemberIdParams{
				Limit:    int32(limit),
				Offset:   int32(offset),
				MemberID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			},
		)
		if err != nil {
			util.WriteError(w, 404, "Cannot get user's rooms", err.Error())
			return
		}
		util.WriteJson(w, 200, rooms)
	default:
		rooms, err := conf.DB_Queries.GetRoomsDataByMemberIdAndRoomName(
			r.Context(),
			database.GetRoomsDataByMemberIdAndRoomNameParams{
				Limit:    int32(limit),
				Offset:   int32(offset),
				Name:     pgtype.Text{String: search, Valid: true},
				MemberID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			},
		)
		if err != nil {
			util.WriteError(w, 404, "Cannot get user's rooms", err.Error())
			return
		}
		util.WriteJson(w, 200, rooms)
	}
}

func HandleLeaveRoom(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	room_id := r.PathValue("room_id")
	room_uuid, err := uuid.Parse(room_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	err = conf.DB_Queries.DeleteRoomMemmber(
		r.Context(),
		database.DeleteRoomMemmberParams{
			RoomID:   pgtype.UUID{Bytes: room_uuid, Valid: true},
			MemberID: pgtype.UUID{Bytes: user_uuid, Valid: true},
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot delete room member", err.Error())
		return
	}
	err = conf.RD_Client.SRem(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user_id, cache.USER_ROOMS),
		room_id,
	).Err()
	if err != nil {
		util.WriteError(w, 500, "Cannot update cache", err.Error())
		return
	}
	room_out_msg, _ := json.Marshal(
		queue.RoomIo{
			UserId: user_id,
			RoomId: room_id,
			Type:   queue.ROOM_OUT,
		})
	err = conf.RBMQ_Channel.Publish(
		queue.EXCHANGE_API_TO_SOCKET,
		queue.ROUTE_ROOM_IO,
		false,
		false,
		amqp.Publishing{
			Body: room_out_msg,
		},
	)
	if err != nil {
		util.WriteError(w, 500, "Cannot publish message to queue", err.Error())
		return
	}
	util.WriteMessage(w, 200, "Success")
}

func HandleGetRoomInvites(w http.ResponseWriter, r *http.Request) {
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
	reqs, err := conf.DB_Queries.GetRoomInvitesByReceiverId(
		r.Context(),
		database.GetRoomInvitesByReceiverIdParams{
			ReceiverID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			Limit:      int32(limit),
			Offset:     int32(offset),
		},
	)
	if err != nil {
		util.WriteError(w, 404, "Cannot get invitation", err.Error())
		return
	}
	util.WriteJson(w, 200, reqs)
}

func HandleAcceptRoomInvite(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	req_id := r.PathValue("invite_id")
	req_uuid, err := uuid.Parse(req_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	req, err := conf.DB_Queries.AcceptRoomInvite(
		r.Context(),
		database.AcceptRoomInviteParams{
			ReceiverID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			InviteID:   pgtype.UUID{Bytes: req_uuid, Valid: true},
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot accept invite", err.Error())
		return
	}
	room_member, err := conf.DB_Queries.CreateRoomMember(
		r.Context(),
		database.CreateRoomMemberParams{
			RoomID:   req.RoomID,
			MemberID: req.ReceiverID,
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot create room member", err.Error())
		return
	}
	err = conf.RD_Client.SAdd(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user_id, cache.USER_ROOMS),
		req.RoomID.String(),
	).Err()
	if err != nil {
		util.WriteError(w, 500, "Cannot update cache", err.Error())
		return
	}
	room_in_msg, _ := json.Marshal(
		queue.RoomIo{
			UserId: user_id,
			RoomId: req.RoomID.String(),
			Type:   queue.ROOM_IN,
		},
	)
	err = conf.RBMQ_Channel.Publish(
		queue.EXCHANGE_API_TO_SOCKET,
		queue.ROUTE_ROOM_IO,
		false,
		false,
		amqp.Publishing{
			Body: room_in_msg,
		},
	)
	if err != nil {
		util.WriteError(w, 500, "Cannot publish message to queue", err.Error())
		return
	}
	util.WriteJson(w, 200, room_member)
}

func HandleRejectRoomInvite(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	req_id := r.PathValue("invite_id")
	req_uuid, err := uuid.Parse(req_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	_, err = conf.DB_Queries.RejectRoomInvite(
		r.Context(),
		database.RejectRoomInviteParams{
			ReceiverID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			InviteID:   pgtype.UUID{Bytes: req_uuid, Valid: true},
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot reject invite", err.Error())
		return
	}
	util.WriteMessage(w, 200, "Success")
}

func HandleGetAllInvitesToRoom(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	room_id := r.PathValue("room_id")
	room_uuid, err := uuid.Parse(room_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	room, err := conf.DB_Queries.FindRoomById(r.Context(), pgtype.UUID{Bytes: room_uuid, Valid: true})
	if err != nil {
		util.WriteError(w, 404, "Cannot find room", err.Error())
		return
	}
	if room.CreatedBy.Bytes != user_uuid {
		util.WriteError(w, 403, "Not the room's owner")
		return
	}
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
	room_invites, err := conf.DB_Queries.GetRoomInvitesByRoomId(
		r.Context(),
		database.GetRoomInvitesByRoomIdParams{
			RoomID: pgtype.UUID{Bytes: room_uuid, Valid: true},
			Limit:  int32(limit),
			Offset: int32(offset),
		},
	)
	if err != nil {
		util.WriteError(w, 404, "Cannot get room's invitations", err.Error())
		return
	}
	util.WriteJson(w, 200, room_invites)
}

func HandleCreateRoomInvite(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_name := user_session["name"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	target_id := r.PathValue("target_id")
	target_uuid, err := uuid.Parse(target_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	room_id := r.PathValue("room_id")
	room_uuid, err := uuid.Parse(room_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	room, err := conf.DB_Queries.FindRoomById(r.Context(), pgtype.UUID{Bytes: room_uuid, Valid: true})
	if err != nil {
		util.WriteError(w, 404, "Cannot find room", err.Error())
		return
	}
	if room.CreatedBy.Bytes != user_uuid {
		util.WriteError(w, 403, "Not the room's owner")
		return
	}
	req, err := conf.DB_Queries.CreateRoomInvite(
		r.Context(),
		database.CreateRoomInviteParams{
			RoomID:     pgtype.UUID{Bytes: room_uuid, Valid: true},
			SenderID:   pgtype.UUID{Bytes: user_uuid, Valid: true},
			ReceiverID: pgtype.UUID{Bytes: target_uuid, Valid: true},
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot create invite", err.Error())
		return
	}
	noti, err := conf.DB_Queries.CreateNotification(r.Context(), database.CreateNotificationParams{
		UserID:  pgtype.UUID{Bytes: target_uuid, Valid: true},
		Message: fmt.Sprintf("%s invite you to room %s", user_name, room.RoomName),
		IDRef:   req.InviteID,
		Type:    database.NotificationTypeRoomInvite,
	})
	if err != nil {
		util.WriteError(w, 400, "Cannot create notification", err.Error())
		return
	}
	notification_msg, _ := json.Marshal(
		queue.Notification{
			NotificationId: noti.NotificationID.String(),
			UserId:         noti.UserID.String(),
			Message:        noti.Message,
			Type:           queue.ROOM_INVITE,
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
		util.WriteError(w, 400, "Cannot publish message to queue", err.Error())
		return
	}
	util.WriteJson(w, 200, req)
}

func HandleDeleteRoomInvite(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	target_id := r.PathValue("target_id")
	target_uuid, err := uuid.Parse(target_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	err = conf.DB_Queries.DeleteRoomInvite(
		r.Context(),
		database.DeleteRoomInviteParams{
			SenderID:   pgtype.UUID{Bytes: user_uuid, Valid: true},
			ReceiverID: pgtype.UUID{Bytes: target_uuid, Valid: true},
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot delete invitation", err.Error())
		return
	}
	util.WriteMessage(w, 200, "Success")
}

func HandleGetRoomMessages(w http.ResponseWriter, r *http.Request) {
	room_id := r.PathValue("room_id")
	room_uuid, err := uuid.Parse(room_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
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
	msgs, err := conf.DB_Queries.GetMesssagesByRoomId(
		r.Context(),
		database.GetMesssagesByRoomIdParams{
			RoomID: pgtype.UUID{Bytes: room_uuid, Valid: true},
			Limit:  int32(limit),
			Offset: int32(offset),
		},
	)
	if err != nil {
		util.WriteError(w, 404, "Cannot get messages", err.Error())
		return
	}
	util.WriteJson(w, 200, msgs)
}
