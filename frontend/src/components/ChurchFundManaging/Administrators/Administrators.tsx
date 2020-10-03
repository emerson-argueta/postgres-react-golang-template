import React, { useState, useEffect, Fragment } from 'react'
import './Administrators.css'
import { IAppState, IAuthState } from '../../../types/Interface/Interfaces'
import { useSelector, useDispatch } from 'react-redux'
import { IAdministrator } from '../../../types/Interface/AdministratorInterfaces'
import { IChurch } from '../../../types/Interface/ChurchInterfaces'
import { AdministratorsTable } from './AdministratorsTable'
import { ChurchSelect } from '../Churches/ChurchSelect'
import { AppButton } from '../AppButton'
import { Button } from '@material-ui/core'
import { FormSelect } from '../FormSelect'
import { Forms } from '../Forms'
import {  churchAdministratorsLoadACT, updateUserACT } from '../../../redux/actions/AppActions'

export const Administrators = (
    { onClick, visible }: { onClick: Function, visible: boolean }
) => {
    const [open, setOpen] = useState(false)
    const [formSelect, setFormSelect] = useState('')
    const [formDetails, setFormDetails] = useState<any>()

    const administrators: Array<IAdministrator> | null = useSelector(
        (state: { auth: IAuthState, app: IAppState }) => {
            

            if (
                state.app.selectedChurch !== 0 &&
                state.app.administrator.churches &&
                state.app.administrator.churches[state.app.selectedChurch] &&
                state.app.administrator.churches[state.app.selectedChurch].administrators
            ) {
                return Object.values(state.app.administrator.churches[state.app.selectedChurch].administrators)
            }

            return null

        }
    )
    const selectedChurch: IChurch | null = useSelector(
        (state: { auth: IAuthState, app: IAppState }) => {
            if (
                state.app.administrator.churches &&
                state.app.selectedChurch !== 0 
            ) {
                return state.app.administrator.churches[state.app.selectedChurch]
            }
            return null

        }
    )
    const dispatch = useDispatch()
    const churchAdministratorsLoad = (churchID:number) => dispatch(churchAdministratorsLoadACT(churchID))
    
    useEffect(() => {
        if (formSelect === '') {
            setFormSelect(formItems[1].key)
        }

        if(
            administrators &&
            selectedChurch && 
            selectedChurch.id
        ){
            const aa = administrators
            if(aa && !aa.pop()?.email){
                churchAdministratorsLoad(selectedChurch.id)
            }
            
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedChurch?.id]);

    const administratorsAppButton = () => {
        return (
            <div className='apphome-container' onClick={() => onClick()}>
                <AppButton
                    name='Admnistrators'
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
                        selectedChurch !== null ?
                            <AdministratorsTable
                                church={selectedChurch}
                                administrators={administrators}
                            /> :
                            null
                    }

                </div>
            </div>
        )
    }


    const formItems: { [index: number]: { text: string, key: string, onSelect: Function } } = {
        1: { text: 'Update My Information', key: 'form-item-2', onSelect: setFormSelect },
    }
    const forms: {
        [form: string]: {
            [field: string]: { type: string, name: string }
        }
    } = {
        [formItems[1].key]: { 
            'firstname': { name: 'First Name', type: 'text' }, 
            'lastname': { name: 'Last Name', type: 'text' },
            'address': { name: 'Address', type: 'text' },
            'phone': { name: 'Phone', type: 'text' },
            'email': { name: 'Email', type: 'text' },
            'password': { name: 'Password', type: 'text' },
        }
    }
    const formsActions: { [form: string]: Function } = {
        [formItems[1].key]: () => {
            const user = {
                firstname:formDetails?.firstname,
                lastname:formDetails?.lastname,
                address:formDetails?.address,
                phone:formDetails?.phone,
                email:formDetails?.email === "" ? null : formDetails?.email,
                password:formDetails?.password === "" ? null : formDetails?.password
            }
            console.log('Test action form-1',formDetails)
            dispatch(updateUserACT(user))
        },
    }
    const rightView = () => {
        return (
            <div key="administrators-rightview" className='rightview-container'>
                <div className='rightview-churchselect'>
                    <ChurchSelect />
                </div>
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
                    visible ?
                        administratorsAppButton() :
                        null
            }
        </Fragment>

    )
}
