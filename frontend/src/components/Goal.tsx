import React, { Fragment, useEffect } from 'react'
import { useSelector } from 'react-redux'
import { IGoal, TAchieverGoal } from '../types/GoalTypes'

export const Goal = ({ achieverGoal }: { achieverGoal: TAchieverGoal }) => {

    const renderAchieverGoal = () => {
        return (
            <div>achieverGoal</div>
        )
    }
    return (
        <Fragment>
            {renderAchieverGoal()}
        </Fragment>
    )
}
