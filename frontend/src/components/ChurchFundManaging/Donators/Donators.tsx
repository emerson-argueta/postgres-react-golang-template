import React, { useState, Fragment, useEffect } from 'react'
import { Button } from '@material-ui/core'
import { useSelector, useDispatch } from 'react-redux'
import { IChurch } from '../../../types/Interface/ChurchInterfaces'
import { IAuthState, IAppState } from '../../../types/Interface/Interfaces'
import { IDonator } from '../../../types/Interface/DonatorInterfaces'
import { DonatorsTable } from './DonatorsTable'
import { ChurchSelect } from '../Churches/ChurchSelect'
import { FormSelect } from '../FormSelect'
import { Forms } from '../Forms'
import { AppButton } from '../AppButton'
import { churchDonatorsLoadACT, createDonatorACT } from '../../../redux/actions/AppActions'

export const Donators = (
    { onClick, visible }: { onClick: Function, visible: boolean }
) => {
    const [open, setOpen] = useState(false)
    const [formSelect, setFormSelect] = useState('')
    const [formDetails, setFormDetails] = useState<any>()

    const donators: Array<IDonator> | null = useSelector(
        (state: { auth: IAuthState, app: IAppState }) => {

            if (
                state.app.selectedChurch !== 0 &&
                state.app.administrator.churches &&
                state.app.administrator.churches[state.app.selectedChurch] &&
                state.app.administrator.churches[state.app.selectedChurch].donators
            ) {
                const donators = state.app.administrator.churches[state.app.selectedChurch].donators || {}

                return Object.values(donators)
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
    const churchDonatorsLoad = (churchID:number) => dispatch(churchDonatorsLoadACT(churchID))
    useEffect(() => {
        if (formSelect === '') {
            setFormSelect(formItems[1].key)
        }
        if(
            donators &&
            selectedChurch && 
            selectedChurch.id
        ){
            const dd = donators
            if(dd && !dd.pop()?.firstname){
                churchDonatorsLoad(selectedChurch.id)
            }
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedChurch?.id]);

    const donatorsAppButton = () => {
        return (
            <div className='apphome-container' onClick={() => onClick()}>
                <AppButton
                    name='Donators'
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
                            <DonatorsTable
                                church={selectedChurch}
                                donators={donators}
                            /> :
                            null
                    }
                    
                </div>
            </div>
        )
    }

    const formItems: { [index: number]: { text: string, key: string, onSelect: Function } } = {
        1: { text: 'New Donator', key: 'form-item-1', onSelect: setFormSelect }
    }
    const forms: {
        [form: string]: {
            [field: string]: { type: string, name: string }
        }
    } = {
        [formItems[1].key]: { 
            'firstname': { name: 'First Name', type: 'text' }, 
            'lastname': { name: 'Last Name', type: 'text' }, 
            'email': { name: 'Email', type: 'text' },
            'address': { name: 'Address', type: 'text' },
            'phone': { name: 'Phone', type: 'text' }
        },

    }
    const formsActions: { [form: string]: Function } = {
        
        [formItems[1].key]: () => {
            const donator:IDonator = {
                firstname:formDetails?.firstname,
                lastname:formDetails?.lastname,
                email:formDetails?.email === "" ? null : formDetails?.email,
                address:formDetails?.address === "" ? null : formDetails?.address,
                phone:formDetails?.phone === "" ? null : formDetails?.phone
            }
            if(selectedChurch?.id){
                dispatch(createDonatorACT(selectedChurch?.id,donator))
            }
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
                        donatorsAppButton() :
                        null

            }
        </Fragment>

    )
}
