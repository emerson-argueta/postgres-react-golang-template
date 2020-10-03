import React, { useState } from 'react'
import MaterialTable, { MTableEditRow, MTableAction } from 'material-table'

export const TestComp = (
    {
        colInVisible,
        onEditClick,
        
    }:{
        colInVisible:boolean,
        onEditClick:Function
    }
) => {
    const [isHidden,setIsHidden] = useState(true)
    const data = [
        {
            id: 5,
            name: "e",
            surname: "Baran",
            birthYear: 1987,
            birthCity: 34,
            sex: "Female",
            type: "child"
        },
        {
            id: 6,
            name: "f",
            surname: "Baran",
            birthYear: 1987,
            birthCity: 34,
            sex: "Female",
            type: "child",
            extra: "extra",
            parentId: 5
        }
    ];

    const testTable = () => {
        return (
            <MaterialTable
                title="Basic Tree Data Preview"
                data={data}
                components={{
                    Groupbar: props => (
                        null
                    ),
                    
                    EditRow: props => {
                        console.log(props)
                        props.columns = 
                        alert('hello everyone here!!!');
                        // setIsHidden(false)
                        // onEditClick()
                        return (
                            
                                <MTableEditRow
                                    {...props}
                                            
                                />
                            
                        )
                    }
                }
                }
                columns={[
                    { title: "Adı", field: "name", hidden: false },
                    { title: "Soyadı", field: "surname", hidden: isHidden },
                    { title: "Cinsiyet", field: "sex", hidden: isHidden },
                    { title: "Tipi", field: "type", removable: false, hidden: isHidden },
                    { title: "Doğum Yılı", field: "birthYear", type: "numeric", hidden: isHidden },
                    {
                        
                        title: "Doğum Yeri",
                        field: "birthCity",
                        lookup: { 34: "İstanbul", 63: "Şanlıurfa" },
                        hidden: colInVisible
                    }
                ]}
                //   parentChildData={(row, rows) => rows.find(a => a.id === row.parentId)}
                options={{
                    // selection: true,
                    // grouping:true
                }}
                //   onTreeExpandChange={()=>{setColInVisible(!colInVisible)}}  
                editable={{
                    onRowUpdate: (newData, oldData) => {
                        return new Promise((resolve, reject) => {
                            setTimeout(() => {
                                const dataUpdate = [...data];
                                const index = oldData?.id || 0;
                                dataUpdate[index] = newData;
                                //   setData([...dataUpdate]);

                                resolve();
                            }, 1000);
                        })
                    }
                }}
            />
        );
    }
    return (
        testTable()
    )
}
