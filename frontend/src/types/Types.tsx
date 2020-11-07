import { IAchiever } from "./AchieverTypes";
import { IAuthState } from "./AuthTypes";
import { IGoal } from "./GoalTypes";

export const API_URL_PREFIX = '/api/v1/communitygoaltracker'

export interface IConfigHeaders {
    headers: {
        [index: string]: string;
    };
}

export interface IAction {
    type: string;
    payload?: any;
}

export interface IError {
    id: string
    status: string
    msg?: string
}

export interface IAppState {
    achiever?: IAchiever
    goals?: { [id: number]: IGoal }
    error: IError | null
    loading: boolean
}

export interface IReduxState {
    app: IAppState
    auth: IAuthState
}