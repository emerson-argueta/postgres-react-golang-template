// import { v4 as uuidv4 } from 'uuid';
import { IAction, IAppState, IError } from "../../types/Interface/Interfaces"
import * as APP_TYPES from "../../types/AppTypes"
import * as CHURCH_TYPES from "../../types/Interface/ChurchInterfaces"
import { IDonator } from "../../types/Interface/DonatorInterfaces"
import { IDonation } from "../../types/Interface/DonationInterfaces"
import { IChurches, IAdministrator } from "../../types/Interface/AdministratorInterfaces"

const initialState: IAppState = {
    administrator: {},
    selectedChurch: 0,
    error: null,
    loading: true
}

export default (state = initialState, action: IAction) => {
    switch (action.type) {

        case APP_TYPES.USER_LOADED: {
            const updatedUser = Object.assign(state.administrator, action.payload.administrator)
            console.log("loaded user");

            return {
                ...state,
                administrator: { ...updatedUser },
                loading: false
            }
        }
        case APP_TYPES.USER_UPDATED: {
            const updatedUser: IAdministrator = Object.assign(state.administrator, { ...action.payload.administrator })
            const updatedChurches = Object.assign(state.administrator.churches || {})

            if (state.administrator.churches) {
                Object.entries(state.administrator.churches).forEach(entry => {
                    const churchID = parseInt(entry[0])
                    const church: CHURCH_TYPES.IChurch = Object.assign({}, entry[1])

                    if (updatedUser.uuid && church.administrators[updatedUser.uuid]) {
                        church.administrators[updatedUser.uuid] = Object.assign(church.administrators[action.payload.administrator.uuid], { ...action.payload.administrator })
                        updatedChurches[churchID] = church
                    }
                })
            }

            return {
                ...state,
                administrator: {
                    ...updatedUser,
                    churches: updatedChurches
                }
            }
        }
        case APP_TYPES.ADMINISTRATOR_DELETED:
        case APP_TYPES.ADMINISTRATOR_UPDATED: {
            if (
                state.administrator.churches &&
                state.administrator.churches[state.selectedChurch] &&
                state.administrator.churches[state.selectedChurch].administrators
            ) {

                const updatedAdministrators = loadChurchAdministrators(
                    state.administrator.churches[state.selectedChurch].administrators,
                    { [action.payload.administratoruuid]: action.payload.administrator }
                )
                const updatedChurches = Object.assign({}, state.administrator.churches)
                updatedChurches[state.selectedChurch].administrators = updatedAdministrators

                const updatedUser = Object.assign(state.administrator, { churches: { ...updatedChurches } })

                return {
                    ...state,
                    administrator: {
                        ...updatedUser
                    }
                }
            }

            return {
                ...state
            }
        }
        case APP_TYPES.CHURCHES_LOADED: {
            const updatedChurches = loadChurches(state.administrator.churches, action.payload)

            return {
                ...state,
                administrator: {
                    ...state.administrator,
                    churches: updatedChurches
                }
            }
        }
        case APP_TYPES.CHURCH_DELETED: {
            const updatedChurches = loadChurches(state.administrator.churches, { [action.payload.churchid]: action.payload.church })
            let selectedChurch: number = state.selectedChurch

            if ((state.selectedChurch === action.payload.churchid) && updatedChurches) {
                selectedChurch = parseInt(Object.keys(updatedChurches).pop() || '0')
            }

            return {
                ...state,
                selectedChurch: selectedChurch,
                administrator: {
                    ...state.administrator,
                    churches: updatedChurches
                }
            }
        }
        case APP_TYPES.CHURCH_UPDATED:
        case APP_TYPES.CHURCH_ADDED:
        case APP_TYPES.CHURCH_CREATED: {
            const updatedChurches = loadChurches(state.administrator.churches, { [action.payload.churchid]: action.payload.church })

            return {
                ...state,
                administrator: {
                    ...state.administrator,
                    churches: updatedChurches
                }
            }
        }
        case APP_TYPES.CHURCH_ADMINISTRATORS_LOADED: {

            if (
                state.administrator.churches &&
                state.administrator.churches[state.selectedChurch] &&
                state.administrator.churches[state.selectedChurch].administrators
            ) {
                const updatedAdministrators = loadChurchAdministrators(
                    state.administrator.churches[state.selectedChurch].administrators,
                    action.payload
                )
                const updateChurches = Object.assign({}, state.administrator.churches)
                updateChurches[state.selectedChurch].administrators = updatedAdministrators

                return {
                    ...state,
                    administrator: {
                        ...state.administrator,
                        churches: updateChurches,
                    }
                }
            }

            return {
                ...state
            }

        }
        case APP_TYPES.CHURCH_DONATORS_LOADED: {
            if (
                state.administrator.churches &&
                state.administrator.churches[state.selectedChurch] &&
                state.administrator.churches[state.selectedChurch].donators

            ) {

                const updatedDonators = loadChurchDonators(
                    state.administrator.churches[state.selectedChurch].donators,
                    action.payload
                )
                const updateChurches = Object.assign({}, state.administrator.churches)
                updateChurches[state.selectedChurch].donators = updatedDonators

                return {
                    ...state,
                    administrator: {
                        ...state.administrator,
                        churches: updateChurches
                    }
                }

            }
            return {
                ...state
            }

        }
        case APP_TYPES.DONATOR_DELETED:
        case APP_TYPES.DONATOR_UPDATED:
        case APP_TYPES.DONATOR_CREATED: {
            const donators = { [action.payload.donatorid]: action.payload.donator }
            if (
                state.administrator.churches &&
                state.administrator.churches[state.selectedChurch] &&
                state.administrator.churches[state.selectedChurch].donators

            ) {

                const updatedDonators = loadChurchDonators(
                    state.administrator.churches[state.selectedChurch].donators,
                    donators
                )
                const updateChurches = Object.assign({}, state.administrator.churches)
                updateChurches[state.selectedChurch].donators = updatedDonators

                return {
                    ...state,
                    administrator: {
                        ...state.administrator,
                        churches: updateChurches
                    }
                }

            }
            return {
                ...state
            }
        }
        case APP_TYPES.CHURCH_DONATIONS_LOADED: {
            if (
                state.administrator.churches &&
                state.administrator.churches[state.selectedChurch] &&
                state.administrator.churches[state.selectedChurch].donators
            ) {
                const updatedDonations = loadAllDonations(
                    state.administrator.churches[state.selectedChurch].donators,
                    action.payload
                )
                const updateChurches = Object.assign({}, state.administrator.churches)
                updateChurches[state.selectedChurch].donators = updatedDonations

                return {
                    ...state,
                    administrator: {
                        ...state.administrator,
                        churches: updateChurches
                    }
                }
            }

            return {
                ...state
            }

        }
        case APP_TYPES.DONATION_CREATED: {
            const donations = { [action.payload.donation.donatorid]: [action.payload.donation] }

            if (
                state.administrator.churches &&
                state.administrator.churches[state.selectedChurch] &&
                state.administrator.churches[state.selectedChurch].donators
            ) {
                const updatedDonations = loadChurchDonations(
                    state.administrator.churches[state.selectedChurch].donators,
                    donations
                )
                const updateChurches = Object.assign({}, state.administrator.churches)
                updateChurches[state.selectedChurch].donators = updatedDonations

                return {
                    ...state,
                    administrator: {
                        ...state.administrator,
                        churches: updateChurches
                    }
                }
            }

            return {
                ...state,
            }
        }
        case APP_TYPES.ADMINISTRATOR_SELECT_CHURCH: {
            return {
                ...state,
                selectedChurch: action.payload,
            }
        }
        case APP_TYPES.CHURCH_REPORT_LOADED: {


            const churchDonationReport: CHURCH_TYPES.IChurchDonationReport = action.payload
            const churchID: number = parseInt(action.payload.churchid)
            if (
                state.administrator.churches &&
                state.administrator.churches[churchID]

            ) {
                console.log("church report loaded reducer", churchID, churchDonationReport);
                const updateChurches = Object.assign({}, state.administrator.churches)
                updateChurches[churchID].churchdonationreport = churchDonationReport

                return {
                    ...state,
                    administrator: {
                        ...state.administrator,
                        churches: updateChurches
                    }
                }

            }
            return {
                ...state
            }
        }
        case APP_TYPES.DONATION_STATEMENT_REPORT_LOADED: {


            const donationStatementReport: CHURCH_TYPES.IDonationStatementReport = action.payload
            const churchID: number = parseInt(action.payload.churchid)
            if (
                state.administrator.churches &&
                state.administrator.churches[churchID]

            ) {
                console.log("donation statement report loaded reducer", churchID, donationStatementReport);
                const updateChurches = Object.assign({}, state.administrator.churches)
                updateChurches[churchID].donationstatementreport = donationStatementReport

                return {
                    ...state,
                    administrator: {
                        ...state.administrator,
                        churches: updateChurches
                    }
                }

            }
            return {
                ...state
            }
        }
        case APP_TYPES.ERROR: {
            const error: IError = { msg: action.payload?.data.error, id: action.payload?.id, status: action.payload?.status }
            return {
                ...state,
                error: error
            }
        }
        case APP_TYPES.CLEAR_ERROR: {
            return {
                ...state,
                error: null
            }
        }
        default:
            return state
    }
}

const loadChurches = (churches: IChurches | undefined, newChurches: IChurches) => {
    if (churches) {
        Object.entries(newChurches).forEach((entry) => {
            const churchID = parseInt(entry[0])
            const newChurch = entry[1]

            if (newChurch && churches[churchID].donators && churches[churchID].administrators) {
                delete newChurch.donators
                delete newChurch.administrators
            }

            if (newChurch) {
                churches[churchID] = Object.assign({ id: churchID }, churches[churchID], newChurch)
            } else {
                delete churches[churchID]
            }

        })
    }


    return churches
}
const loadChurchAdministrators = (administrators: CHURCH_TYPES.IAdministrators, newAdministrators: { [uuid: string]: IAdministrator | null }) => {
    if (administrators) {
        Object.entries(newAdministrators).forEach(entry => {
            const administratorID = entry[0]
            const newAdministrator = entry[1]

            if (newAdministrator) {
                administrators[administratorID] = Object.assign(administrators[administratorID] || {}, { ...newAdministrator })
            } else {
                delete administrators[administratorID]
            }

        })
    }


    return administrators
}

const loadChurchDonators = (donators: CHURCH_TYPES.IDonators | undefined, newDonators: { [id: string]: IDonator }) => {
    if (donators) {
        Object.entries(newDonators).forEach(entry => {
            const donatorID = parseInt(entry[0])
            const newDonator = entry[1]

            if (newDonator) {
                donators[donatorID] = Object.assign(donators[donatorID] || {}, { ...newDonator })
            } else {
                delete donators[donatorID]
            }

        })
    }


    return donators
}

const loadChurchDonations = (donators: CHURCH_TYPES.IDonators | undefined, donations: { [id: string]: Array<IDonation> }) => {
    if (donators) {
        Object.entries(donations).forEach(entry => {
            const donatorID = parseInt(entry[0])
            const newDonations = entry[1]
            const donations = donators[donatorID].donations || []

            if (newDonations) {
                newDonations.forEach(newDonation => { donations.push(newDonation) })
            }
            donators[donatorID] = Object.assign(donators[donatorID], { donations: donations })
        })
    }


    return donators
}
const loadAllDonations = (donators: CHURCH_TYPES.IDonators | undefined, donations: { [id: string]: Array<IDonation> }) => {
    if (donators) {
        Object.entries(donations).forEach(entry => {
            const donatorID = parseInt(entry[0])
            const newDonations = entry[1]

            donators[donatorID] = Object.assign(donators[donatorID] || {}, { donations: newDonations })
        })
    }


    return donators
}