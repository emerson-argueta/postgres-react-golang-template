import { Button } from '@material-ui/core'
import React, { Fragment, useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { retreiveGoalAchievers } from '../../redux/actions/AppActions'
import { RootState } from '../../redux/reducers'
import * as ACHIEVER_TYPES from '../../types/AchieverTypes'
import * as GOAL_TYPES from '../../types/GoalTypes'
import { Goal } from '../Goal'

interface TMetadata {
    achieversStats?: GOAL_TYPES.TAchieverStats
    name?: string
    id?: number
}
export const GoalPage = ({ id }: { id: number }) => {
    const goal = useSelector((state: RootState) => {
        return state.app.goals && state.app.goals[id]
    })
    const achievers = useSelector((state: RootState) => {
        return state.app.achievers
    })


    const [metadata, setMetadata] = useState<TMetadata>({})
    const [selectedAchiever, setSelectedAchiever] = useState<string>()
    const [openAchiever, setOpenAchiever] = useState<boolean>(false)

    const dispatch = useDispatch()
    useEffect(() => {
        if (goal?.achievers) {
            const achieverStats: GOAL_TYPES.TAchieverStats = caculateAchieverStats(goal.achievers)
            setMetadata({ name: goal.name, id: goal.id, achieversStats: achieverStats })
        }

        const achieverUUIDs = goal?.achievers && Object.getOwnPropertyNames(goal.achievers)
        if (achieverUUIDs && achieverUUIDs.length > 0) {
            dispatch(retreiveGoalAchievers(id))
        }
    }, [goal, dispatch, id])

    const renderAchieverList = (achievers: ACHIEVER_TYPES.TAchievers) => {
        return Object.entries(achievers).map(([achieverUUID, achiever]) => {
            return (
                // Todo change to modal
                <Button
                    key={achieverUUID}
                    onClick={() => {
                        setSelectedAchiever(achieverUUID);
                        setOpenAchiever(!openAchiever);
                    }}
                >
                    {achiever.firstname + " " + achiever.lastname}
                </Button>
            )
        })
    }
    const renderAchieverGoal = (achieverUUID: string, id: number) => {
        return (
            openAchiever && <Goal id={id} achieverUUID={achieverUUID} />
        )

    }
    const renderMetadata = (metatdata: TMetadata) => {
        return (
            <Fragment>
                <div>{metadata.name}</div>
                <div>{"Achievers with this goal: " + metadata.achieversStats?.countAchievers}</div>
                <div>{"Achievers who completed the goal: " + metadata.achieversStats?.achieversCompleted}</div>
            </Fragment>
        )
    }

    return (
        <Fragment>
            {metadata && renderMetadata(metadata)}
            {achievers && renderAchieverList(achievers)}
            {selectedAchiever && renderAchieverGoal(selectedAchiever, id)}
        </Fragment>
    )
}

const caculateAchieverStats = (achievers: GOAL_TYPES.TAchievers): GOAL_TYPES.TAchieverStats => {
    const countAchievers: number = Object.getOwnPropertyNames(achievers).length
    const achieversCompletedReducer = (accumulator: number, achieverGoal: GOAL_TYPES.TAchieverGoal) => {

        return accumulator + (achieverGoal.progress === 100 ? 1 : 0)
    }
    const achieversCompleted = Object.values(achievers).reduce(
        achieversCompletedReducer, 0
    )

    return {
        countAchievers: countAchievers,
        achieversCompleted: achieversCompleted
    }
}