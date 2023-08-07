
import { redirect } from '@sveltejs/kit'
import * as querystring from 'querystring'

var client_id = import.meta.env.VITE_SPOTIFY_CLIENT_ID
var client_secret = import.meta.env.VITE_SPOTIFY_CLIENT_SECRET
var redirect_uri = import.meta.env.VITE_SPOTIFY_REDIRECT_URI

var stateKey = 'spotify_auth_state';

export async function load({ url, cookies }) {
  var code = url.searchParams.get('code') || null;
  var state = url.searchParams.get('state') || null;
  var storedState = cookies ? cookies.get(stateKey) : null;

  if (state === null || state !== storedState) {
    throw redirect(307, '/' +
      querystring.stringify({
        error: 'state_mismatch'
      }));
  } else {
    cookies.delete(stateKey)

    const data = {
      code: code,
      redirect_uri: redirect_uri,
      grant_type: 'authorization_code'
    }

    // This should be moved to the Server!
    const response = await fetch('https://accounts.spotify.com/api/token', 
    {
      method: "POST",
      mode: "no-cors",
      cache: "no-cache",
      headers: {
        'Authorization': 'Basic ' + btoa(client_id + ':' + client_secret),
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: Object.entries(data).map(e => e.join('=')).join('&'),
    })

    if (response.status === 200) {
      const body = await response.json()

      var access_token = body.access_token,
          refresh_token = body.refresh_token;

      const meResponse = await fetch('https://api.spotify.com/v1/me', {
        headers: { 'Authorization': 'Bearer ' + access_token },
      })

      console.log(await meResponse.json())

      throw redirect (307, '/?' +
        querystring.stringify({
          access_token: access_token,
          refresh_token: refresh_token
        }));

    } else {
      throw redirect (307, '/?' +
        querystring.stringify({
          error: 'invalid_token'
        }));
    }
  }
}