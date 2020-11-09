import { Reducer } from "redux"
import * as APP_TYPES from "../../types/AppTypes"
import { IGoalAPIResponse, TGoals } from "../../types/GoalTypes"

const initialState: APP_TYPES.IAppState = {
    achiever: undefined,
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
            const goals = retrievedGoals && retrievedGoals.map(g => {
                return goalAPIResponseToGoals(state.goals || {}, g)
            })

            return {
                ...state,
                goals: goals
            }
        }
        case APP_TYPES.CREATE_GOAL_FAIL:
        case APP_TYPES.RETRIEVE_GOALS_FAIL: {
            return {
                ...state
            }
        }
        default:
            return state
    }
}

const goalAPIResponseToGoals = (goals: TGoals, newGoal: IGoalAPIResponse): TGoals => {
    const ng: TGoals | undefined = (newGoal.id && {
        [newGoal.id]: {
            achievers: newGoal.achievers,
            id: newGoal.id,
            name: newGoal.name
        }
    }) || undefined

    return (newGoal.id && Object.assign(goals, { [newGoal.id]: ng })) || goals
}