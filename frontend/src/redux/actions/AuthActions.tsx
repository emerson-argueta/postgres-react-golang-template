import axios, { AxiosResponse } from 'axios'
import * as AUTH_TYPES from '../../types/AuthTypes'
import * as TYPES from '../../types/Types'
import { IAchiever, TAchieverAPIRequest, TAchieverAPIResponse } from '../../types/AchieverTypes'
import { Dispatch } from 'react'


export const userLoginACT = (achiever: IAchiever) => async (dispatch: Dispatch<AUTH_TYPES.TAuthActions>) => {
    const url = TYPES.API_URL_PREFIX + AUTH_TYPES.LOGIN_URL_POSTFIX
    const req: TAchieverAPIRequest = { achiever: achiever }

    try {
        const res: AxiosResponse<TAchieverAPIResponse> = await axios.post(url, req)
        dispatch({ type: AUTH_TYPES.CLEAR_ERROR });
        dispatch({ type: AUTH_TYPES.LOGIN_SUCCESS, payload: res.data });
    } catch (err) {
        dispatch({ type: AUTH_TYPES.LOGIN_FAIL, payload: { ...err?.response, id: AUTH_TYPES.LOGIN_FAIL } })
    }

}

export const userLogoutACT = () => async (dispatch: Dispatch<AUTH_TYPES.TAuthActions>) => {
    dispatch({ type: AUTH_TYPES.LOGOUT })
}

export const userRegisterACT = (achiever: IAchiever) => async (dispatch: Dispatch<AUTH_TYPES.TAuthActions>) => {
    const url = TYPES.API_URL_PREFIX + AUTH_TYPES.REGISTER_URL_POSTFIX
    const req: TAchieverAPIRequest = { achiever: achiever }

    try {
        const res: AxiosResponse<TAchieverAPIResponse> = await axios.post(url, req)
        dispatch({ type: AUTH_TYPES.CLEAR_ERROR })
        dispatch({ type: AUTH_TYPES.REGISTER_SUCCESS, payload: res.data })

    } catch (err) {
        dispatch({ type: AUTH_TYPES.REGISTER_FAIL, payload: { ...err?.response, id: AUTH_TYPES.REGISTER_FAIL } })
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
        const res: AxiosResponse<TAchieverAPIResponse> = await axios.post(url, req)
        dispatch({ type: AUTH_TYPES.USER_REFRESH, payload: res.data });
        if (retryAction) {
            dispatch(retryAction)
        }
    } catch (err) {
        dispatch({ type: AUTH_TYPES.AUTH_ERROR, payload: { ...err?.response, id: AUTH_TYPES.AUTH_ERROR } })
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