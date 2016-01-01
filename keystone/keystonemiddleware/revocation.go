package keystonemiddleware

// Fetch the token revocation list from identity server
// So that even token is not expired locally, but it was revoked
// from identity server, we will know it and do not let it
// pass the validation.

// One way to implement this is to use a channel for revoked token list
// and use a go-routine for checking token revokes

// Currently we leave it non-implemented.