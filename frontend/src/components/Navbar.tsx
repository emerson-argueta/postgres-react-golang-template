import React, { Fragment, useEffect } from 'react'
import { useSelector } from 'react-redux'
import { AppBar, Typography, Toolbar } from '@material-ui/core'
import { Login } from './auth/Login';
import { Header } from './Header';

import { useStyles } from './navbar/MaterialUIStyles';
import { Logout } from './auth/Logout';

import './navbar/Navbar.css'
import { Register } from './auth/Register';
import * as AUTH_TYPES from '../types/AuthTypes';
import * as APP_TYPES from '../types/Types';

export const Navbar = () => {
    console.log("rendering");

    const auth: AUTH_TYPES.IAuthState | null = useSelector((state: { app: APP_TYPES.IAppState, auth: AUTH_TYPES.IAuthState }) => { return state.auth })
    const app: APP_TYPES.IAppState | null = useSelector((state: { app: APP_TYPES.IAppState, auth: AUTH_TYPES.IAuthState }) => { return state.app })
    const user = app?.achiever || {}

    useEffect(() => {

    }, [auth.isAuthenticated]);

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
            <Register />
            <Login />
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