import React, { Fragment, useState } from 'react'
import { IAchiever, TGoals } from '../types/AchieverTypes';
import { Goal } from './Goal';
import { GoalPage } from './templates/GoalPage';

export const Achiever = ({ achiever }: { achiever: IAchiever }) => {
    const [openGoal, setOpenGoal] = useState<{ [goalID: number]: boolean }>()

    const renderAchiever = () => {

        return (
            <Fragment>
                <div>{achiever.firstname}</div>
                <div>{achiever.lastname}</div>
                <div>{achiever.address}</div>
                <div>{achiever.phone}</div>
            </Fragment>
        )
    }
    const renderGoals = (goals: TGoals) => {
        return (
            Object.getOwnPropertyNames(goals).map(key => {
                const goalID: number = parseInt(key)
                return (
                    openGoal < GoalPage id = { goalID } />
                )
            })
        )
    }

return (
    <Fragment>
        {renderAchiever()}
        {achiever.goals && renderGoals(achiever.goals)}
    </Fragment>
)
}
