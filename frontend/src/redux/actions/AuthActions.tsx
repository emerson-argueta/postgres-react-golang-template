import axios from 'axios'
import * as INTERFACES from "../../types/Interface/Interfaces"
import * as ADMINISTRATOR_INTERFACES from "../../types/Interface/AdministratorInterfaces"
import * as AUTH_TYPES from '../../types/AuthTypes'
import * as appActions from './AppActions'


export const userLoginACT = (administrator: ADMINISTRATOR_INTERFACES.IAdministrator | null) => (dispatch: Function, getState: () => { app: INTERFACES.IAppState, auth: INTERFACES.IAuthState }) => {

    const req = { administrator: { ...administrator } }
    axios.post("/api/administrator/login", req)
        .then(res => {
            dispatch({ type: AUTH_TYPES.CLEAR_ERROR })
            dispatch({ type: AUTH_TYPES.LOGIN_SUCCESS, payload: res.data })
            dispatch(appActions.userLoadACT())
        })
        .catch(err => {
            dispatch({ type: AUTH_TYPES.LOGIN_FAIL, payload: { ...err?.response, id: AUTH_TYPES.LOGIN_FAIL } })
        })

}

export const userLogoutACT = () => (dispatch: Function, getState: () => { app: INTERFACES.IAppState, auth: INTERFACES.IAuthState }) => {
    dispatch({ type: AUTH_TYPES.LOGOUT })
    dispatch({ type: AUTH_TYPES.LOGOUT_SUCCESS })
}

export const userRegisterACT = (administrator: ADMINISTRATOR_INTERFACES.IAdministrator | null) => (dispatch: Function, getState: () => { app: INTERFACES.IAppState, auth: INTERFACES.IAuthState }) => {

    const req = { administrator: { ...administrator } }
    axios.post("/api/administrator", req)
        .then(res => {

            dispatch({ type: AUTH_TYPES.CLEAR_ERROR })
            dispatch({ type: AUTH_TYPES.REGISTER_SUCCESS, payload: res.data })
            dispatch(appActions.userLoadACT())
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

export const userTokenRefreshACT = (retryAction?: Function) => (dispatch: Function, getState: Function) => {

    const token: ADMINISTRATOR_INTERFACES.IToken = {
        accesstoken: localStorage.getItem(AUTH_TYPES.TRUSTDONATIONS_ACCESS_TOKEN) || undefined,
        refreshtoken: localStorage.getItem(AUTH_TYPES.TRUSTDONATIONS_REFRESH_TOKEN) || undefined
    }

    const req = { token: token }
    axios
        .post('/api/administrator/authorization', req)
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
    // Get token from localstorage
    const accessToken = localStorage.getItem(AUTH_TYPES.TRUSTDONATIONS_ACCESS_TOKEN)

    // Headers
    const config: INTERFACES.IConfigHeaders = {
        headers: {
            'Content-type': 'application/json'
        }
    };

    // If token, add to headers
    if (accessToken) {
        config.headers['Authorization'] = "Bearer " + accessToken;
    }

    return config;
};