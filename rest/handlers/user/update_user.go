package user

import (
	"encoding/json"
	"net/http"

	"ecoscan.com/repo"
)

func (h *UserHandler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	uID := r.Context().Value("userID")

	// decode the body 
	var reqUpdateUser repo.User

	err := json.NewDecoder(r.Body).Decode(&reqUpdateUser)
	if err != nil{
		http.Error(w, "Invalid update user request body", http.StatusBadRequest)
		return
	}

	reqUpdateUser.ID = uID.(int64)

	query := `
    UPDATE users 
    SET name = $1, email = $2, updated_at = NOW()
    WHERE id = $3
    RETURNING id, name, email, points, created_at, updated_at
	`
	row := h.DB.QueryRowContext(r.Context(), query, reqUpdateUser.Name,
	reqUpdateUser.Email,
	reqUpdateUser.ID,
	)

	var updatedUser repo.User
	err = row.Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
        &updatedUser.Points,
        &updatedUser.CreatedAt,
        &updatedUser.UpdatedAt,
	)
	if err != nil{
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)

}
