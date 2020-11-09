import React, { Fragment } from 'react'
import { useSelector } from 'react-redux'
import { RootState } from '../redux/reducers'
import { TAchieverGoal } from '../types/GoalTypes'

type TProps = {
    id: number,
    achieverUUID: string
}
export const Goal = ({ id, achieverUUID }: TProps) => {
    const achieverGoal = useSelector((state: RootState) => {
        const achievers = state.app.goals && state.app.goals[id].achievers
        return achievers && achievers[achieverUUID]
    })

    const renderAchieverGoal = (achieverGoal: TAchieverGoal) => {
        return (
            <div key={achieverUUID}>
                <div>
                    {achieverGoal.progress}
                </div>
                <div>
                    {achieverGoal.state}
                </div>
                {
                    achieverGoal.messages && Object.entries(achieverGoal.messages).map(([timestamp, message]) => {
                        return (
                            <div key={timestamp}>
                                {timestamp + '--->' + message}
                            </div>
                        )
                    })
                }
            </div>
        )
    }
    return (
        <Fragment>
            {achieverGoal && renderAchieverGoal(achieverGoal)}
        </Fragment>
    )
}
