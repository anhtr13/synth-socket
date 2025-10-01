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

func HandleCreateGroup(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	payload := CreateGroupPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	if payload.GroupName == "" {
		util.WriteError(w, 400, "Invalid group name")
		return
	}
	_, err = conf.DB_Queries.FindGroupByCreatorAndName(
		r.Context(),
		database.FindGroupByCreatorAndNameParams{
			CreatedBy: pgtype.UUID{Bytes: user_uuid, Valid: true},
			GroupName: payload.GroupName,
		},
	)
	if err == nil {
		util.WriteError(w, 400, "Duplicated group name")
		return
	}
	cur_time := pgtype.Timestamp{
		Time:             time.Now(),
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
	group, err := conf.DB_Queries.CreateGroup(
		r.Context(),
		database.CreateGroupParams{
			GroupName:    payload.GroupName,
			GroupPicture: pgtype.Text{String: payload.GroupPicture},
			CreatedBy:    pgtype.UUID{Bytes: user_uuid, Valid: true},
			CreatedAt:    cur_time,
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot insert into database", err.Error())
		return
	}
	_, err = conf.DB_Queries.CreateGroupMember(
		r.Context(),
		database.CreateGroupMemberParams{
			GroupID:  group.GroupID,
			MemberID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			JoinedAt: cur_time,
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot insert into database", err.Error())
		return
	}
	err = conf.RD_Client.SAdd(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", conf.RDKEY_PREFIX_USER, user_id, conf.RDKEY_SUFFIX_GROUPS),
		group.GroupID.String(),
	).Err()
	if err != nil {
		util.WriteError(w, 500, "Cannot update cache", err.Error())
		return
	}
	room_in_msg, _ := json.Marshal(
		conf.QueueMsg_RoomIo{
			UserId: user_id,
			RoomId: group.GroupID.String(),
			Type:   conf.ROOM_IN,
		},
	)
	err = conf.RBMQ_Channel.Publish(
		conf.RBMQ_EXCHANGE_API_SOCKET,
		conf.RBMQ_KEY_ROOM_IO,
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
	util.WriteJson(w, 201, group)
}

func HandleGetGroupMembers(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	group_id := r.PathValue("group_id")
	group_uuid, err := uuid.Parse(group_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
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
	ok, err := conf.RD_Client.SIsMember(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", conf.RDKEY_PREFIX_USER, user_id, conf.RDKEY_SUFFIX_GROUPS),
		group_id,
	).Result()
	if err != nil || !ok {
		util.WriteError(w, 403, "You are not in this group")
		return
	}
	members, err := conf.DB_Queries.GetAllGroupMemberInfo(
		r.Context(),
		database.GetAllGroupMemberInfoParams{
			GroupID: pgtype.UUID{Bytes: group_uuid, Valid: true},
			Limit:   int32(limit),
			Offset:  int32(offset),
		},
	)
	if err != nil {
		util.WriteError(w, 404, "Cannot select from database", err.Error())
		return
	}
	util.WriteJson(w, 200, members)
}

func HandleGetMyGroups(w http.ResponseWriter, r *http.Request) {
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
	groups, err := conf.DB_Queries.GetAllUserGroups(r.Context(), database.GetAllUserGroupsParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		MemberID: pgtype.UUID{Bytes: user_uuid, Valid: true},
	})
	if err != nil {
		util.WriteError(w, 404, "Cannot select from database", err.Error())
		return
	}
	util.WriteJson(w, 200, groups)
}

func HandleLeaveGroup(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	group_id := r.PathValue("group_id")
	group_uuid, err := uuid.Parse(group_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	err = conf.DB_Queries.DeleteGroupMemmber(r.Context(), database.DeleteGroupMemmberParams{
		GroupID:  pgtype.UUID{Bytes: group_uuid, Valid: true},
		MemberID: pgtype.UUID{Bytes: user_uuid, Valid: true},
	})
	if err != nil {
		util.WriteError(w, 400, "Cannot delete row from database", err.Error())
		return
	}
	err = conf.RD_Client.SRem(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", conf.RDKEY_PREFIX_USER, user_id, conf.RDKEY_SUFFIX_GROUPS),
		group_id,
	).Err()
	if err != nil {
		util.WriteError(w, 500, "Cannot update cache", err.Error())
		return
	}
	room_out_msg, _ := json.Marshal(conf.QueueMsg_RoomIo{
		UserId: user_id,
		RoomId: group_id,
		Type:   conf.ROOM_OUT,
	})
	err = conf.RBMQ_Channel.Publish(
		conf.RBMQ_EXCHANGE_API_SOCKET,
		conf.RBMQ_KEY_ROOM_IO,
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

func HandleGetMyGroupInvites(w http.ResponseWriter, r *http.Request) {
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
	reqs, err := conf.DB_Queries.GetUserGroupInvites(r.Context(), database.GetUserGroupInvitesParams{
		ReceiverID: pgtype.UUID{Bytes: user_uuid, Valid: true},
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		util.WriteError(w, 404, "Cannot select from database", err.Error())
		return
	}
	util.WriteJson(w, 200, reqs)
}

func HandleAcceptGroupInvite(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	req_id := r.PathValue("invite_id")
	req_uuid, err := uuid.Parse(req_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	req, err := conf.DB_Queries.AcceptGroupInvite(r.Context(), database.AcceptGroupInviteParams{
		ReceiverID: pgtype.UUID{Bytes: user_uuid, Valid: true},
		InviteID:   pgtype.UUID{Bytes: req_uuid, Valid: true},
	})
	if err != nil {
		util.WriteError(w, 400, "Cannot update row in database", err.Error())
		return
	}
	group_member, err := conf.DB_Queries.CreateGroupMember(
		r.Context(),
		database.CreateGroupMemberParams{
			GroupID:  req.GroupID,
			MemberID: req.ReceiverID,
			JoinedAt: pgtype.Timestamp{Time: time.Now(), InfinityModifier: pgtype.Finite, Valid: true},
		},
	)
	if err != nil {
		util.WriteError(w, 400, "Cannot insert to database", err.Error())
		return
	}
	err = conf.RD_Client.SAdd(
		r.Context(),
		fmt.Sprintf("%s:%s:%s", conf.RDKEY_PREFIX_USER, user_id, conf.RDKEY_SUFFIX_GROUPS),
		req.GroupID.String(),
	).Err()
	if err != nil {
		util.WriteError(w, 500, "Cannot update cache", err.Error())
		return
	}
	room_in_msg, _ := json.Marshal(
		conf.QueueMsg_RoomIo{
			UserId: user_id,
			RoomId: req.GroupID.String(),
			Type:   conf.ROOM_IN,
		},
	)
	err = conf.RBMQ_Channel.Publish(
		conf.RBMQ_EXCHANGE_API_SOCKET,
		conf.RBMQ_KEY_ROOM_IO,
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
	util.WriteJson(w, 200, group_member)
}

func HandleRejectGroupInvite(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	req_id := r.PathValue("invite_id")
	req_uuid, err := uuid.Parse(req_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	_, err = conf.DB_Queries.RejectGroupInvite(r.Context(), database.RejectGroupInviteParams{
		ReceiverID: pgtype.UUID{Bytes: user_uuid, Valid: true},
		InviteID:   pgtype.UUID{Bytes: req_uuid, Valid: true},
	})
	if err != nil {
		util.WriteError(w, 400, "Cannot update row in database", err.Error())
		return
	}
	util.WriteMessage(w, 200, "Success")
}

func HandleGetAllInvitesToGroup(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	group_id := r.PathValue("group_id")
	group_uuid, err := uuid.Parse(group_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	group, err := conf.DB_Queries.FindGroupById(r.Context(), pgtype.UUID{Bytes: group_uuid, Valid: true})
	if err != nil {
		util.WriteError(w, 404, "Cannot find from database", err.Error())
		return
	}
	if group.CreatedBy.Bytes != user_uuid {
		util.WriteError(w, 403, "You're not own this group")
		return
	}
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
	group_reqs, err := conf.DB_Queries.GetGroupInvites(r.Context(), database.GetGroupInvitesParams{
		GroupID: pgtype.UUID{Bytes: group_uuid, Valid: true},
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		util.WriteError(w, 404, "Cannot select from database", err.Error())
		return
	}
	util.WriteJson(w, 200, group_reqs)
}

func HandleCreateGroupInvite(w http.ResponseWriter, r *http.Request) {
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
	group_id := r.PathValue("group_id")
	group_uuid, err := uuid.Parse(group_id)
	if err != nil {
		util.WriteError(w, 400, "Cannot parse uuid", err.Error())
		return
	}
	group, err := conf.DB_Queries.FindGroupById(r.Context(), pgtype.UUID{Bytes: group_uuid, Valid: true})
	if err != nil {
		util.WriteError(w, 404, "Cannot select from database", err.Error())
		return
	}
	if group.CreatedBy.Bytes != user_uuid {
		util.WriteError(w, 403, "You're not own this group")
		return
	}
	req, err := conf.DB_Queries.CreateGroupInvite(r.Context(), database.CreateGroupInviteParams{
		GroupID:    pgtype.UUID{Bytes: group_uuid, Valid: true},
		SenderID:   pgtype.UUID{Bytes: user_uuid, Valid: true},
		ReceiverID: pgtype.UUID{Bytes: target_uuid, Valid: true},
		CreatedAt:  pgtype.Timestamp{Time: time.Now(), InfinityModifier: pgtype.Finite, Valid: true},
	})
	if err != nil {
		util.WriteError(w, 400, "Cannot insert into database", err.Error())
		return
	}
	cur_time := time.Now()
	noti, err := conf.DB_Queries.CreateNotification(r.Context(), database.CreateNotificationParams{
		UserID:    pgtype.UUID{Bytes: target_uuid, Valid: true},
		Message:   fmt.Sprintf("%s invite you to %s", user_name, group.GroupName),
		IDRef:     req.InviteID,
		Type:      database.NotificationTypeGroupRequest,
		CreatedAt: pgtype.Timestamp{Time: cur_time, Valid: true, InfinityModifier: pgtype.Finite},
	})
	if err != nil {
		util.WriteError(w, 400, "Cannot insert into database", err.Error())
		return
	}
	notification_msg, _ := json.Marshal(
		conf.QueueMsg_Notification{
			NotificationId: noti.NotificationID.String(),
			UserId:         noti.UserID.String(),
			Message:        noti.Message,
			Type:           conf.REQ_GROUP,
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
		util.WriteError(w, 400, "Cannot publish message to queue", err.Error())
		return
	}
	util.WriteJson(w, 200, req)
}

func HandleDeleteGroupInvite(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	target_id := r.PathValue("target_id")
	target_uuid, err := uuid.Parse(target_id)
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	err = conf.DB_Queries.DeleteGroupInvite(r.Context(), database.DeleteGroupInviteParams{
		SenderID:   pgtype.UUID{Bytes: user_uuid, Valid: true},
		ReceiverID: pgtype.UUID{Bytes: target_uuid, Valid: true},
	})
	if err != nil {
		util.WriteError(w, 400, "Cannot delete row from database", err.Error())
		return
	}
	util.WriteMessage(w, 200, "Success")
}
