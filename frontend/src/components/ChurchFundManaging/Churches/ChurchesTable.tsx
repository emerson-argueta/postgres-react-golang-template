import React from 'react'
import MaterialTable, { Column } from 'material-table'
import { IChurch } from '../../../types/Interface/ChurchInterfaces'
import { useDispatch } from 'react-redux'
import { updateChurchACT, deleteChurchACT } from '../../../redux/actions/AppActions'


export const ChurchesTable = (
    { administratorName, churches }:
        {
            administratorName: string,
            churches: Array<IChurch> | null
        }
) => {
    const dispatch = useDispatch()
    
    let data: Array<any> = []
    if (churches) {
        data = churches.map((c) => {
            return {
                id: c.id || null,
                name: c.name || null,
                email: c.email || null,
                type: c.type || null,
                address: c.address || null,
                phone: c.phone || null
                
            }

        })
    }
    const columns: Array<Column<any>> = [
        { title: 'Name', field: 'name'},
        { title: 'Email', field: 'email',editable:'never'},
        { title: 'Type', field: 'type'},
        { title: 'Address', field: 'address'},
        { title: 'Phone', field: 'phone'},
        

    ]
    return (
        <div>
            <MaterialTable
                title={administratorName + ' Churches'}
                columns={columns}
                data={data}
                // options={{
                //   filtering: true
                // }}
                editable={{
                    onRowUpdate: (newChurch, oldChurch) =>
                        new Promise((resolve) => {
                            setTimeout(() => {
                                resolve();
                                if (oldChurch) {
                                    dispatch(updateChurchACT(newChurch))
                                }
                            }, 600);
                        }),
                    onRowDelete: oldChurch =>
                        new Promise((resolve, reject) => {
                            setTimeout(() => {
                                resolve();
                                dispatch(deleteChurchACT(oldChurch.id))
                            }, 1000);
                        })
                    
                }}
            />
        </div>
    )
}
