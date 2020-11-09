import React, { Fragment, useEffect, useState } from 'react'
import { useSelector } from 'react-redux'
import { RootState } from '../../redux/reducers'
import { TAchieverGoal, TAchievers, TAchieverStats } from '../../types/GoalTypes'
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
    const achieverUUIDs = goal?.achievers && Object.getOwnPropertyNames(goal.achievers)

    const [metadata, setMetadata] = useState<TMetadata>({})
    const [selectedAchiever, setSelectedAchiever] = useState<string>()
    const [openAchiever, setOpenAchiever] = useState<boolean>(false)


    useEffect(() => {
        if (goal?.achievers) {
            const achieverStats: TAchieverStats = caculateAchieverStats(goal.achievers)
            setMetadata({ name: goal.name, id: goal.id, achieversStats: achieverStats })
        }
    }, [goal])

    const renderAchieverList = (achieverUUIDs: Array<string>) => {
        return achieverUUIDs.map((achieverUUID) => {
            return (
                // Todo change to modal
                <div
                    key={achieverUUID}
                    onClick={() => {
                        setSelectedAchiever(achieverUUID);
                        setOpenAchiever(!openAchiever);
                    }}
                >
                    {"need achiever name here"}
                </div>
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
            <div>{metadata}</div>
        )
    }

    return (
        <Fragment>
            {achieverUUIDs && renderAchieverList(achieverUUIDs)}
            {metadata && renderMetadata(metadata)}
            {selectedAchiever && renderAchieverGoal(selectedAchiever, id)}
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