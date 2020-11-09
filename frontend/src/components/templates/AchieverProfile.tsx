import { Button } from '@material-ui/core'
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

    const [selectedGoal, setSelectedGoal] = useState<number>()
    const [openGoal, setOpenGoal] = useState<boolean>(false)

    const dispatch = useDispatch()
    useEffect(() => {
        console.log("using effect in achevier profile");

        if (achiever) {

            dispatch(retrieveGoalsACT())

        }
    }, [achiever, dispatch]);

    const renderAchiever = () => {
        return (
            achiever && <Achiever achiever={achiever} />
        )
    }
    const renderGoalList = (goals: TGoals) => {
        return Object.entries(goals).map(([key, value]) => {
            const goalID = parseInt(key)
            const goal = value

            return (
                // Todo change to modal

                <Button
                    key={goalID}
                    onClick={() => {
                        setSelectedGoal(goalID);
                        setOpenGoal(!openGoal);
                    }}
                >
                    {goal.name}
                </Button>
            )
        })
    }
    const renderGoal = (goalID: number) => {
        console.log("rendering goal", goalID);

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
