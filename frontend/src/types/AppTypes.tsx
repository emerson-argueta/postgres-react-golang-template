import { IAchiever, IAchieverAPIResponse, TAchievers } from "./AchieverTypes";
import { IGoalAPIResponse, TGoals } from "./GoalTypes";
import { IError } from "./Types";

export const ACHIEVER_URL_POSTFIX = '/achiever'
export const LOAD_ACHIEVER = 'APP_LOAD_ACHIEVER'
interface IAchieverAction {
    type: typeof LOAD_ACHIEVER
    payload?: IAchieverAPIResponse
}

export const GOAL_URL_POSTFIX = '/goal'
export const CREATE_GOAL_SUCCESS = "APP_CREATE_GOAL_SUCCESS"
export const CREATE_GOAL_FAIL = "APP_CREATE_GOAL_FAIL"
interface ICreateGoalAction {
    type: typeof CREATE_GOAL_SUCCESS | typeof CREATE_GOAL_FAIL
    payload?: IGoalAPIResponse
    error?: IError
}
export const RETRIEVE_GOALS_SUCCESS = "APP_RETREIVE_GOALS_SUCCESS"
export const RETRIEVE_GOALS_FAIL = "APP_RETREIVE_GOALS_FAIL"
interface IRetrieveGoalAction {
    type: typeof RETRIEVE_GOALS_SUCCESS | typeof RETRIEVE_GOALS_FAIL
    payload?: Array<IGoalAPIResponse>
    error?: IError
}
export const GOAL_ACHIEVERS_URL_POSTFIX = '/achiever?goalID='
export const RETRIEVE_GOAL_ACHIEVERS_SUCCESS = "APP_RETRIEVE_GOAL_ACHIEVERS_SUCCESS"
export const RETRIEVE_GOAL_ACHIEVERS_FAIL = "APP_RETRIEVE_GOAL_ACHIEVERS_FAIL"
interface IRetrieveGoalAchieversAction {
    type: typeof RETRIEVE_GOAL_ACHIEVERS_SUCCESS | typeof RETRIEVE_GOAL_ACHIEVERS_FAIL
    payload?: Array<IAchieverAPIResponse>
    error?: IError
}

export type TAppActions = IAchieverAction | ICreateGoalAction | IRetrieveGoalAction | IRetrieveGoalAchieversAction

export interface IAppState {
    achiever?: IAchiever
    achievers?: TAchievers
    goals?: TGoals
    error?: IError
    loading: boolean
}
