import {combineReducers} from 'redux';
import AppReducer from './AppReducer';
import AuthReducer from './AuthReducer';
import { IAction } from '../../types/Interface/Interfaces';
import { LOGOUT } from '../../types/AuthTypes';


const appReducer = combineReducers(
  {
    app:AppReducer,
    auth:AuthReducer
  }
)

export default (state:any, action: IAction) => {
  if(action.type===LOGOUT){
    state = undefined
  }
  
  return appReducer(state,action)
}