import React, { useState, Fragment } from 'react'
import { TextField } from '@material-ui/core';
import { useStyles } from './MaterialUIStyles';

export const RegisterForm = ({ onChange,handleEnterKeyPress }: { onChange: Function,handleEnterKeyPress:Function }) => {
    const [formDetails, setFormDetails] = useState({ email: '', password: '', firstname: '', lastname: '' })

    const handleChangeFormDetails = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const fd: any = formDetails
        fd[e.target.name] = e.target.value
        setFormDetails(fd)
        onChange(fd)
    }

    const handleKeyPress = (e:React.KeyboardEvent<HTMLDivElement>)=>{
        if(e.key === "Enter"){
            handleEnterKeyPress()
        }
    }
    const classes = useStyles()
    return (
        <Fragment>
            <form >
                <TextField
                    className={classes.textField}
                    variant="outlined"
                    name="firstname"
                    id="standard-basic"
                    helperText="First Name"
                    type="text"
                    autoComplete="name"
                    onChange={handleChangeFormDetails}
                    autoFocus={true}
                />
                <TextField
                    className={classes.textField}
                    variant="outlined"
                    name="lastname"
                    id="standard-basic"
                    helperText="Last Name"
                    type="text"
                    autoComplete="name"
                    onChange={handleChangeFormDetails}
                />
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
                    onKeyPress={handleKeyPress}
                />
            </form>
        </Fragment>
    )
}
