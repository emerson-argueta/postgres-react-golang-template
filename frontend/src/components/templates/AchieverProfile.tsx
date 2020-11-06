import React, { Fragment } from 'react'
import { render } from 'react-dom'
import { useSelector } from 'react-redux'
import { IAchiever } from '../../types/AchieverTypes'
import { Achiever } from '../Achiever'

export const AchieverProfile = () => {
    const achiever: IAchiever | null = useSelector((state: any) => {
        return state.achiever
    })

    const renderAchiever = () => {
        return (
            // TODO: pass achiever data from redux store and pass to Achiever
            // component.
            <Achiever achiever={achiever || {}} />
        )
    }

    return (
        <Fragment>
            {renderAchiever()}
        </Fragment>
    )
}
