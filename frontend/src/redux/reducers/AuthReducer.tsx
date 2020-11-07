import { TAchieverAPIResponse } from "../../types/AchieverTypes"
import * as AUTH_TYPES from "../../types/AuthTypes"
import * as TYPES from "../../types/Types"

const initialState: AUTH_TYPES.IAuthState = {

    token: {
        Accesstoken: localStorage.getItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN),
        Refreshtoken: localStorage.getItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN)
    },
    isAuthenticated: false,
    error: null,
    loading: true
}

export default (state = initialState, action: TYPES.IAction) => {

    switch (action.type) {
        case AUTH_TYPES.REGISTER_SUCCESS:
        case AUTH_TYPES.LOGIN_SUCCESS: {
            const achieverResponse: TAchieverAPIResponse = action.payload

            localStorage.setItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN, achieverResponse.authorization?.accesstoken || "")
            localStorage.setItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN, achieverResponse.authorization?.refreshtoken || "")

            return {
                ...state,
                token: action.payload.token,
                isAuthenticated: true,
                loading: false
            }
        }
        case AUTH_TYPES.USER_REFRESH: {
            const achieverResponse: TAchieverAPIResponse = action.payload

            localStorage.setItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN, achieverResponse.authorization?.accesstoken || "")
            localStorage.setItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN, achieverResponse.authorization?.refreshtoken || "")

            return {
                ...state,
                token: action.payload.token,
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
        case AUTH_TYPES.AUTH_ERROR: {
            localStorage.removeItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN)
            localStorage.removeItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN)

            const error: TYPES.IError = { msg: action.payload?.data.error, id: action.payload?.id, status: action.payload?.status }

            return {
                ...state,
                isAuthenticated: false,
                token: null,
                error: error
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
