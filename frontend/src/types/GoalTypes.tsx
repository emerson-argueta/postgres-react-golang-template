import { IAuthorization } from "./AuthTypes"

export interface IGoal {
    id: number
    name: string
    achievers: TAchievers
}

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

export interface IGoalAPIResponse {
    id?: string
    name?: string
    achievers?: string
    error?: string
}
export interface IGoalAPIRequest {
    id?: number
    name?: string
    state?: string
    progress?: number
    message?: string
    timestamp?: string
    authorization?: IAuthorization
}