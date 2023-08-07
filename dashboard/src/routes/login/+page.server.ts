
import { redirect } from '@sveltejs/kit'
import * as querystring from 'querystring'

var client_id = import.meta.env.VITE_SPOTIFY_CLIENT_ID
var client_secret = import.meta.env.VITE_SPOTIFY_CLIENT_SECRET
var redirect_uri = import.meta.env.VITE_SPOTIFY_REDIRECT_URI

console.log(import.meta.env)

console.log(client_id, client_secret, redirect_uri)

var generateRandomString = function(length : number) {
  var text = '';
  var possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';

  for (var i = 0; i < length; i++) {
    text += possible.charAt(Math.floor(Math.random() * possible.length));
  }
  return text;
};

export function load() {

    var state = generateRandomString(16);
    var scope = 'user-read-private user-read-email';

    throw redirect (307, 'https://accounts.spotify.com/authorize?' +
    querystring.stringify({
      response_type: 'code',
      client_id: client_id,
      scope: scope,
      redirect_uri: redirect_uri,
      state: state
    }))
}