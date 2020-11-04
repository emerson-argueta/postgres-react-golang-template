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
// TGoals is an object where the keys represent the goal ID and the
// value indicates true if the goal has not been abandoned, false
// otherwise.
export type TGoals = { [goalID: number]: boolean }

export type TAchieverAPIResponse = {
    achiever: IAchiever
    authorization: IAuthorization
}