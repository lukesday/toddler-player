import { redirect } from '@sveltejs/kit'
import * as querystring from 'querystring'

const localServerUri = import.meta.env.VITE_SERVER_URI

export async function load() {
    const response = await fetch(`${localServerUri}/api/spotify/me`, 
    {
      method: "GET",
      mode: "no-cors",
      cache: "no-cache"
    })

    return {
        loggedIn: response.status === 200,
        userData: response.status === 200 ? await response.json() : null,
    }
}