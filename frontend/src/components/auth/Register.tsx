import React, { Fragment, useState, useRef, useEffect } from 'react'
import { Button, Dialog, DialogContent, DialogActions } from '@material-ui/core'
import { useDispatch, useSelector } from 'react-redux'
import { userRegisterACT } from '../../redux/actions/AuthActions'
import { RegisterForm } from './register/RegisterForm'
import { Alert } from '@material-ui/lab'
import { REGISTER_FAIL } from '../../types/AuthTypes'
import * as AUTH_TYPES from '../../types/AuthTypes'
import { IAchiever } from '../../types/AchieverTypes'
import { RootState } from '../../redux/reducers'


export const Register = () => {
    const [open, setOpen] = useState(false)
    const [msg, setMsg] = useState(null)
    const [formDetails, setFormDetails] = useState(null)

    const auth: AUTH_TYPES.IAuthState = useSelector((state: RootState) => { return state.auth })

    const dispatch = useDispatch()
    const register = (user: IAchiever) => dispatch(userRegisterACT(user))

    useEffect(() => {
        if (auth.error?.id === REGISTER_FAIL) {
            setMsg(auth.error.msg as any);
        } else {
            setMsg(null);
        }

        if (open) {
            if (auth.isAuthenticated) {
                setOpen(false);
            }
        }
    }, [auth.error, auth.isAuthenticated, open]);

    const loginButtonRef = useRef<HTMLButtonElement>(null)
    const handleChange = (val: any) => {
        setFormDetails(val)
    }


    const handleEnterKeyPress = () => {
        if (loginButtonRef.current !== null) {
            loginButtonRef.current.focus();
        }
    }
    const registerDialog = () => {

        return (
            <Dialog
                id="standard-dialog"
                open={open}
                onClose={() => setOpen(false)}

            >
                <DialogContent>
                    {msg ? <Alert severity="error">{msg}</Alert> : null}
                    <RegisterForm onChange={handleChange} handleEnterKeyPress={handleEnterKeyPress} />
                </DialogContent>
                <DialogActions>

                    <Button
                        ref={loginButtonRef}
                        variant="contained"
                        color="default"
                        disableElevation
                        onClick={() => { register(formDetails || {}) }}
                    >
                        Register
                        </Button>

                    <Button
                        variant="contained"
                        color="default"
                        disableElevation
                        onClick={() => { setOpen(false) }}
                    >
                        Cancel
                        </Button>
                </DialogActions>
            </Dialog>
        )

    }

    return (
        <Fragment>
            <Button
                color="inherit"
                onClick={() => setOpen(true)}
            >
                Register
            </Button>
            {registerDialog()}
        </Fragment>
    )
}
