import { IAdministrator } from "./AdministratorInterfaces";

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
  id:string,
  status: string,
  msg: string
}

export interface ITransaction {
  id: string,
  text: string,
  amount: number
}

export interface IAppState {
  selectedChurch:number | 0,
  administrator: IAdministrator,

  error: IError | null,
  loading: boolean
}

export interface IAuthState {
  isAuthenticated: boolean,
  token: any,
  error: IError | null,
  loading: boolean
}

export interface IAccountStatement {

}