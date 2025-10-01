package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/anhtr13/synth-socket/api-service/api/handler"
	"github.com/anhtr13/synth-socket/api-service/internal/conf"
)

func init() {
	conf.InitConnection()
}

// load all user's friendships, user's groups to redis
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
	group_members_table, err := conf.DB_Queries.GetGroupMembersTable(ctx)
	if err != nil {
		return err
	}
	user_groups := map[string][]string{}
	for _, gm := range group_members_table {
		user_id := gm.MemberID.String()
		if user_groups[user_id] == nil {
			user_groups[user_id] = []string{}
		}
		user_groups[user_id] = append(user_groups[user_id], gm.GroupID.String())
	}
	for user_id, groups := range user_groups {
		conf.RD_Client.SAdd(
			ctx,
			fmt.Sprintf("%s:%s:%s", conf.RDKEY_PREFIX_USER, user_id, conf.RDKEY_SUFFIX_GROUPS),
			groups,
		).Err()
	}

	friendship_table, err := conf.DB_Queries.GetFriendshipTable(ctx)
	if err != nil {
		return err
	}
	user_friendships := map[string][]string{}
	for _, ft := range friendship_table {
		user1_id := ft.User1ID.String()
		user2_id := ft.User2ID.String()
		friendship_id := ft.FriendshipID.String()
		if user_friendships[user1_id] == nil {
			user_friendships[user1_id] = []string{}
		}
		if user_friendships[user2_id] == nil {
			user_friendships[user2_id] = []string{}
		}
		user_friendships[user1_id] = append(user_friendships[user1_id], friendship_id)
		user_friendships[user2_id] = append(user_friendships[user2_id], friendship_id)
	}
	for user_id, friendships := range user_friendships {
		conf.RD_Client.SAdd(
			ctx,
			fmt.Sprintf("%s:%s:%s", conf.RDKEY_PREFIX_USER, user_id, conf.RDKEY_SUFFIX_FRIENDSHIPS),
			friendships,
		).Err()
	}

	return nil
}

// Graceful shutdown to ensures that ongoing requests are completed and resources are properly released before the server terminates
func gracefulShutdown(apiServer *http.Server, done chan bool) {
	defer conf.CloseAllConnections()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func main() {
	err := loadDataToCache()
	if err != nil {
		log.Fatal("Cannot load data to cache:", err)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.PORT),
		Handler:      handler.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	forever := make(chan bool, 1)
	go gracefulShutdown(server, forever)

	log.Println("Server is running on port", server.Addr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("Http server error: %s", err))
	}

	<-forever

	log.Println("Graceful shutdown complete.")
}
