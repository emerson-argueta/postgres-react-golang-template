import React, { Fragment } from 'react'
import { useSelector } from 'react-redux'
import { RootState } from '../redux/reducers'
import { AchieverProfile } from './templates/AchieverProfile'

export const CommunityGaolTracker = () => {
    const isAuthenticated = useSelector((state: RootState) => { return state.auth.isAuthenticated })

    const renderAchieverProfile = () => {
        return (
            <AchieverProfile />
        )
    }
    return (
        <Fragment>
            {isAuthenticated && renderAchieverProfile()}
        </Fragment>
    )
}
