import React, { Fragment, useState, useRef, useEffect } from 'react'
import { Button, Dialog, DialogContent, DialogActions } from '@material-ui/core'
import { useDispatch, useSelector } from 'react-redux'
import { IAuthState, IAppState } from '../../../types/Interface/Interfaces'
import { userRegisterACT } from '../../../redux/actions/AuthActions'
import { RegisterForm } from './RegisterForm'
import {Alert} from '@material-ui/lab'
import { REGISTER_FAIL } from '../../../types/AuthTypes'
import { IAdministrator } from '../../../types/Interface/AdministratorInterfaces'

export const RegisterDialog = () => {
    const [open, setOpen] = useState(false)
    const [formDetails,setFormDetails] = useState(null)
    const [msg, setMsg] = useState(null)

    const auth:IAuthState = useSelector((state:{app:IAppState,auth:IAuthState})=>{
        return state.auth
    })
    
    const dispatch = useDispatch()
    const register = (administrator: IAdministrator|null) => dispatch(userRegisterACT(administrator))

    useEffect(() => {
        // Check for register error
        if (auth.error?.id === REGISTER_FAIL) {
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
    }, [auth.error,auth.isAuthenticated,open]);
    
    const handleChange = (val:any) =>{
        setFormDetails(val)
    }
      
    const loginButtonRef = useRef<HTMLButtonElement>(null)
    const handleEnterKeyPress = () =>{
        if (loginButtonRef.current !== null) {
            loginButtonRef.current.focus(); 
        }
    }
    return (
        <Fragment>
            <Button
                color="inherit"
                onClick={() => setOpen(true)}
            >
                Register
            </Button>
            <Dialog
                // onEnter={handleOnEnter}
                id="standard-dialog"
                open={open}
                onClose={() => setOpen(false)}

            >
                <DialogContent>
                    {msg ? <Alert severity="error">{msg}</Alert> : null}
                    <RegisterForm onChange={handleChange} handleEnterKeyPress={handleEnterKeyPress}/>
                </DialogContent>
                <DialogActions>

                    <Button
                        ref={loginButtonRef}
                        variant="contained"
                        color="default"
                        disableElevation
                        onClick={() => {register(formDetails)}}
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
        </Fragment>
    )
}
