import React from 'react'
import { Button } from '@material-ui/core'

export const AppButton = (
    {name,setOpen}:
    {name:string,setOpen:(open:boolean)=>void}
) => {
    return (
        <div className="appbuttons">
            <Button
                variant="contained"
                color="default"
                disableElevation
                onClick={() => { setOpen(true) }}
            >
                {name}
            </Button>
            
        </div>
    )
}
