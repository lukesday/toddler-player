
import { redirect } from '@sveltejs/kit'
import * as querystring from 'querystring'

const redirect_uri = import.meta.env.VITE_SPOTIFY_REDIRECT_URI
const localServerUri = import.meta.env.VITE_SERVER_URI

var stateKey = 'spotify_auth_state';

export async function load({ url, cookies }) {
  var code = url.searchParams.get('code') || null;
  var state = url.searchParams.get('state') || null;
  var storedState = cookies ? cookies.get(stateKey) : null;

  if (state === null || state !== storedState) {
    throw redirect(307, '/?' +
      querystring.stringify({
        error: 'state_mismatch'
      }));
  } else {
    cookies.delete(stateKey)

    const data = {
      code: code,
      redirect: redirect_uri,
      grant_type: 'authorization_code',
    }

    const response = await fetch(`${localServerUri}/api/spotify/auth`, 
    {
      method: "POST",
      mode: "no-cors",
      cache: "no-cache",
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: Object.entries(data).map(e => e.join('=')).join('&'),
    })

    if (response.status !== 200) {
      throw redirect(307, '/?' +
        querystring.stringify({
          error: 'auth_failure'
        }));
    }

    const body = await response.json()

    cookies.set("session_id", body.session_id, {path: '/'})

    throw redirect(307, '/');
  }
}