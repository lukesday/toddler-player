const localServerUri = import.meta.env.VITE_SERVER_URI

export async function load({ cookies }) {

    const sessionId = cookies.get("session_id")

    if (!sessionId) {
      return {
          loggedIn: false,
          userData: null,
      }
    }

    const response = await fetch(`${localServerUri}/api/spotify/me`, 
    {
      method: "GET",
      mode: "no-cors",
      cache: "no-cache",
      headers: {
        'Session-Id': sessionId,
      },
    })

    if (response.status === 200) {
      const body = await response.json()
      cookies.set("session_id", body.session_id, {path: '/'})

      return {
          loggedIn: true,
          userData: body,
      }
    }
    
    return {
        loggedIn: false,
        userData: null,
    }
}