import React, { Fragment, useState } from 'react'
import { IAchiever } from '../types/AchieverTypes';

export const Achiever = ({ achiever }: { achiever: IAchiever }) => {

    const renderAchiever = () => {

        return (
            <Fragment>
                <div>{achiever.firstname}</div>
                <div>{achiever.lastname}</div>
                <div>{achiever.address}</div>
                <div>{achiever.phone}</div>
            </Fragment>
        )
    }

    return (
        <Fragment>
            {renderAchiever()}
        </Fragment>
    )
}
