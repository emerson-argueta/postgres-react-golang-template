import axios from 'axios'
import * as AUTH_TYPES from '../../types/AuthTypes'
import * as TYPES from '../../types/Types'
import { IAchiever, TAchieverAPIResponse } from '../../types/AchieverTypes'


export const userLoginACT = (achiever: IAchiever) => (dispatch: Function, getState: () => { app: TYPES.IAppState, auth: AUTH_TYPES.IAuthState }) => {
    const url = TYPES.API_URL_PREFIX + AUTH_TYPES.LOGIN_URL_POSTFIX
    const req = { achiever: { ...achiever } }

    axios.post(url, req)
        .then(res => {
            const achieverResponse: TAchieverAPIResponse = res.data
            dispatch({ type: AUTH_TYPES.CLEAR_ERROR });
            dispatch({ type: AUTH_TYPES.LOGIN_SUCCESS, payload: achieverResponse });
            // TODO: dispatch(appActions.userLoadACT(achieverResponse))
        })
        .catch(err => {
            dispatch({ type: AUTH_TYPES.LOGIN_FAIL, payload: { ...err?.response, id: AUTH_TYPES.LOGIN_FAIL } })
        })

}

export const userLogoutACT = () => (dispatch: Function, getState: () => { app: TYPES.IAppState, auth: AUTH_TYPES.IAuthState }) => {
    dispatch({ type: AUTH_TYPES.LOGOUT })
    dispatch({ type: AUTH_TYPES.LOGOUT_SUCCESS })
}

export const userRegisterACT = (achiever: IAchiever) => (dispatch: Function, getState: () => { app: TYPES.IAppState, auth: AUTH_TYPES.IAuthState }) => {
    const url = TYPES.API_URL_PREFIX + AUTH_TYPES.REGISTER_URL_POSTFIX
    const req = { achiever: { ...achiever } }

    axios.post(url, req)
        .then(res => {

            dispatch({ type: AUTH_TYPES.CLEAR_ERROR })
            dispatch({ type: AUTH_TYPES.REGISTER_SUCCESS, payload: res.data })
            // TODO: dispatch(appActions.userLoadACT())
        })
        .catch(err => {
            dispatch({ type: AUTH_TYPES.REGISTER_FAIL, payload: { ...err?.response, id: AUTH_TYPES.REGISTER_FAIL } })
        })
}

export const userLoadingACT = () => (dispatch: Function, getState: Function) => {
    dispatch({ type: AUTH_TYPES.USER_LOADING });
};
export const userLoadedACT = () => (dispatch: Function, getState: Function) => {
    dispatch({ type: AUTH_TYPES.USER_LOADED });
};

export const userReAuthorizeACT = (retryAction?: Function) => (dispatch: Function, getState: Function) => {
    const url = TYPES.API_URL_PREFIX + AUTH_TYPES.REAUTHORIZE_URL_POSTFIX

    const token: AUTH_TYPES.IAuthorization = {
        accesstoken: localStorage.getItem(AUTH_TYPES.TRUSTDONATIONS_ACCESS_TOKEN),
        refreshtoken: localStorage.getItem(AUTH_TYPES.TRUSTDONATIONS_REFRESH_TOKEN)
    }

    const req = { token: token }
    axios
        .post(url, req)
        .then(res => {
            dispatch({ type: AUTH_TYPES.USER_REFRESH, payload: res.data })
            if (retryAction) {
                dispatch(retryAction)
            }
        })
        .catch(err => {
            dispatch({ type: AUTH_TYPES.AUTH_ERROR, payload: { ...err?.response, id: AUTH_TYPES.AUTH_ERROR } })
        });
}

// Setup config/headers and token
export const tokenConfig = () => {
    const accessToken = localStorage.getItem(AUTH_TYPES.TRUSTDONATIONS_ACCESS_TOKEN)

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