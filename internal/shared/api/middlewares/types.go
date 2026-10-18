package middlewares

// They serve as safe keys to store and retrieve token claims in the context.Context.

// This way, the middlewares can validate the token, and then the handlers can use those claims without parsing the JWT again.

// Access token context key

type accessTokenClaimsContext string

const AccessTokenClaimsIDKey accessTokenClaimsContext = "access-token-claims"

// Refresh token context key

type refreshTokenClaimsContext string

const RefreshTokenClaimsKey refreshTokenClaimsContext = "refresh-token-claims"
