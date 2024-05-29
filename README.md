# Globus Go Library
## Description
This library should provide a way for an app to request and monitor transfers from Globus after authentication.

For authentication with Globus, the standard OAuth2 implementation for Go is used (github.com/golang/oauth2).

## CLI app
The `cmd/` subfolder contains a full implementation of all capabilities of this library in the form of a command line application.

The `client credential / code grant` based authentication requires the user to authenticate each time, as passing a refresh/auth token is not supported at this time. The library itself should be capable of doing this eventually.   
