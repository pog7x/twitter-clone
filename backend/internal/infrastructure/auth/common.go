package auth

type (
	loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	sessionClaims struct {
		SessionID string `json:"session_id"`
	}
)
