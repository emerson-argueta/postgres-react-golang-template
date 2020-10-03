import React from 'react'
import { IAdministrator, AccessTypes, RoleTypes } from '../../../types/Interface/AdministratorInterfaces'
import MaterialTable, { Column } from 'material-table';
import { useDispatch } from 'react-redux';
import { updateChurchAdministratorACT, deleteChurchAdministratorACT } from '../../../redux/actions/AppActions';
import { IChurch } from '../../../types/Interface/ChurchInterfaces';


export const AdministratorsTable = (
  { church, administrators }:
    {
      church: IChurch,
      administrators: Array<IAdministrator> | null
    }
) => {
  const dispatch = useDispatch()
  
  let data: Array<any> = []
  if (administrators) {
    data = administrators.map((a) => {
      return {
        email: a.email,
        name: a.firstname + ' ' + a.lastname,
        access: a.access,
        role: a.role,
        uuid:a.uuid
      }

    })
  }
  const columns: Array<Column<any>> = [
    { title: 'Email', field: 'email', editable: 'never' },
    { title: 'Name', field: 'name', editable: 'never' },
    { title: 'Access', field: 'access', lookup: AccessTypes },
    { title: 'Role', field: 'role', lookup: RoleTypes, editable: 'never' },

  ]
  return (
    <div>
      <MaterialTable
        title={(church?.name || 'No') + ' Administrators'}
        columns={columns}
        data={data}
        options={{
          exportButton:true,
          exportAllData:true
        }}
        editable={{
          onRowUpdate: (newAdministrator, oldAdministrator) =>
            new Promise((resolve) => {
              setTimeout(() => {
                resolve();
                if (oldAdministrator) {
                  dispatch(updateChurchAdministratorACT(church.id||0,newAdministrator.uuid,{access:newAdministrator.access,role:newAdministrator.role}))
                }
              }, 600);
            }),
          onRowDelete: oldAdministrator =>
            new Promise((resolve, reject) => {
                setTimeout(() => {
                    resolve();
                    dispatch(deleteChurchAdministratorACT(church.id||0,oldAdministrator.uuid))
                }, 1000);
            })
        }}
      />
    </div>
  )
}
