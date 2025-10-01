package handler

import (
	"net/http"

	"github.com/anhtr13/synth-socket/api-service/api/guard"
	"github.com/anhtr13/synth-socket/api-service/api/middleware"
)

func RegisterRoutes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /auth/signup", HandleSignUp)
	router.HandleFunc("POST /auth/signin", HandleSignIn)
	router.HandleFunc("POST /auth/signout", guard.Auth(HandleSignOut))

	router.HandleFunc("GET /@/access_token", HandleGetAccessToken) // ?refresh_token=<refresh_token>
	router.HandleFunc("GET /@/info", guard.Auth(HandleGetMyInfo))
	router.HandleFunc("PATCH /@/info", guard.Auth(HandleUpdateMyInfo))

	router.HandleFunc("GET /@/friend", guard.Auth(HandleGetMyFriends)) // ?search=friend_name
	// router.HandleFunc("DELETE /@/friend/{friend_id}", guard.Auth(HandleDeleteFriend))
	router.HandleFunc("GET /@/friend_request", guard.Auth(HandleGetMyFriendRequests))
	router.HandleFunc("POST /@/friend_request/{request_id}", guard.Auth(HandleAcceptFriendRequest))
	router.HandleFunc("DELETE /@/friend_request/{request_id}", guard.Auth(HandleRejectFriendRequest))

	router.HandleFunc("GET /@/group", guard.Auth(HandleGetMyGroups)) // ?search=group_name
	router.HandleFunc("DELETE /@/group/{group_id}", guard.Auth(HandleLeaveGroup))
	router.HandleFunc("GET /@/group_invite", guard.Auth(HandleGetMyGroupInvites))
	router.HandleFunc("POST /@/group_invite/{invite_id}", guard.Auth(HandleAcceptGroupInvite))
	router.HandleFunc("DELETE /@/group_invite/{invite_id}", guard.Auth(HandleRejectGroupInvite))

	router.HandleFunc("GET /@/notification", guard.Auth(HandlerGetMyNotifications))
	router.HandleFunc("POST /@/notification/{notification_id}", guard.Auth(HandlerSeenNotification))

	router.HandleFunc("POST /group", guard.Auth(HandleCreateGroup))
	router.HandleFunc("GET /group/{group_id}/member", guard.Auth(HandleGetGroupMembers))
	router.HandleFunc("GET /group/{group_id}/invite", guard.Auth(HandleGetAllInvitesToGroup))
	router.HandleFunc("POST /group/{group_id}/invite/{target_id}", guard.Auth(HandleCreateGroupInvite))
	router.HandleFunc("DELETE /group/{group_id}/invite/{target_id}", guard.Auth(HandleDeleteGroupInvite))

	router.HandleFunc("GET /user", guard.Auth(HandleSearchUserByName)) // ?search=user_name
	router.HandleFunc("GET /user/{user_id}", guard.Auth(HandleGetUserInfo))
	router.HandleFunc("POST /user/{user_id}/friend_request", guard.Auth(HandleCreateFriendRequest))
	router.HandleFunc("DELETE /user/{user_id}/friend_request", guard.Auth(HandleDeleteFriendRequest))

	api_v1 := http.NewServeMux()
	api_v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	middlewareStack := middleware.CreateMdwStack(middleware.Auth)

	return middlewareStack(api_v1)
}
