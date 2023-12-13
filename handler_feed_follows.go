package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	// "github.com/Bayan2019/rss_blog/internal/auth"
	"github.com/Bayan2019/rss_blog/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed : %s", err))
	}

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

// func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {

// 	feeds, err := apiCfg.DB.GetFeeds(r.Context())
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds : %s", err))
// 	}

// 	respondWithJSON(w, 201, databaseFeedsToFeeds(feeds))
// }