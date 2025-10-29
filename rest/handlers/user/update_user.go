package user

import (
	"encoding/json"
	"net/http"

	"ecoscan.com/repo"
)

func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	uID := r.Context().Value("userID")

	// decode the body 
	var reqUpdateUser repo.User

	err := json.NewDecoder(r.Body).Decode(&reqUpdateUser)
	if err != nil{
		http.Error(w, "Invalid update user request body", http.StatusBadRequest)
		return
	}

	reqUpdateUser.ID = uID.(int64)



	



}
