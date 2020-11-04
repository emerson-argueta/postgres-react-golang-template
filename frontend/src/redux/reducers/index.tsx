import { combineReducers } from 'redux';
import AuthReducer from './AuthReducer';
import { IAction } from '../../types/Types';
import { LOGOUT } from '../../types/AuthTypes';


const appReducer = combineReducers(
  {
    auth: AuthReducer
  }
)

export default (state: any, action: IAction) => {
  if (action.type === LOGOUT) {
    state = undefined
  }

  return appReducer(state, action)
}