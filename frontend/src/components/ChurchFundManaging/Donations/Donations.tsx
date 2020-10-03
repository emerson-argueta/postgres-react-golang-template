import React, { useState, Fragment, useEffect } from 'react'
import { Button } from '@material-ui/core'
import { useSelector, useDispatch } from 'react-redux'
import * as CHURCH_TYPES from '../../../types/Interface/ChurchInterfaces'
import * as APP_TYPES from '../../../types/Interface/Interfaces'
import { DonationsTable } from './DonationsTable'
import { FormSelect } from '../FormSelect'
import { Forms } from '../Forms'
import { AppButton } from '../AppButton'
import { IDonator } from '../../../types/Interface/DonatorInterfaces'
import { ChurchSelect } from '../Churches/ChurchSelect'
import * as APP_ACTIONS from '../../../redux/actions/AppActions'
import * as DONATION_TYPES from '../../../types/Interface/DonationInterfaces'

export const Donations = (
    { onClick, visible }: { onClick: Function, visible: boolean }
) => {
    const [open, setOpen] = useState(false)
    const [formSelect, setFormSelect] = useState('')
    const [formDetails, setFormDetails] = useState<any>()

    const donatorsMap: CHURCH_TYPES.IDonators | null = useSelector(
        (state: { auth: APP_TYPES.IAuthState, app: APP_TYPES.IAppState }) => {

            if (
                state.app.selectedChurch !== 0 &&
                state.app.administrator.churches &&
                state.app.administrator.churches[state.app.selectedChurch] &&
                state.app.administrator.churches[state.app.selectedChurch].donators
            ) {
                const donators = state.app.administrator.churches[state.app.selectedChurch].donators || {}

                return donators
            }

            return null
        }
    )
    const donators = Object.values(donatorsMap || {}) as Array<IDonator>

    const selectedChurch: CHURCH_TYPES.IChurch | null = useSelector(
        (state: { auth: APP_TYPES.IAuthState, app: APP_TYPES.IAppState }) => {
            if (
                state.app.administrator.churches &&
                state.app.selectedChurch !== 0
            ) {
                return state.app.administrator.churches[state.app.selectedChurch]
            }
            return null

        }
    )
    const churchDonationReport: CHURCH_TYPES.IChurchDonationReport | null = useSelector(
        (state: { auth: APP_TYPES.IAuthState, app: APP_TYPES.IAppState }) => {

            if (
                state.app.selectedChurch !== 0 &&
                state.app.administrator.churches &&
                state.app.administrator.churches[state.app.selectedChurch] &&
                state.app.administrator.churches[state.app.selectedChurch].churchdonationreport
            ) {
                const churchDonationReport = state.app.administrator.churches[state.app.selectedChurch].churchdonationreport || {}

                return churchDonationReport
            }

            return null

        }
    )
    const donationStatementReport: CHURCH_TYPES.IDonationStatementReport | null = useSelector(
        (state: { auth: APP_TYPES.IAuthState, app: APP_TYPES.IAppState }) => {

            if (
                state.app.selectedChurch !== 0 &&
                state.app.administrator.churches &&
                state.app.administrator.churches[state.app.selectedChurch] &&
                state.app.administrator.churches[state.app.selectedChurch].donationstatementreport
            ) {
                const donationStatementReport = state.app.administrator.churches[state.app.selectedChurch].donationstatementreport || {}

                return donationStatementReport
            }

            return null

        }
    )

    const dispatch = useDispatch()
    const churchDonationsLoad = (churchID: number) => dispatch(APP_ACTIONS.churchDonationsLoadACT(churchID))
    useEffect(() => {
        if (formSelect === '') {
            setFormSelect(formItems[1].key)
        }

        if (
            donators.length !== 0 &&
            selectedChurch &&
            selectedChurch.id
        ) {
            const dd = donators
            if (dd && !dd.find(donator => donator.donations)) {
                churchDonationsLoad(selectedChurch.id)
            }

        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedChurch?.id]);

    const donationsAppButton = () => {
        return (
            <div className='apphome-container' onClick={() => onClick()}>
                <AppButton
                    name='Donations'
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
                            <DonationsTable
                                church={selectedChurch}
                                donators={donatorsMap}
                                donationReport={formSelect === 'form-item-2' && churchDonationReport ? churchDonationReport : null}
                                statementReport={formSelect === 'form-item-3' && donationStatementReport ? donationStatementReport : null}
                            /> :
                            null
                    }

                </div>
            </div>
        )
    }

    const formItems: { [index: number]: { text: string, key: string, onSelect: Function } } = {
        1: { text: 'New Donation', key: 'form-item-1', onSelect: setFormSelect },
        2: { text: 'New Church Donation Report', key: 'form-item-2', onSelect: setFormSelect },
        3: { text: 'New Church Donator Statement Report', key: 'form-item-3', onSelect: setFormSelect },
        4: { text: 'Add New Donation Category', key: 'form-item-4', onSelect: setFormSelect },
        5: { text: 'Remove Donation Category', key: 'form-item-5', onSelect: setFormSelect },
    }
    const forms: {
        [form: string]: {
            [field: string]: { type: string, name: string, options?: Array<any>, action?: Function }
        }
    } = {
        // New Donation
        [formItems[1].key]: {
            'donatorname': {
                name: 'Donator Name',
                type: 'autocomplete',
                options: donators.length !== 0 ? donators.map(d => {
                    return {
                        id: 'donatorid',
                        value: d?.id || '',
                        label: (d?.firstname || '') + ' ' + (d?.lastname || '')
                    }
                }) : []
            },
            'date': { name: 'Date', type: 'datetime-local' },
            'amount': { name: 'Amount', type: 'number' },
            'type': {
                name: 'Type', type: 'autocomplete', options: Object.values(DONATION_TYPES.DonationTypes).map(v => {
                    return { id: 'type', value: v, label: v }
                })
            },
            'category': {
                name: 'Category', type: 'autocomplete', options: Object.values(selectedChurch?.donationcategories || {}).map(v => {
                    return { id: 'category', value: v, label: v }
                })
            },
            'details': { name: 'Details', type: 'text' }
        },
        // New Church Donation Report
        [formItems[2].key]: {
            'reporttype': {
                name: 'Report', type: 'autocomplete', options: Object.values(CHURCH_TYPES.ReportType).map(v => {
                    return { id: 'reporttype', value: v, label: v }
                })
            },
            'lower': { name: 'from', type: 'datetime-local' },
            'upper': { name: 'to', type: 'datetime-local' },
            'timeperiod': {
                name: 'donations summed by', type: 'autocomplete', options: Object.values(CHURCH_TYPES.SumFilter).map(v => {
                    return { id: 'timeperiod', value: v, label: v }
                })
            },
            'category': {
                name: 'from', type: 'autocomplete', options: Object.values(Object.assign({ all: "all categories" }, selectedChurch?.donationcategories)).map(v => {
                    return { id: 'category', value: v, label: v }
                })
            }
        },
        // New Donation Statement Report
        [formItems[3].key]: {
            'donatorname': {
                name: 'Donation report for ',
                type: 'autocomplete',
                options: donators.length !== 0 ? [
                    { id: 'donatorid', value: 'all donators', label: "all donators" },
                    ...donators.map(d => {
                        return {
                            id: 'donatorid',
                            value: d?.id || '',
                            label: (d?.firstname || '') + ' ' + (d?.lastname || '')
                        }
                    })
                ] : []
            },
            'lower': { name: 'from', type: 'datetime-local' },
            'upper': { name: 'to', type: 'datetime-local' }
        },
        // Add New Donation Category
        [formItems[4].key]: { 'category': { name: 'Category', type: 'text' } },
        // Remove Donation Category
        [formItems[5].key]: {
            'category': {
                name: 'Category',
                type: 'autocomplete',
                options: Object.values(selectedChurch?.donationcategories || {}).map(v => {
                    return { id: 'category', value: v, label: v }
                })
            }
        }

    }
    const formsActions: { [form: string]: Function } = {
        // New Donation
        [formItems[1].key]: () => {
            const donation: DONATION_TYPES.IDonation = {
                donatorid: parseInt(formDetails?.donatorid),
                churchid: selectedChurch?.id,
                date: formDetails?.date + 'Z',
                amount: parseFloat(formDetails?.amount),
                type: formDetails?.type,
                category: formDetails?.category,
                currency: 'US',
                details: formDetails?.details,
            }
            dispatch(APP_ACTIONS.createDonationACT(donation))

            console.log('Test action New Donation', donation, ' <---> ', formDetails)
        },
        // New Church Donation Report
        [formItems[2].key]: () => {

            const churchDonationReport: CHURCH_TYPES.IChurchDonationReport = {
                churchid: selectedChurch?.id || 0,
                reporttype: formDetails?.reporttype,
                timerange: { upper: formDetails?.upper + 'Z', lower: formDetails?.lower + 'Z' },
                sumfilter: { timeperiod: formDetails?.timeperiod, multiplier: 1 },
                donationcategories: formDetails?.category === "all categories" && selectedChurch?.donationcategories ? Object.values(selectedChurch.donationcategories) : [formDetails?.category]
            }
            dispatch(APP_ACTIONS.getChurchDonationReport(churchDonationReport))

            console.log('Test action New Church Donation Report', formDetails)
        },
        // New Church Donator Statement Report    
        [formItems[3].key]: () => {
            const donationStatementReport: CHURCH_TYPES.IDonationStatementReport = {
                churchid: selectedChurch?.id || 0,
                timerange: { upper: formDetails?.upper + 'Z', lower: formDetails?.lower + 'Z' },
                donatorids: formDetails?.donatorid === "all donators" && donators.length !== 0 ? donators.map(donator => donator.id || 0) : [parseInt(formDetails?.donatorid)]
            }
            dispatch(APP_ACTIONS.getDonationStatementReport(donationStatementReport))

            console.log('Test action New Donation Statement Report', donationStatementReport, formDetails)
        },
        // Add New Donation Category
        [formItems[4].key]: () => {
            const church: any = { id: selectedChurch?.id }
            if (selectedChurch?.donationcategories) {
                const donationCategories = selectedChurch?.donationcategories
                donationCategories[formDetails?.category] = formDetails?.category
                church.donationcategories = donationCategories
            } else {
                church.donationcategories = { [formDetails?.category]: formDetails?.category }
            }
            dispatch(APP_ACTIONS.updateChurchACT(church))

            console.log('Test action Add New Donation Category', formDetails)
        },
        // Remove Donation Category
        [formItems[5].key]: () => {
            const church: any = { id: selectedChurch?.id }
            if (selectedChurch?.donationcategories) {
                const donationCategories = selectedChurch?.donationcategories
                delete donationCategories[formDetails?.category]
                church.donationcategories = donationCategories
            }

            dispatch(APP_ACTIONS.updateChurchACT(church))

            console.log('Test action Remove Donation Category', formDetails)
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
                        donationsAppButton() :
                        null

            }
        </Fragment>

    )
}
