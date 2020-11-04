import { IError } from "./Types"

export const USER_LOADING = "AUTH_USER_LOADING"
export const USER_LOADED = "AUTH_USER_LOADED"
export const USER_REFRESH = "AUTH_USER_REFRESH"
export const AUTH_ERROR = "AUTH_ERROR"

export const LOGIN_URL_POSTFIX = '/achiever/login'
export const LOGIN_SUCCESS = "AUTH_LOGIN_SUCCESS"
export const LOGIN_FAIL = "AUTH_LOGIN_FAIL"
export const LOGOUT = "AUTH_LOGOUT"
export const LOGOUT_SUCCESS = "AUTH_LOGOUT_SUCCESS"

export const REGISTER_URL_POSTFIX = '/achiever'
export const REGISTER_SUCCESS = "AUTH_REGISTER_SUCCESS"
export const REGISTER_FAIL = "AUTH_REGISTER_FAIL"

export const CLEAR_ERROR = "AUTH_CLEAR_ERROR"

export const REAUTHORIZE_URL_POSTFIX = '/achiever/reauthorize'
export const TRUSTDONATIONS_ACCESS_TOKEN = "AUTH_TRUSTDONATIONS_ACCESS_TOKEN"
export const TRUSTDONATIONS_REFRESH_TOKEN = "AUTH_TRUSTDONATIONS_REFRESH_TOKEN"


export interface IAuthState {
    isAuthenticated: boolean
    token: any
    error: IError | null
    loading: boolean
}
export interface IAuthorization {
    accesstoken: string | null
    refreshtoken: string | null
}