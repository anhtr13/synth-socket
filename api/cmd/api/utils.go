package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/anhtr13/synth-socket/api/internal/cache"
	"github.com/anhtr13/synth-socket/api/internal/conf"
)

// load all user's friendships, user's rooms to redis
func loadDataToCache() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	foo, err := conf.RD_Client.Get(
		ctx,
		"foo",
	).Result()
	if err == nil || foo == "bar" {
		return nil
	}

	conf.RD_Client.Set(
		ctx,
		"foo",
		"bar",
		0,
	).Err()

	room_members_table, err := conf.DB_Queries.GetRoomMembersTable(ctx)
	if err != nil {
		return err
	}
	user_rooms := map[string][]string{}
	for _, rm := range room_members_table {
		user_id := rm.MemberID.String()
		if user_rooms[user_id] == nil {
			user_rooms[user_id] = []string{}
		}
		user_rooms[user_id] = append(user_rooms[user_id], rm.RoomID.String())
	}
	for user_id, rooms := range user_rooms {
		conf.RD_Client.SAdd(
			ctx,
			fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user_id, cache.USER_ROOMS),
			rooms,
		).Err()
	}

	friendship_table, err := conf.DB_Queries.GetFriendshipTable(ctx)
	if err != nil {
		return err
	}
	user_friends := map[string][]string{}
	for _, ft := range friendship_table {
		user1_id := ft.User1ID.String()
		user2_id := ft.User2ID.String()
		if user_friends[user1_id] == nil {
			user_friends[user1_id] = []string{}
		}
		if user_friends[user2_id] == nil {
			user_friends[user2_id] = []string{}
		}
		user_friends[user1_id] = append(user_friends[user1_id], user2_id)
		user_friends[user2_id] = append(user_friends[user2_id], user1_id)
	}
	for user_id, friendships := range user_friends {
		conf.RD_Client.SAdd(
			ctx,
			fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user_id, cache.USER_FRIENDS),
			friendships,
		).Err()
	}

	return nil
}

// Graceful shutdown to ensures that ongoing requests are completed and resources are properly released before the server terminates
func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}
