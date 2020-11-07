import axios, { AxiosResponse } from 'axios'
import * as AUTH_TYPES from '../../types/AuthTypes'
import * as TYPES from '../../types/Types'
import { IAchiever, TAchieverAPIRequest, IAchieverAPIResponse } from '../../types/AchieverTypes'
import { Dispatch } from 'react'


export const userLoginACT = (achiever: IAchiever) => async (dispatch: Dispatch<AUTH_TYPES.TAuthActions>) => {
    const url = TYPES.API_URL_PREFIX + AUTH_TYPES.LOGIN_URL_POSTFIX
    const req: TAchieverAPIRequest = { achiever: achiever }

    try {
        const res: AxiosResponse<IAchieverAPIResponse> = await axios.post(url, req)
        dispatch({ type: AUTH_TYPES.CLEAR_ERROR });
        dispatch({ type: AUTH_TYPES.LOGIN_SUCCESS, payload: res.data });
    } catch (err) {
        const error: TYPES.IError = { id: AUTH_TYPES.LOGIN_FAIL, status: err?.response.status, msg: err?.response?.data?.error }
        dispatch({ type: AUTH_TYPES.LOGIN_FAIL, error: error })
    }

}

export const userLogoutACT = () => async (dispatch: Dispatch<AUTH_TYPES.TAuthActions>) => {
    dispatch({ type: AUTH_TYPES.LOGOUT })
}

export const userRegisterACT = (achiever: IAchiever) => async (dispatch: Dispatch<AUTH_TYPES.TAuthActions>) => {
    const url = TYPES.API_URL_PREFIX + AUTH_TYPES.REGISTER_URL_POSTFIX
    const req: TAchieverAPIRequest = { achiever: achiever }

    try {
        const res: AxiosResponse<IAchieverAPIResponse> = await axios.post(url, req)
        dispatch({ type: AUTH_TYPES.CLEAR_ERROR })
        dispatch({ type: AUTH_TYPES.REGISTER_SUCCESS, payload: res.data })

    } catch (err) {
        const error: TYPES.IError = { id: AUTH_TYPES.REGISTER_FAIL, status: err?.response.status, msg: err?.response?.data?.error }
        dispatch({ type: AUTH_TYPES.REGISTER_FAIL, error: error })
    }
}

export const userReAuthorizeACT = (retryAction?: Function) => async (dispatch: Dispatch<AUTH_TYPES.TAuthActions | Function>) => {
    const url = TYPES.API_URL_PREFIX + AUTH_TYPES.REAUTHORIZE_URL_POSTFIX

    const token: AUTH_TYPES.IAuthorization = {
        accesstoken: localStorage.getItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN) || undefined,
        refreshtoken: localStorage.getItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_REFRESH_TOKEN) || undefined
    }

    const req: TAchieverAPIRequest = { authorization: token }
    try {
        const res: AxiosResponse<IAchieverAPIResponse> = await axios.post(url, req)
        dispatch({ type: AUTH_TYPES.REAUTHORIZE_SUCCESS, payload: res.data });
        if (retryAction) {
            dispatch(retryAction)
        }
    } catch (err) {
        const error: TYPES.IError = { id: AUTH_TYPES.REAUTHORIZE_FAIL, status: err?.response.status, msg: err?.response?.data?.error }
        dispatch({ type: AUTH_TYPES.REAUTHORIZE_FAIL, error: error })
    }
}

export const userLoadingACT = () => async (dispatch: Dispatch<AUTH_TYPES.TAuthActions>) => {
    dispatch({ type: AUTH_TYPES.USER_LOADING });
};
export const userLoadedACT = () => async (dispatch: Dispatch<AUTH_TYPES.TAuthActions>) => {
    dispatch({ type: AUTH_TYPES.USER_LOADED });
};



// Setup config/headers and token
export const tokenConfig = () => {
    const accessToken = localStorage.getItem(AUTH_TYPES.COMMUNITY_GOAL_TRACKER_ACCESS_TOKEN)

    const config: TYPES.IConfigHeaders = {
        headers: {
            'Content-type': 'application/json'
        }
    };

    if (accessToken) {
        config.headers['Authorization'] = "Bearer " + accessToken;
    }

    return config;
};