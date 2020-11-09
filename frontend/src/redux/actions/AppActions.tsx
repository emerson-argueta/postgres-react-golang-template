import axios, { AxiosResponse } from "axios";
import { Dispatch } from "redux";
import { IAchieverAPIResponse } from "../../types/AchieverTypes";
import * as APP_TYPES from "../../types/AppTypes"
import { IGoalAPIResponse } from "../../types/GoalTypes";
import * as TYPES from "../../types/Types"
import { tokenConfig } from "./AuthActions";

export const loadAchieverACT = (res: IAchieverAPIResponse): APP_TYPES.TAppActions => {
    return { type: APP_TYPES.LOAD_ACHIEVER, payload: res }
}

export const retrieveGoalsACT = () => async (dispatch: Dispatch<APP_TYPES.TAppActions>) => {
    const url = TYPES.API_URL_PREFIX + APP_TYPES.GOAL_URL_POSTFIX

    try {
        const res: AxiosResponse<Array<IGoalAPIResponse>> = await axios.get(url, tokenConfig())
        dispatch({ type: APP_TYPES.RETRIEVE_GOALS_SUCCESS, payload: res.data });
    } catch (err) {
        const error: TYPES.IError = { id: APP_TYPES.RETRIEVE_GOALS_FAIL, status: err?.response.status, msg: err?.response?.data?.error }
        dispatch({ type: APP_TYPES.RETRIEVE_GOALS_FAIL, error: error })
    }
}

export const retreiveGoalAchievers = (goalID: number) => async (dispatch: Dispatch<APP_TYPES.TAppActions>) => {
    const url = TYPES.API_URL_PREFIX + APP_TYPES.GOAL_ACHIEVERS_URL_POSTFIX + goalID

    try {
        const res: AxiosResponse<Array<IAchieverAPIResponse>> = await axios.get(url, tokenConfig())
        dispatch({ type: APP_TYPES.RETRIEVE_GOAL_ACHIEVERS_SUCCESS, payload: res.data });
    } catch (err) {
        const error: TYPES.IError = { id: APP_TYPES.RETRIEVE_GOAL_ACHIEVERS_FAIL, status: err?.response.status, msg: err?.response?.data?.error }
        dispatch({ type: APP_TYPES.RETRIEVE_GOAL_ACHIEVERS_FAIL, error: error })
    }
}