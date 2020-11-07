import { Reducer } from "redux"
import * as AUTH_TYPES from "../../types/AuthTypes"

const initialState: AUTH_TYPES.IAuthState = {
    authorization: {
        accesstoken: localStorage.getItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN) || undefined,
        refreshtoken: localStorage.getItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN) || undefined
    },
    isAuthenticated: false,
    error: undefined,
    loading: true
}

export const AuthReducer: Reducer<AUTH_TYPES.IAuthState, AUTH_TYPES.TAuthActions> = (state = initialState, action): AUTH_TYPES.IAuthState => {

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
                authorization: undefined,
                error: action.error
            }
        }
        case AUTH_TYPES.LOGOUT: {
            localStorage.removeItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN)
            localStorage.removeItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN)
            return {
                ...state,
                authorization: undefined,
                isAuthenticated: false
            }
        }
        case AUTH_TYPES.CLEAR_ERROR: {
            return {
                ...state,
                error: undefined
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
