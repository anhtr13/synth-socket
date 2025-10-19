package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/anhtr13/synth-socket/api/internal/conf"
	"github.com/anhtr13/synth-socket/api/internal/database"
	"github.com/anhtr13/synth-socket/api/internal/util"
)

func HandlerGetNotifications(w http.ResponseWriter, r *http.Request) {
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
	notis, err := conf.DB_Queries.GetUnSeenNotificationsByUserId(r.Context(), database.GetUnSeenNotificationsByUserIdParams{
		UserID: pgtype.UUID{Bytes: user_uuid, Valid: true},
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		util.WriteError(w, 404, "Cannot select from database", err.Error())
		return
	}
	util.WriteJson(w, 200, notis)
}

func HandlerSeenNotification(w http.ResponseWriter, r *http.Request) {
	user_session := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
	user_id := user_session["id"].(string)
	user_uuid, _ := uuid.Parse(user_id)
	notification_id := r.PathValue("notification_id")
	notification_uuid, err := uuid.Parse(notification_id)
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	err = conf.DB_Queries.SeenNotification(r.Context(), database.SeenNotificationParams{
		NotificationID: pgtype.UUID{Bytes: notification_uuid, Valid: true},
		UserID:         pgtype.UUID{Bytes: user_uuid, Valid: true},
	})
	if err != nil {
		util.WriteError(w, 400, "Cannot update from database", err.Error())
		return
	}
	util.WriteMessage(w, 200, "Success")
}
