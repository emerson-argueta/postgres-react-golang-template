import { IAuthorization } from "./AuthTypes"

export interface IAchiever {
    uuid?: string
    role?: string
    firstname?: string
    lastname?: string
    address?: string
    phone?: string
    goals?: TGoals
    email?: string
    password?: string
}
export type TAchievers = { [uuid: string]: IAchiever }
// TGoals is an object where the keys represent the goal ID and the
// value indicates true if the goal has not been abandoned, false
// otherwise.
export type TGoals = { [goalID: number]: boolean }

export interface IAchieverAPIResponse {
    achiever?: IAchiever
    authorization?: IAuthorization
}
export type TAchieverAPIRequest = {
    achiever?: IAchiever
    authorization?: IAuthorization
}