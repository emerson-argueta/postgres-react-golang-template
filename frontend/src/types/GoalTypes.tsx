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