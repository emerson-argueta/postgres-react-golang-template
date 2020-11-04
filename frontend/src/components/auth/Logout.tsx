import React, { Fragment } from 'react'
import { Button } from '@material-ui/core'
import { useDispatch } from 'react-redux'
import { userLogoutACT } from '../../redux/actions/AuthActions'

export const Logout = () => {
    const dispatch = useDispatch()
    const logout = () => dispatch(userLogoutACT())

    return (
        <Fragment>
            <Button
                color="inherit"
                onClick={logout}
            >
                Logout
            </Button>
        </Fragment>
    )
}
