import { TAchieverAPIResponse } from "./AchieverTypes"
import { IError } from "./Types"

export const USER_LOADING = "AUTH_USER_LOADING"
export const USER_LOADED = "AUTH_USER_LOADED"
export const USER_REFRESH = "AUTH_USER_REFRESH"
interface IUserAction {
    type: typeof USER_LOADING | typeof USER_LOADED | typeof USER_REFRESH
    payload?: TAchieverAPIResponse
}

export const AUTH_ERROR = "AUTH_ERROR"
interface IErrorAction {
    type: typeof AUTH_ERROR
    payload?: IError
}

export const LOGIN_URL_POSTFIX = '/achiever/login'
export const LOGIN_SUCCESS = "AUTH_LOGIN_SUCCESS"
export const LOGIN_FAIL = "AUTH_LOGIN_FAIL"
interface ILoginAction {
    type: typeof LOGIN_SUCCESS | typeof LOGIN_FAIL
    payload: TAchieverAPIResponse
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
    payload: TAchieverAPIResponse
}

export const CLEAR_ERROR = "AUTH_CLEAR_ERROR"
interface IClearErrorAction {
    type: typeof CLEAR_ERROR
}

export const REAUTHORIZE_URL_POSTFIX = '/achiever/reauthorize'
export const COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN = "AUTH_COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN"
export const COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN = "AUTH_COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN"


export interface IAuthState {
    isAuthenticated?: boolean
    token?: any
    error: IError | null
    loading: boolean
}
export interface IAuthorization {
    accesstoken?: string
    refreshtoken?: string
}

export type TAuthActions = IUserAction | ILoginAction | ILogoutAction | IRegisterAction | IClearErrorAction | IErrorAction
