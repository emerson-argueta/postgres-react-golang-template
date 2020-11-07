import { IAchiever } from "./AchieverTypes";
import { IGoal, IGoalAPIResponse } from "./GoalTypes";
import { IError } from "./Types";


export const CREATE_GOAL_URL_POSTFIX = '/goal'
export const CREATE_GOAL_SUCCESS = "APP_CREATE_GOAL_SUCCESS"
export const CREATE_GOAL_FAIL = "APP_CREATE_GOAL_FAIL"
interface ICreateGoalAction {
    type: typeof CREATE_GOAL_SUCCESS | typeof CREATE_GOAL_FAIL
    payload?: IGoalAPIResponse
    error?: IError
}

export type TAppActions = ICreateGoalAction

export interface IAppState {
    achiever?: IAchiever
    goals?: { [id: number]: IGoal }
    error?: IError
    loading: boolean
}
