package domain


type TokenConfiguration struct {
	// Keycloak client ID
	ClientID string `mapstructure:"client_id"`
	
	// Keycloak client secret
	ClientSecret string `mapstructure:"client_secret"`
	
	// Audience for WebSocket tokens
	WebSocketAudience string `mapstructure:"websocket_audience"`
	
	// Token lifetime in seconds
	TokenLifetime int `mapstructure:"token_lifetime"`
}