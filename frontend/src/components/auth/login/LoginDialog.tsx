import React, { useState, useRef, Fragment, useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'

import { IAdministrator } from '../../../types/Interface/AdministratorInterfaces'
import { userLoginACT } from '../../../redux/actions/AuthActions'
import { LoginForm } from './LoginForm'
import { Button, Dialog, DialogContent, DialogActions } from '@material-ui/core'
import { Alert } from '@material-ui/lab'
import * as AUTH_TYPES from '../../../types/AuthTypes'

export const LoginDialog = () => {
    const [open, setOpen] = useState(false)
    const [msg, setMsg] = useState(null)
    const [formDetails, setFormDetails] = useState(null)

    const auth: AUTH_TYPES.IAuthState = useSelector((state: { app: AUTH_TYPES.IAppState, auth: AUTH_TYPES.IAuthState }) => {
        return state.auth
    })

    const dispatch = useDispatch()
    const login = (administrator: IAdministrator | null) => dispatch(userLoginACT(administrator))

    useEffect(() => {
        // Check for register error
        if (auth.error?.id === AUTH_TYPES.LOGIN_FAIL) {
            setMsg(auth.error.msg as any);
        } else {
            setMsg(null);
        }

        // If authenticated, close modal
        if (open) {
            if (auth.isAuthenticated) {
                setOpen(false);
            }
        }
    }, [auth.error, auth.isAuthenticated, open]);


    const loginButtonRef = useRef<HTMLButtonElement>(null)
    const handleOnEnter = () => {
        if (loginButtonRef.current !== null) {
            loginButtonRef.current.focus();
        }
    };
    return (
        <Fragment>
            <Button
                color="inherit"
                onClick={() => setOpen(true)}
            >
                Login
            </Button>
            <Dialog
                onEnter={handleOnEnter}
                id="standard-dialog"
                open={open}
                onClose={() => setOpen(false)}

            >
                <DialogContent>
                    {msg ? <Alert severity="error">{msg}</Alert> : null}
                    <LoginForm
                        onChange={setFormDetails}
                    />
                </DialogContent>

                <DialogActions>

                    <Button
                        ref={loginButtonRef}
                        variant="contained"
                        color="default"
                        disableElevation
                        onClick={() => login(formDetails)}
                    >
                        Login
                    </Button>

                    <Button
                        variant="contained"
                        color="default"
                        disableElevation
                        onClick={() => setOpen(false)}
                    >
                        Cancel
                </Button>
                </DialogActions>
            </Dialog>
        </Fragment>
    )
}
