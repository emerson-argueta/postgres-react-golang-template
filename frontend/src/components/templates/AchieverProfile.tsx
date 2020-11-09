import React, { Fragment, useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { retrieveGoalsACT } from '../../redux/actions/AppActions'
import { RootState } from '../../redux/reducers'
import { TGoals } from '../../types/GoalTypes'
import { Achiever } from '../Achiever'
import { GoalPage } from './GoalPage'

export const AchieverProfile = () => {
    const achiever = useSelector((state: RootState) => { return state.app.achiever })
    const goals = useSelector((state: RootState) => { return state.app.goals })

    const dispatch = useDispatch()
    useEffect(() => {
        console.log("using effect in achevier profile");

        if (achiever) {

            dispatch(retrieveGoalsACT())

        }
    }, [achiever, dispatch]);

    const [selectedGoal, setSelectedGoal] = useState<number>()
    const [openGoal, setOpenGoal] = useState<boolean>(false)

    const renderAchiever = () => {
        return (
            achiever && <Achiever achiever={achiever} />
        )
    }
    const renderGoalList = (goals: TGoals) => {
        return Object.entries(goals).map(([key, value]) => {
            const goalID = parseInt(key)
            const goal = value
            console.log("rendering goal list", goal);

            return (
                // Todo change to modal
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
