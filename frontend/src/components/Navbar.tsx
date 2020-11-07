import React, { Fragment, useEffect } from 'react'
import { useSelector } from 'react-redux'
import { AppBar, Typography, Toolbar } from '@material-ui/core'
import { Login } from './auth/Login';
import { Header } from './Header';
import { Logout } from './auth/Logout';
import { Register } from './auth/Register';

import './navbar/Navbar.css'
import { useStyles } from './navbar/MaterialUIStyles';

import { RootState } from '../redux/reducers';

export const Navbar = () => {
    console.log("rendering");

    const auth = useSelector((state: RootState) => { return state.auth })
    const user = useSelector((state: RootState) => { return state.app.achiever })

    useEffect(() => {

    }, [auth]);

    const classes = useStyles()
    const authLinks = (
        <Fragment>
            <Typography className="welcome-user" variant="h6">
                <strong >
                    {auth && auth.isAuthenticated ? `Welcome ${user?.firstname}` : ''}
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