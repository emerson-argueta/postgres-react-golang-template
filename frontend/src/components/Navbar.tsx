import React, { Fragment } from 'react'
import { useSelector } from 'react-redux'
import { AppBar, Typography, Toolbar } from '@material-ui/core'
import { LoginDialog } from './auth/login/LoginDialog';
import { Header } from './Header';

import { useStyles } from './navbar/MaterialUIStyles';
import { Logout } from './auth/Logout';

import './Navbar.css'
import { RegisterDialog } from './auth/register/RegisterDialog';
import * as AUTH_TYPES from '../types/AuthTypes';
import * as APP_TYPES from '../types/Types';

export const Navbar = () => {
    const auth: AUTH_TYPES.IAuthState | null = useSelector((state: { app: APP_TYPES.IAppState, auth: AUTH_TYPES.IAuthState }) => { return state.auth })
    const app: APP_TYPES.IAppState | null = useSelector((state: { app: APP_TYPES.IAppState, auth: AUTH_TYPES.IAuthState }) => { return state.app })
    const user = app.user

    const classes = useStyles()
    const authLinks = (
        <Fragment>
            <Typography className="welcome-user" variant="h6">
                <strong >
                    {auth && auth.isAuthenticated ? `Welcome ${user.firstname}` : ''}
                </strong>
            </Typography>
            <Logout />
        </Fragment>
    );

    const guestLinks = (
        <Fragment>
            <RegisterDialog />
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