const localServerUri = import.meta.env.VITE_SERVER_URI

const queryAPI = async (resource, sessionId) => {
    
    const response = await fetch(`${localServerUri}${resource}`, 
    {
      method: "GET",
      mode: "no-cors",
      cache: "no-cache",
      headers: {
        'Session-Id': sessionId,
      },
    })
    
    if (response.status !== 200) {
        return { 
            error: response.error,
            data: null
        }
    }
    
    return {
        error: null,
        data: response.json()
    }
}

export { queryAPI }