import { combineReducers } from 'redux';
import { AuthReducer } from './AuthReducer';
import { AppReducer } from './AppReducer';

export const rootReducer = combineReducers({
  auth: AuthReducer,
  app: AppReducer,
})

export type RootState = ReturnType<typeof rootReducer>