package auth

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type sessionClaims struct {
	SessionID string `json:"session_id"`
}
