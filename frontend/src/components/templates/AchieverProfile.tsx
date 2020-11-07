import React, { Fragment } from 'react'
import { render } from 'react-dom'
import { useSelector } from 'react-redux'
import { RootState } from '../../redux/reducers'
import { IAchiever, TGoals } from '../../types/AchieverTypes'
import { Achiever } from '../Achiever'
import { GoalPage } from './GoalPage'

export const AchieverProfile = () => {
    const achiever = useSelector((state: RootState) => { return state.app.achiever })
    const goals = achiever?.goals

    const renderAchiever = () => {
        return (
            <Achiever achiever={achiever || {}} />
        )
    }
    const renderGoals = (goals: TGoals) => {
        return Object.entries(goals).map(([key, value]) => {
            const goalID = parseInt(key)
            const isInProgress = value

            if (isInProgress) {
                return (
                    <GoalPage id={goalID} />
                )
            }
            return (
                null
            )
        })
    }


    return (
        <Fragment>
            {renderAchiever()}
            {goals && Object.getOwnPropertyNames(goals).length > 0 && renderGoals(goals)}
        </Fragment>
    )
}
