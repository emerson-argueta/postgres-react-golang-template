import React, { useState, useEffect } from 'react'
import { TextField, Button } from '@material-ui/core';
import Autocomplete from '@material-ui/lab/Autocomplete';

export const Forms = (
    { formSelect, forms, onChange, formsActions }:
        {
            formSelect: string,
            forms: {
                [form: string]: {
                    [field: string]: { type: string, name: string, action?: Function, options?: Array<any> }
                }
            },
            onChange: Function,
            formsActions: { [form: string]: Function },
            

        }
) => {
    const [formDetails, setFormDetails] = useState<any>()
    const handleChangeFormDetails = (e: React.ChangeEvent<any> | { target: { id: any, value: any } }) => {
        const fd: any = formDetails
        fd[e.target.id] = e.target.value

        setFormDetails(fd)
        onChange(fd)
    }
    const currentDateUTCStr = new Date().toISOString()
    useEffect(() => {
        if(!forms || !forms[formSelect]){return}
        
        const ifd: { [index: string]: any } = {}
        const fd = Object.entries(forms[formSelect]).reduce((acc, curr) => {
            const key = curr[0]
            const field = curr[1]
            if (field.type === 'autocomplete' && field.options) {
                const firstOption:{value:any,id:any,label:string} = field.options.find((o, i) => i === 0)
                acc[firstOption?.id] = field.options.find((o, i) => i === 0)?.value
            }else if (field.type === 'datetime-local') {
                acc[key] = currentDateUTCStr.replace(/Z+$/,'')
            } else {
                acc[key] = ''
            }

            return acc
        }, ifd)
        setFormDetails(fd)
        onChange(fd)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [formSelect]);

    
    const renderForm = () => {
        if(!forms || !forms[formSelect]){return null}
        return Object.entries(forms[formSelect]).map((entry) => {
            const key = entry[0]
            const field = entry[1]
            if (field.type === 'autocomplete') {
                const options = field.options
                if (!options) { return null }

                return (
                    <Autocomplete
                        key={key}
                        id={key}
                        options={options}
                        getOptionLabel={(option) => option.label || ''}
                        onChange={(e, newValue) => {
                            if (newValue && newValue.value) {
                                const e = { target: { id: newValue.id, value: newValue.value } }
                                handleChangeFormDetails(e)
                            }

                        }}
                        style={{ width: 300 }}
                        renderInput={(params) => {
                            return (
                                <TextField
                                    {...params}
                                    label={field.name}
                                    variant="outlined"
                                />
                            )
                        }
                        }
                    />
                )
            }
            
            return (
                <TextField
                    defaultValue={field.type === 'datetime-local' ? currentDateUTCStr.replace(/Z+$/,'') :undefined}
                    key={key}
                    variant="outlined"
                    name={field.name}
                    id={key}
                    helperText={field.name}
                    type={field.type}
                    onChange={(e) => {
                        handleChangeFormDetails(e)
                    }}
                />
            )

        })
    }
    return (
        <div>
            <form>
                {renderForm()}
                <Button
                    variant="contained"
                    color="default"
                    disableElevation
                    onClick={() => formsActions[formSelect]()}
                >
                    Submit
                </Button>
            </form>
        </div>
    )
}
