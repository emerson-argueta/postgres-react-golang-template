import { IAchieverAPIResponse } from "./AchieverTypes"
import { IError } from "./Types"

export const USER_LOADING = "AUTH_USER_LOADING"
export const USER_LOADED = "AUTH_USER_LOADED"
interface IUserAction {
    type: typeof USER_LOADING | typeof USER_LOADED
    payload?: IAchieverAPIResponse
}

export const LOGIN_URL_POSTFIX = '/achiever/login'
export const LOGIN_SUCCESS = "AUTH_LOGIN_SUCCESS"
export const LOGIN_FAIL = "AUTH_LOGIN_FAIL"
interface ILoginAction {
    type: typeof LOGIN_SUCCESS | typeof LOGIN_FAIL
    payload?: IAchieverAPIResponse
    error?: IError
}

export const LOGOUT = "AUTH_LOGOUT"
interface ILogoutAction {
    type: typeof LOGOUT
}

export const REGISTER_URL_POSTFIX = '/achiever'
export const REGISTER_SUCCESS = "AUTH_REGISTER_SUCCESS"
export const REGISTER_FAIL = "AUTH_REGISTER_FAIL"
interface IRegisterAction {
    type: typeof REGISTER_SUCCESS | typeof REGISTER_FAIL
    payload?: IAchieverAPIResponse
    error?: IError
}

export const CLEAR_ERROR = "AUTH_CLEAR_ERROR"
interface IClearErrorAction {
    type: typeof CLEAR_ERROR
}

export const REAUTHORIZE_URL_POSTFIX = '/achiever/reauthorize'
export const REAUTHORIZE_SUCCESS = "AUTH_REAUTHORIZE_SUCCESS"
export const REAUTHORIZE_FAIL = "AUTH_REAUTHORIZE_FAIL"
interface IErrorAction {
    type: typeof REAUTHORIZE_SUCCESS | typeof REAUTHORIZE_FAIL
    payload?: IAchieverAPIResponse
    error?: IError
}

export const COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN = "AUTH_COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN"
export const COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN = "AUTH_COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN"


export interface IAuthState {
    isAuthenticated?: boolean
    authorization?: IAuthorization
    error: IError | null
    loading: boolean
}
export interface IAuthorization {
    accesstoken?: string
    refreshtoken?: string
}

export type TAuthActions = IUserAction | ILoginAction | ILogoutAction | IRegisterAction | IClearErrorAction | IErrorAction
