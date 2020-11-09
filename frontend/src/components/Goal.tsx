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
            <div key={achieverUUID}>{achieverGoal}</div>
        )
    }
    return (
        <Fragment>
            {achieverGoal && renderAchieverGoal(achieverGoal)}
        </Fragment>
    )
}
