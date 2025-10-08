import { json, error } from '@sveltejs/kit';
 
const localServerUri = import.meta.env.VITE_SERVER_URI

export async function POST({ request, cookies }) {

    const sessionId = cookies.get("session_id")

    if (!sessionId) {
        return error(400, "invalid session id")
    }

    const requestBody = await request.json()

    console.log(requestBody)

    const response = await fetch(`${localServerUri}/api/automation/${requestBody.id}`, 
    {
      method: "Delete",
      cache: "no-cache",
      headers: {
        'Session-Id': sessionId,
        'Content-Type': 'application/json'
      },
    })
    
    if (response.status !== 200) {
        console.log(response.status)
        return json({ success: false }, {status: response.status})
    }

    return json({ success: true }, {status: 200})
}