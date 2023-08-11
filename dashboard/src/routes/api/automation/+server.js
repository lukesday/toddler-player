import { error, json } from '@sveltejs/kit'
import * as querystring from 'querystring'

const localServerUri = import.meta.env.VITE_SERVER_URI

export async function GET({ cookies }) {

    const sessionId = cookies.get("session_id")

    if (!sessionId) {
      throw error(400, 'Bad Request')
    }

    const response = await fetch(`${localServerUri}/api/automation`, 
    {
      method: "GET",
      mode: "no-cors",
      cache: "no-cache",
      headers: {
        'Session-Id': sessionId,
      },
    })

    if (response.status !== 200) {
        throw error(500, 'Server Error')
    }
    
    const body = await response.json()

    return json({
        automations: body
    })
}