# spotify-oauth-proxy

Managed OAuth access for your spotify account

## WARNING

This project is unmaintained. Please go to https://github.com/conradludgate/oauth2-proxy for the newer project

## How it works

1)  Go to https://spotify.conradludgate.com/login
2)  Choose the [scopes](https://developer.spotify.com/documentation/general/guides/scopes/) you wish to use.
3)  Click confirm and login with your Spotify account
4)  Copy and save the API Key. You will only see it once and it's needed to access the API

And you're all set. Now you can make a request to

```
GET https://spotify.conradludgate.com/api/token
Authorization: Basic <base64 encoded user_id:api_key>
```

And it will respond with your account's access token, token type and expiry time.

```json
{"access_token": "XXXXX", "token_type": "Bearer", "expires": 1610115227}
```
