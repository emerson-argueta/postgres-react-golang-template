import { Reducer } from "redux"
import { IAchieverAPIResponse, TAchievers } from "../../types/AchieverTypes"
import * as APP_TYPES from "../../types/AppTypes"
import { IGoalAPIResponse, TGoals } from "../../types/GoalTypes"

const initialState: APP_TYPES.IAppState = {
    achiever: undefined,
    achievers: undefined,
    goals: undefined,
    error: undefined,
    loading: true
}

export const AppReducer: Reducer<APP_TYPES.IAppState, APP_TYPES.TAppActions> = (state = initialState, action) => {

    switch (action.type) {
        case APP_TYPES.LOAD_ACHIEVER: {
            const achiever = action.payload?.achiever
            return {
                ...state,
                achiever: achiever
            }
        }
        case APP_TYPES.CREATE_GOAL_SUCCESS: {
            const newGoal = action.payload
            const goals = state.goals && newGoal && goalAPIResponseToGoals(state.goals, newGoal)

            return {
                ...state,
                goals: goals
            }
        }
        case APP_TYPES.RETRIEVE_GOALS_SUCCESS: {
            const retrievedGoals = action.payload
            const goals = retrievedGoals && retrievedGoals.reduce((acc, g) => {
                return goalAPIResponseToGoals(acc, g)
            }, state.goals || {})

            return {
                ...state,
                goals: goals
            }
        }
        case APP_TYPES.RETRIEVE_GOAL_ACHIEVERS_SUCCESS: {
            const retrievedAchievers = action.payload
            const achievers = retrievedAchievers && retrievedAchievers.reduce((acc, a) => {
                return achieverAPIResponseToAchievers(acc, a)
            }, {} as TAchievers)

            return {
                ...state,
                achievers: achievers
            }
        }
        case APP_TYPES.CREATE_GOAL_FAIL:
        case APP_TYPES.RETRIEVE_GOALS_FAIL:
        case APP_TYPES.RETRIEVE_GOAL_ACHIEVERS_FAIL: {
            return {
                ...state,
                error: action.error
            }
        }
        default:
            return state
    }
}

const goalAPIResponseToGoals = (goals: TGoals, newGoal: IGoalAPIResponse): TGoals => {
    const ng: TGoals | undefined = (newGoal.id && {
        [newGoal.id]: {
            // achievers: newGoal.achievers,
            // id: newGoal.id,
            // name: newGoal.name
            ...newGoal
        }
    }) || undefined

    return (newGoal.id && Object.assign(goals, { ...ng })) || goals
}

const achieverAPIResponseToAchievers = (achievers: TAchievers, res: IAchieverAPIResponse): TAchievers => {
    const newAchiever = res.achiever
    const na: TAchievers | undefined = (newAchiever?.uuid && {
        [newAchiever.uuid]: {
            ...newAchiever
        }
    }) || undefined

    return (newAchiever?.uuid && Object.assign(achievers, { ...na })) || achievers
}