package authmw

// bearerSchema is standart schema jwt header token
// implemented on auth service middleware
const bearerSchema = "Bearer"

// contextClaimKey key value store/get token on context
const contextClaimKey = "ctx.mw.auth.claim"
