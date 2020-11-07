import React, { Fragment, useEffect } from 'react'
import { TAchieverGoal } from '../types/GoalTypes'

export const Goal = ({ achieverGoal }: { achieverGoal: TAchieverGoal }) => {

    const renderAchieverGoal = () => {
        return (
            <div>{achieverGoal}</div>
        )
    }
    return (
        <Fragment>
            {renderAchieverGoal()}
        </Fragment>
    )
}
