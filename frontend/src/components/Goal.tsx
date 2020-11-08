import React, { Fragment } from 'react'
import { TAchieverGoal } from '../types/GoalTypes'

type TProps = {
    achieverGoal: TAchieverGoal,
    achieverUUID: string
}
export const Goal = ({ achieverGoal, achieverUUID }: TProps) => {

    const renderAchieverGoal = () => {
        return (
            <div key={achieverUUID}>{achieverGoal}</div>
        )
    }
    return (
        <Fragment>
            {renderAchieverGoal()}
        </Fragment>
    )
}
