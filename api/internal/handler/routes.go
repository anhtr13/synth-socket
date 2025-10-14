package handler

import (
	"net/http"

	"github.com/anhtr13/synth-socket/api/internal/guard"
	"github.com/anhtr13/synth-socket/api/internal/middleware"
)

func RegisterRoutes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /auth/signup", HandleSignUp)
	router.HandleFunc("POST /auth/login", HandleLogin)
	router.HandleFunc("POST /auth/access_token", HandleCreateAccessToken) // ?refresh_token=refresh_token
	router.HandleFunc("POST /auth/logout", guard.Auth(HandleLogOut))

	router.HandleFunc("GET /me/info", guard.Auth(HandleGetPersonalInfo))
	router.HandleFunc("PATCH /me/info", guard.Auth(HandleUpdatePersonalInfo))

	router.HandleFunc("GET /friend", guard.Auth(HandleGetFriends)) // ?search=friend_name
	// router.HandleFunc("DELETE /friend/{friend_id}", guard.Auth(HandleDeleteFriend))
	router.HandleFunc("GET /friend_request", guard.Auth(HandleGetFriendRequests))
	router.HandleFunc("POST /friend_request/{request_id}", guard.Auth(HandleAcceptFriendRequest))
	router.HandleFunc("DELETE /friend_request/{request_id}", guard.Auth(HandleRejectFriendRequest))

	router.HandleFunc("POST /room", guard.Auth(HandleCreateRoom))
	router.HandleFunc("GET /room/all", guard.Auth(HandleGetAllRooms)) // ?search=room_name
	router.HandleFunc("DELETE /room/all/{room_id}", guard.Auth(HandleLeaveRoom))
	router.HandleFunc("GET /room/all/{room_id}/member", guard.Auth(HandleGetRoomMembers))
	router.HandleFunc("GET /room/all/{room_id}/message", guard.Auth(HandleGetRoomMessages))
	router.HandleFunc("GET /room/owned", guard.Auth(HandleGetOwnedRooms))
	router.HandleFunc("GET /room/owned/{room_id}/invite", guard.Auth(HandleGetAllInvitesToRoom))
	router.HandleFunc("POST /room/owned/{room_id}/invite/{target_id}", guard.Auth(HandleCreateRoomInvite))
	router.HandleFunc("DELETE /room/owned/{room_id}/invite/{target_id}", guard.Auth(HandleDeleteRoomInvite))

	router.HandleFunc("GET /room_invite", guard.Auth(HandleGetRoomInvites))
	router.HandleFunc("POST /room_invite/{invite_id}", guard.Auth(HandleAcceptRoomInvite))
	router.HandleFunc("DELETE /room_invite/{invite_id}", guard.Auth(HandleRejectRoomInvite))

	router.HandleFunc("GET /notification", guard.Auth(HandlerGetNotifications))
	router.HandleFunc("POST /notification/{notification_id}", guard.Auth(HandlerSeenNotification))

	router.HandleFunc("GET /user", guard.Auth(HandleGetAllUsersInfo)) // ?search=user_name
	router.HandleFunc("GET /user/{user_id}", guard.Auth(HandleGetUserInfo))
	router.HandleFunc("POST /user/{user_id}/friend_request", guard.Auth(HandleCreateFriendRequest))
	router.HandleFunc("DELETE /user/{user_id}/friend_request", guard.Auth(HandleDeleteFriendRequest))

	api_v1 := http.NewServeMux()
	api_v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	middlewareStack := middleware.CreateMdwStack(middleware.Auth)

	return middlewareStack(api_v1)
}
