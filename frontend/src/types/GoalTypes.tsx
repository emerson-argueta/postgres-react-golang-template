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

export type TMessages = { [timestamp: string]: string }