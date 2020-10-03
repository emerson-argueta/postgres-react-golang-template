import axios from "axios"
import * as APP_TYPES from "../../types/AppTypes"
import * as INTERFACES from "../../types/Interface/Interfaces"
import * as AUTH_ACTIONS from "./AuthActions"
import { Access, Role } from "../../types/Interface/AdministratorInterfaces";
import { IDonation } from "../../types/Interface/DonationInterfaces";
import * as CHURCH_TYPES from "../../types/Interface/ChurchInterfaces";
import { IDonator } from "../../types/Interface/DonatorInterfaces";

// Load administrator using jwt token.
export const userLoadACT = () => (dispatch: Function, getState: Function) => {
    // Administrator loading
    dispatch(AUTH_ACTIONS.userLoadingACT());

    const url = '/api/administrator'
    axios
        .get(url, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.USER_LOADED, payload: res.data })
            if (res.data.administrator.churches) {
                dispatch(churchesLoadACT())
            }
            dispatch(AUTH_ACTIONS.userLoadedACT())
        })
        .catch(err => {
            if (err.response?.status === 401) {
                dispatch(AUTH_ACTIONS.userTokenRefreshACT(userLoadACT()));
                return
            }
            dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
};
export const updateUserACT = (user: any) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator'
    const req = { administrator: { ...user } }
    axios
        .patch(url, req, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.USER_UPDATED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, updateUserACT(user)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}

export const updateChurchAdministratorACT = (
    churchID: number,
    administratorUUID: string,
    churchAdministrator: { access: Access, role: Role },
) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/administrators/' + churchID + '/' + administratorUUID

    let req: any = { administrator: { ...churchAdministrator } }

    axios
        .patch(url, req, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.ADMINISTRATOR_UPDATED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, updateChurchAdministratorACT(churchID, administratorUUID, churchAdministrator)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const deleteChurchAdministratorACT = (
    churchID: number,
    administratorUUID: string
) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/administrators/' + churchID + '/' + administratorUUID

    let req: any = {}
    axios
        .patch(url, req, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.ADMINISTRATOR_DELETED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, deleteChurchAdministratorACT(churchID, administratorUUID)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}

const churchesLoadACT = () => (dispatch: Function, getState: Function) => {
    const url = '/api/administrator/churches'
    axios
        .get(url, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            const selectedChurch: number = res.data ?
                parseInt(Object.keys(res.data)[0]) :
                0

            dispatch(administratorSelectChurchACT(selectedChurch))
            dispatch({ type: APP_TYPES.CHURCHES_LOADED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, churchesLoadACT()))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}

export const administratorSelectChurchACT = (churchID: number) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    dispatch({ type: APP_TYPES.ADMINISTRATOR_SELECT_CHURCH, payload: churchID })
}

export const churchAdministratorsLoadACT = (churchID: number) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/administrators/' + churchID

    axios
        .get(url, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.CHURCH_ADMINISTRATORS_LOADED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, churchAdministratorsLoadACT(churchID)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const churchDonatorsLoadACT = (churchID: number) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/donators/' + churchID

    axios
        .get(url, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.CHURCH_DONATORS_LOADED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, churchDonatorsLoadACT(churchID)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const createDonatorACT = (churchID: number, donator: IDonator) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/donator'
    let req = { donator: { ...donator }, church: { id: churchID } }
    axios
        .post(url, req, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.DONATOR_CREATED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, createDonatorACT(churchID, donator)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const updateDonatorACT = (churchID: number, donator: IDonator) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/donator'
    let req = { donator: { ...donator }, church: { id: churchID } }
    axios
        .patch(url, req, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.DONATOR_UPDATED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, updateDonatorACT(churchID, donator)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const deleteDonatorACT = (churchID: number, donatorID: number) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/donator/' + churchID + '/' + donatorID

    axios
        .delete(url, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.DONATOR_DELETED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, deleteDonatorACT(churchID, donatorID)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const churchDonationsLoadACT = (churchID: number) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/donations/' + churchID

    axios
        .get(url, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.CHURCH_DONATIONS_LOADED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, churchDonationsLoadACT(churchID)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const createDonationACT = (donation: IDonation) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/donation'
    let req = { donation: { ...donation } }
    axios
        .post(url, req, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.DONATION_CREATED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, createDonationACT(donation)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const updateDonationACT = (donation: IDonation) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/donation'
    let req = { donation: { ...donation } }
    axios
        .patch(url, req, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.DONATION_UPDATED, payload: res.data })
            dispatch(churchDonationsLoadACT(donation.churchid || 0))
        })
        .catch(err => {
            dispatch(handleError(err, updateDonationACT(donation)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}

export const createChurchACT = (church: any) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/church'
    let req = { church: { ...church } }
    axios
        .post(url, req, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.CHURCH_CREATED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, createChurchACT(church)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const addChurchACT = (email: string, password: string) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/church/add'
    let req = { church: { email: email, password: password } }
    axios
        .post(url, req, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.CHURCH_ADDED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, addChurchACT(email, password)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const updateChurchACT = (church: CHURCH_TYPES.IChurch) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/church'
    let req = { church: { ...church } }
    axios
        .patch(url, req, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.CHURCH_UPDATED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, updateChurchACT(church)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const deleteChurchACT = (churchID: number) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const url = '/api/administrator/church/' + churchID

    axios
        .delete(url, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.CHURCH_DELETED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, deleteChurchACT(churchID)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const getChurchDonationReport = (churchDonationReport: CHURCH_TYPES.IChurchDonationReport) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const multiplier = churchDonationReport.sumfilter.multiplier
    const categories = churchDonationReport.donationcategories.reduce((acc, cur) => {
        acc += "categories=" + cur + "&"
        return acc
    }, "")

    const url = '/api/administrator/donations/report/' + churchDonationReport.churchid + '/'
        + churchDonationReport.timerange.lower + '/' + churchDonationReport.timerange.upper + '/'
        + churchDonationReport.sumfilter.timeperiod + '/' + multiplier + '?' + categories

    axios
        .get(url, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.CHURCH_REPORT_LOADED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, getChurchDonationReport(churchDonationReport)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}
export const getDonationStatementReport = (donationStatementReport: CHURCH_TYPES.IDonationStatementReport) => (dispatch: Function, getState: () => { auth: INTERFACES.IAuthState, app: INTERFACES.IAppState }) => {
    const donatorIDs = donationStatementReport.donatorids.reduce((acc, cur) => {
        acc += "donatorIDs=" + cur + "&"
        return acc
    }, "")

    const url = '/api/administrator/donations/statement/' + donationStatementReport.churchid + '/'
        + donationStatementReport.timerange.lower + '/' + donationStatementReport.timerange.upper + '?' + donatorIDs

    axios
        .get(url, AUTH_ACTIONS.tokenConfig())
        .then(res => {
            dispatch({ type: APP_TYPES.DONATION_STATEMENT_REPORT_LOADED, payload: res.data })
        })
        .catch(err => {
            dispatch(handleError(err, getDonationStatementReport(donationStatementReport)))
            // dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
        });
}

export const clearErrorACT = () => (dispatch: Function, getState: () => { app: INTERFACES.IAppState, auth: INTERFACES.IAuthState }) => {
    dispatch({ type: APP_TYPES.CLEAR_ERROR })
}
const handleError = (err: any, retryAction: Function) => (dispatch: Function, getState: () => { app: INTERFACES.IAppState, auth: INTERFACES.IAuthState }) => {
    if (err.response?.status === 401) {
        console.log("retrying action");

        dispatch(AUTH_ACTIONS.userTokenRefreshACT(retryAction));
        return
    }
    dispatch({ type: APP_TYPES.ERROR, payload: { ...err?.response, id: APP_TYPES.ERROR } })
}