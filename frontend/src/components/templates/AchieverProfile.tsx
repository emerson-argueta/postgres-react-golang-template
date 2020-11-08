import React, { Fragment, useState } from 'react'
import { useSelector } from 'react-redux'
import { RootState } from '../../redux/reducers'
import { TGoals } from '../../types/GoalTypes'
import { Achiever } from '../Achiever'
import { GoalPage } from './GoalPage'

export const AchieverProfile = () => {
    const achiever = useSelector((state: RootState) => { return state.app.achiever })
    const goals = useSelector((state: RootState) => { return state.app.goals })

    const [selectedGoal, setSelectedGoal] = useState<number>()
    const [openGoal, setOpenGoal] = useState<boolean>(false)

    const renderAchiever = () => {
        return (
            <Achiever achiever={achiever || {}} />
        )
    }
    const renderGoalList = (goals: TGoals) => {
        return Object.entries(goals).map(([key, value]) => {
            const goalID = parseInt(key)
            const goal = value

            return (
                <div
                    key={goalID}
                    onClick={() => {
                        setSelectedGoal(goalID);
                        setOpenGoal(!openGoal);
                    }}
                >
                    {goal.name}
                </div>
            )
        })
    }
    const renderGoal = (goalID: number) => {
        return (
            openGoal && <GoalPage id={goalID} />
        )
    }

    return (
        <Fragment>
            {renderAchiever()}
            {goals && Object.getOwnPropertyNames(goals).length > 0 && renderGoalList(goals)}
            {selectedGoal && renderGoal(selectedGoal)}
        </Fragment>
    )
}
