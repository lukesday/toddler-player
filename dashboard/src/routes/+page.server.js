
import { queryAPI } from '../util.js'

export async function load({ cookies }) {

    const sessionId = cookies.get("session_id")

    if (sessionId) {
        const response = await queryAPI('/api/automations', sessionId)
    
        if (response.error === null) {
            return {
                automationList: response.data,
                subTitle: "Logged In",
            }
        }
    }

    return {
        subTitle: "Logged Out",
    }
}

