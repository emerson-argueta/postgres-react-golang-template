import { IAuthorization } from "./AuthTypes"

export interface IGoal {
    id?: number
    name?: string
    achievers?: TAchievers
}
export type TGoals = { [id: number]: IGoal }

export type TAchievers = {
    [achieverUUID: string]: TAchieverGoal
}
export type TAchieverGoal = {
    state: string
    progress: number
    messages: TMessages
}

export type TAchieverStats = {
    countAchievers?: number,
    achieversCompleted?: number
}
export type TMessages = { [timestamp: string]: string }

export interface IGoalAPIResponse extends IGoal {
    error?: string
}
export type TGoalAPIRequest = {
    id?: number
    name?: string
    state?: string
    progress?: number
    message?: string
    timestamp?: string
    authorization?: IAuthorization
}