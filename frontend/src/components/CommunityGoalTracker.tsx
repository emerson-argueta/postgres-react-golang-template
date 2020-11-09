import React, { Fragment } from 'react'
import { useSelector } from 'react-redux'
import { RootState } from '../redux/reducers'
import { AchieverProfile } from './templates/AchieverProfile'

export const CommunityGaolTracker = () => {
    const auth = useSelector((state: RootState) => { return state.auth })

    const renderAchieverProfile = () => {
        return (
            <AchieverProfile />
        )
    }
    return (
        <Fragment>
            {auth.isAuthenticated && renderAchieverProfile()}
        </Fragment>
    )
}
