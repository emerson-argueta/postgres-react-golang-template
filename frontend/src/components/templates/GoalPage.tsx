import React, { Fragment, useEffect, useState } from 'react'
import { useSelector } from 'react-redux'
import { RootState } from '../../redux/reducers'
import { IGoal, TAchieverGoal, TAchievers, TAchieverStats } from '../../types/GoalTypes'
import { Goal } from '../Goal'

interface TMetadata {
    achieversStats?: TAchieverStats
    name?: string
    id?: number
}
export const GoalPage = ({ id }: { id: number }) => {

    const goal = useSelector((state: RootState) => {
        return state.app.goals && state.app.goals[id]
    })
    const achievers = goal?.achievers

    const [metadata, setMetadata] = useState<TMetadata>({})

    useEffect(() => {
        if (goal) {
            const achieverStats: TAchieverStats = caculateAchieverStats(goal.achievers)
            setMetadata({ name: goal.name, id: goal.id, achieversStats: achieverStats })
        }
    }, [goal])

    const renderAchievers = (achievers: TAchievers) => {
        return Object.entries(achievers).map(([achieverUUID, achieverGoal]) => {
            return (
                // TODO: pass the necessary achiever data to create goal component
                <Goal achieverUUID={achieverUUID} achieverGoal={achieverGoal} />
            )
        })
    }
    const renderMetadata = (metatdata: TMetadata) => {
        return (
            <div>{metadata}</div>
        )
    }
    return (
        <Fragment>
            {achievers && renderAchievers(achievers)}
            {metadata && renderMetadata(metadata)}
        </Fragment>
    )
}

const caculateAchieverStats = (achievers: TAchievers): TAchieverStats => {
    const countAchievers: number = Object.getOwnPropertyNames(achievers).length
    const achieversCompletedReducer = (accumulator: number, achieverGoal: TAchieverGoal) => {

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