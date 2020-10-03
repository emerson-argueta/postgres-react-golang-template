import React, { Fragment, useEffect, useState } from 'react'
import { IAuthState, IAppState } from '../../types/Interface/Interfaces'
import { useSelector, useDispatch } from 'react-redux'
import { userLoadACT } from '../../redux/actions/AppActions'
import { Administrators } from './Administrators/Administrators'
import { Donators } from './Donators/Donators'
import { Donations } from './Donations/Donations'
import { Churches } from './Churches/Churches'


export const ChurchFundManaging = () => {
    const [menuItemSelected, setMenuItemSelected] = useState(false)

    const isAuthenticated: boolean = useSelector((state: { auth: IAuthState, app: IAppState }) => { return state.auth.isAuthenticated })

    const dispatch = useDispatch()
    const userLoad = () => dispatch(userLoadACT())

    useEffect(() => {
        console.log("reloading main menu");

        userLoad();
        return () => {
            setMenuItemSelected(false)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [isAuthenticated]);
    return (
        <Fragment>
            <div className="app-container" >
                {
                    isAuthenticated ?
                        [
                            <Donators key="donators" onClick={() => setMenuItemSelected(!menuItemSelected)} visible={!menuItemSelected} />,
                            <Churches key="churches" onClick={() => setMenuItemSelected(!menuItemSelected)} visible={!menuItemSelected} />,
                            <Donations key="donations" onClick={() => setMenuItemSelected(!menuItemSelected)} visible={!menuItemSelected} />,
                            <Administrators key="adminstrators" onClick={() => setMenuItemSelected(!menuItemSelected)} visible={!menuItemSelected} />
                        ]
                        : null
                }
            </div>


        </Fragment>
    )
}