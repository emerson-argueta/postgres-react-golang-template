import React, { useState, Fragment, useEffect } from 'react'
import { Button } from '@material-ui/core'
import { useSelector, useDispatch } from 'react-redux'
import { IChurch } from '../../../types/Interface/ChurchInterfaces'
import { IAuthState, IAppState } from '../../../types/Interface/Interfaces'
import { ChurchesTable } from './ChurchesTable'
import { FormSelect } from '../FormSelect'
import { Forms } from '../Forms'
import { AppButton } from '../AppButton'
import { IAdministrator } from '../../../types/Interface/AdministratorInterfaces'
import { createChurchACT, addChurchACT } from '../../../redux/actions/AppActions'

export const Churches = (
    { onClick, visible }: { onClick: Function, visible: boolean }
) => {
    const [open, setOpen] = useState(false)
    const [formSelect, setFormSelect] = useState('')
    const [formDetails, setFormDetails] =  useState<any>()

    const churches: Array<IChurch> | null = useSelector(
        (state: { auth: IAuthState, app: IAppState }) => {

            if (
                state.app.administrator.churches
            ) {
                return Object.values(state.app.administrator.churches)
            }

            return null

        }
    )
    const administrator: IAdministrator = useSelector(
        (state: { auth: IAuthState, app: IAppState }) => {
            return state.app.administrator
        }
    )
    const dispatch = useDispatch()
    useEffect(() => {
        if (formSelect === '') {
            setFormSelect(formItems[1].key)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);
    const churchesAppButton = () => {
        return (
            <div className='apphome-container' onClick={()=>onClick()}>
                <AppButton
                    name='Churches'
                    setOpen={setOpen}
                />
            </div>
        )

    }
    const leftView = () => {
        return (
            <div key="administrators-leftview" className='leftview-container'>

                <div className='leftview-header'>
                    <Button
                        onClick={
                            () => {
                                setOpen(false)
                                onClick()
                            }
                        }
                    >
                        Menu
                    </Button>
                </div>
                <div className="leftview-table">
                    {
                        (churches !== null && Object.keys(churches).length !== 0) ?
                            <ChurchesTable
                                administratorName={ (administrator.firstname || '') + (administrator.lastname || '') }
                                churches={churches}
                            /> :
                            null
                    }
                </div>
            </div>
        )
    }

    const formItems: { [index: number]: { text: string, key: string, onSelect: Function } } = {
        1: { text: 'Create New Church', key: 'form-item-1', onSelect: setFormSelect },
        2: { text: 'Add Existing Church', key: 'form-item-2', onSelect: setFormSelect }
    }
    const forms: {
        [form: string]: {
            [field: string]: { type: string, name: string }
        }
    } = {
        [formItems[1].key]: { 
            'type': { name: 'Type', type: 'text' }, 
            'name': { name: 'Name', type: 'text' }, 
            'address': { name: 'Address', type: 'text' },
            'phone': { name: 'Phone', type: 'text' }, 
            'email': { name: 'Email', type: 'text' }, 
            'password': { name: 'Password', type: 'text' } 
        },
        [formItems[2].key]: { 
            'email': { name: 'Email', type: 'text' }, 
            'password': { name: 'Password', type: 'text' } 
        }

    }
    const formsActions: { [form: string]: Function } = {
        [formItems[1].key]: () => {
            console.log('Test action form-1',formDetails)
            const church:any = {
                type:formDetails?.type,
                name:formDetails?.name,
                address:formDetails?.address,
                phone:formDetails?.phone,
                email:formDetails?.email,
                password:formDetails?.password
            }
            
            dispatch(createChurchACT(church))
            
        },
        [formItems[2].key]: () => {
            const password = formDetails?.password
            const email = formDetails?.email
             
            dispatch(addChurchACT(email,password))
        },
    }
    const rightView = () => {
        return (
            <div key="administrators-rightview" className='rightview-container'>
                <div className='rightview-formselect'>
                    {<FormSelect formItems={formItems} />}
                </div>
                <div className='rightview-form'>
                    <Forms
                        formSelect={formSelect}
                        forms={forms}
                        onChange={setFormDetails}
                        formsActions={formsActions}
                    />
                </div>
            </div>
        )
    }

    return (
        <Fragment>
            {
                open ?
                    [leftView(), rightView()] :
                    visible?
                    churchesAppButton():
                    null     
            }
        </Fragment>
    )
}
