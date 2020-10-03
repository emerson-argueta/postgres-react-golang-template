import React from 'react'
import Autocomplete from '@material-ui/lab/Autocomplete';
import { IChurch } from '../../../types/Interface/ChurchInterfaces';
import { useSelector, useDispatch } from 'react-redux';
import { IAuthState, IAppState } from '../../../types/Interface/Interfaces';
import { TextField } from '@material-ui/core';
import { administratorSelectChurchACT } from '../../../redux/actions/AppActions';

export const ChurchSelect = () => {
    const dispatch = useDispatch()
    const administratorSelectChurch = (churchID: number) => dispatch(administratorSelectChurchACT(churchID))


    const options: Array<IChurch> | [] = useSelector(
        (state: { auth: IAuthState, app: IAppState }) => {
            if (state.app.administrator.churches) {
                return Object.values(state.app.administrator.churches)
            }
            return []
        }
    )
    const selectedChurch: number = useSelector(
        (state: { auth: IAuthState, app: IAppState }) => state.app.selectedChurch
    )

    return (
        <Autocomplete
            value={options.find((o) => o.id === selectedChurch)}
            id="combo-box-demo"
            options={options}
            getOptionLabel={(option) => option.name || ''}
            onChange={(event, newValue) => {
                if (newValue && newValue.id)
                    administratorSelectChurch(newValue.id)
            }}
            style={{ width: 300 }}
            renderInput={(params) => <TextField
                {...params}
                label={'Select church'}
                variant="outlined"
            />}
        />
    )
}
