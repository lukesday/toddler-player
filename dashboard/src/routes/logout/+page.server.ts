
import { redirect } from '@sveltejs/kit'

export function load({ cookies }) {
  cookies.set("session_id", "", { path: "/", maxAge: 0}) 
}