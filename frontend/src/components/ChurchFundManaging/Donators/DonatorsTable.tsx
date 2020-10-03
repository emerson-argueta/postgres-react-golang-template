import React from 'react'
import { IDonator } from '../../../types/Interface/DonatorInterfaces'
import MaterialTable, { Column } from 'material-table'
import { useDispatch } from 'react-redux'
import { updateDonatorACT, deleteDonatorACT } from '../../../redux/actions/AppActions'
import { IChurch } from '../../../types/Interface/ChurchInterfaces'

export const DonatorsTable = (
    { church, donators }:
        {
            church: IChurch,
            donators: Array<IDonator> | null
        }
) => {

    const dispatch = useDispatch()

    let data: Array<any> = []
    if (donators) {
        data = donators.map((d) => {
            return {
                id: d?.id || null,
                firstname: d?.firstname || null,
                lastname: d?.lastname || null,
                email: d?.email || null,
                address: d?.address || null,
                phone: d?.phone || null

            }

        })
    }
    const columns: Array<Column<any>> = [
        { title: 'Firstname', field: 'firstname' },
        { title: 'Lastname', field: 'lastname' },
        { title: 'Email', field: 'email' },
        { title: 'Address', field: 'address' },
        { title: 'Phone', field: 'phone' }

    ]
    return (
        <div>
            <MaterialTable
                title={(church?.name||'No')  + ' Donators'}
                columns={columns}
                data={data}
                // options={{
                //   filtering: true
                // }}
                editable={{
                    onRowUpdate: (newDonator, oldDonator) =>
                        new Promise((resolve) => {
                            setTimeout(() => {
                                resolve();
                                if (oldDonator) {
                                    console.log('update donator on lv table --->', newDonator, oldDonator);
                                    if (church.id !== 0) {
                                        dispatch(updateDonatorACT(church.id || 0, newDonator))
                                    }


                                }
                            }, 600);
                        }),
                    onRowDelete: oldDonator =>
                        new Promise((resolve, reject) => {
                            setTimeout(() => {
                                resolve();
                                dispatch(deleteDonatorACT(church.id || 0, oldDonator.id))
                            }, 1000);
                        })
                }}
            />
        </div>
    )
}
