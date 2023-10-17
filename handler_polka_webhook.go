package main

import (
	"encoding/json"
	"net/http"

	"github.com/jhonipereira/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	type response struct {
		User
	}

	token, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find API KEY")
		return
	}
	if auth.ValidateAPIKEY(token, cfg.polkaKey) == false {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate API KEY")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 200, "")
		return
	}

	user, err := cfg.DB.GetUserByID(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't found the user")
		return
	}

	_, err = cfg.DB.UpdateUserToPremium(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't updated user")
		return
	}

	respondWithJSON(w, http.StatusOK, "")
}
