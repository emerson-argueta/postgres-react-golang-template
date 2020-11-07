import React, { Fragment, useEffect, useState } from 'react'
import { useSelector } from 'react-redux'
import { IAchiever } from '../../types/AchieverTypes'
import { IGoal, TAchieverGoal, TAchievers, TAchieverStats } from '../../types/GoalTypes'
import { Goal } from '../Goal'

interface TMetadata {
    achieversStats?: TAchieverStats
    name?: string
    id?: number
}
export const GoalPage = ({ id }: { id: number }) => {
    // TODO: create type or interface for metaData
    const [metadata, setMetadata] = useState<TMetadata>({})

    // TODO: create interface for application state
    const goal: IGoal | null = useSelector((state: any) => {
        return state.goals[id]
    })
    const achieversGoal = goal?.achievers

    useEffect(() => {
        // TODO: extract transform and load achievers data into metadata
        if (goal) {
            const achieverStats: TAchieverStats = caculateAchieverStats(goal.achievers)
            setMetadata({ name: goal.name, id: goal.id, achieversStats: achieverStats })
        }
    }, [goal])

    const renderAchievers = () => {
        return achieversGoal.map((achieverGoal) => {
            return (
                // TODO: pass the necessary achiever data to create goal component
                <Goal achieverGoal={achieverGoal} />
            )
        })
    }

    return (
        <Fragment>
            {renderAchievers()}
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