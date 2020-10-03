import React, { Fragment } from 'react'
import { IAuthState, IAppState } from '../../types/Interface/Interfaces'
import { useSelector } from 'react-redux'
import { AppBar, Typography, Toolbar } from '@material-ui/core'
import { LoginDialog } from '../Auth/Login/LoginDialog';
import { Header } from '../Header/Header';

import { useStyles } from './MaterialUIStyles';
import { Logout } from '../Auth/Logout/Logout';

import  './AppNavBar.css'
import { RegisterDialog } from '../Auth/Register/RegisterDialog';

export const AppNavBar = () => {
    const auth: IAuthState | null = useSelector((state: { app: IAppState, auth: IAuthState }) => { return state.auth })
    const app: IAppState | null = useSelector((state: { app: IAppState, auth: IAuthState }) => { return state.app })

    const classes = useStyles()
    const authLinks = (
        <Fragment>
            <Typography className="welcome-administrator" variant="h6">
                <strong >
                {auth && auth.isAuthenticated ? `Welcome ${app.administrator.firstname}` : ''}
                </strong>
            </Typography>
            <Logout />
        </Fragment>
    );

    const guestLinks = (
        <Fragment>
            <RegisterDialog/>
            <LoginDialog />
        </Fragment>
    );

    return (
        <div >

            <AppBar
                color={'transparent'}
                position={'static'}
            >
                <Toolbar>
                    <Typography variant="h6" className={classes.title}><strong><Header /></strong></Typography>

                    {auth && auth.isAuthenticated ? authLinks : guestLinks}

                </Toolbar>
            </AppBar>


        </div>
    )
}