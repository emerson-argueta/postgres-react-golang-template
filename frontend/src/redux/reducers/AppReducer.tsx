import { Reducer } from "redux"
import * as APP_TYPES from "../../types/AppTypes"

const initialState: APP_TYPES.IAppState = {
    achiever: undefined,
    goals: undefined,
    error: undefined,
    loading: true
}

export const AppReducer: Reducer<APP_TYPES.IAppState, APP_TYPES.TAppActions> = (state = initialState, action) => {

    switch (action.type) {
        default:
            return state
    }
}
