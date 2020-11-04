import React, { useState, Fragment } from 'react'
import { TextField } from '@material-ui/core';
import { useStyles } from './MaterialUIStyles';

export const LoginForm = ({ onChange }: { onChange: Function }) => {
    const [formDetails, setFormDetails] = useState({ email: '', password: '' })

    const handleChangeFormDetails = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const fd: any = formDetails
        fd[e.target.name] = e.target.value
        setFormDetails(fd)
        onChange(fd)
    }

    const classes = useStyles()
    return (
        <Fragment>
            <form >
                <TextField
                    className={classes.textField}
                    variant="outlined"
                    name="email"
                    id="standard-basic"
                    helperText="Email Address"
                    type="email"
                    autoComplete="email"
                    onChange={handleChangeFormDetails}

                />
                <TextField
                    className={classes.textField}
                    variant="outlined"
                    name="password"
                    id="standard-password-input"
                    helperText="Password"
                    type="password"
                    autoComplete="current-password"
                    onChange={handleChangeFormDetails}
                />
            </form>
        </Fragment>
    )
}
