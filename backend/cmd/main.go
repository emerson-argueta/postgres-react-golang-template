package main

import (
	community_goal_tracker_subscriptions "emersonargueta/m/v1/modules/communitygoaltracker/subscriptions"
	"emersonargueta/m/v1/shared/infrastructure/http"
)

func main() {
	community_goal_tracker_subscriptions.New()

	httpServer := http.NewServer()

	if err := httpServer.Open(); err != nil {
		panic(err)
	}

}
