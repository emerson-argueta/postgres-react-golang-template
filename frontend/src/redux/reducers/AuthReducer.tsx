import { IAchieverAPIResponse } from "../../types/AchieverTypes"
import * as AUTH_TYPES from "../../types/AuthTypes"
import * as TYPES from "../../types/Types"

const initialState: AUTH_TYPES.IAuthState = {

    authorization: {
        accesstoken: localStorage.getItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN) || undefined,
        refreshtoken: localStorage.getItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN) || undefined
    },
    isAuthenticated: false,
    error: null,
    loading: true
}

export default (state = initialState, action: AUTH_TYPES.TAuthActions) => {

    switch (action.type) {
        case AUTH_TYPES.REGISTER_SUCCESS:
        case AUTH_TYPES.LOGIN_SUCCESS: {

            localStorage.setItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN, action.payload?.authorization?.accesstoken || "")
            localStorage.setItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN, action.payload?.authorization?.refreshtoken || "")

            return {
                ...state,
                authorization: action.payload?.authorization,
                isAuthenticated: true,
                loading: false
            }
        }
        case AUTH_TYPES.REAUTHORIZE_SUCCESS: {

            localStorage.setItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN, action.payload?.authorization?.accesstoken || "")
            localStorage.setItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN, action.payload?.authorization?.refreshtoken || "")

            return {
                ...state,
                authorization: action.payload?.authorization,
                isAuthenticated: true,
                loading: false
            }
        }
        case AUTH_TYPES.USER_LOADED: {
            return {
                ...state,
                isAuthenticated: true,
                loading: false
            }
        }
        case AUTH_TYPES.REGISTER_FAIL:
        case AUTH_TYPES.LOGIN_FAIL:
        case AUTH_TYPES.REAUTHORIZE_FAIL: {
            localStorage.removeItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN)
            localStorage.removeItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN)

            return {
                ...state,
                isAuthenticated: false,
                token: null,
                error: action.payload
            }
        }
        case AUTH_TYPES.LOGOUT: {
            localStorage.removeItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN)
            localStorage.removeItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN)
            return {
                ...state,
                administrator: null,
                token: null,
                isAuthenticated: false
            }
        }
        case AUTH_TYPES.CLEAR_ERROR: {
            return {
                ...state,
                error: null
            }
        }
        case AUTH_TYPES.USER_LOADING: {
            return {
                ...state,
                loading: true
            }
        }
        default:
            return state
    }
}
