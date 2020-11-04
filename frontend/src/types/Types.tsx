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
    msg: string
}

export interface IAppState {
    achiever: any
    error: IError | null
    loading: boolean
}