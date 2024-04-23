package middlewares

type loginRequest struct {
	Username string `json:"username"`
	Passwrod string `json:"password"`
}

type sessionClaims struct {
	SessionID string `json:"session_id"`
}
