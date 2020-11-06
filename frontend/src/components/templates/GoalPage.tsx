import React, { Fragment, useEffect, useState } from 'react'
import { useSelector } from 'react-redux'
import { IAchiever } from '../../types/AchieverTypes'
import { TAchieverGoal } from '../../types/GoalTypes'
import { Goal } from '../Goal'

export const GoalPage = ({ id }: { id: number }) => {
    // TODO: create type or interface for metaData
    const [metaData, setMetaData] = useState<any>({})

    // TODO: create interface for application state
    const achieversGoal: Array<TAchieverGoal> = useSelector((state: any) => {
        return state.goals[id].achievers
    })

    useEffect(() => {
        // TODO: extract transform and load achievers data into metadata
    }, [achieversGoal])

    const renderAchievers = () => {
        return achieversGoal.map(achieverGoal => {
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
