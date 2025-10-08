import { redirect } from '@sveltejs/kit';
import * as querystring from 'querystring';

import { queryAPI } from '../../../util.js';

export async function load({ cookies }) {

    const sessionId = cookies.get("session_id")

    if (!sessionId) {
        throw redirect(307, '/?' +
        querystring.stringify({
            error: 'auth_failure'
        }))
    }

    const nfcResponse = await queryAPI('/api/nfc/unused', sessionId)

    if (nfcResponse.error !== null) {
        console.log(nfcResponse.error)
        throw redirect(307, '/?' +
        querystring.stringify({
            error: 'system_failure'
        }))
    }

    return {
        nfcList: nfcResponse.data,
    }
}

