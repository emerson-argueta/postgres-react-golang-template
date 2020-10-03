import React from 'react'
import { List, ListItemText, ListItem } from '@material-ui/core'

export const FormSelect = (
    { formItems }:
    {
      formItems: {[index:number]:{text:string,key:string,onSelect:Function}},
    
    }
) => {

    const [selectedIndex, setSelectedIndex] = React.useState(1);

    const handleFormItemClick = (
        event: React.MouseEvent<HTMLDivElement, MouseEvent>,
        index: number,
        formItem:{text:string,key:string,onSelect:Function}
    ) => {
        setSelectedIndex(index);
        formItem.onSelect(formItem.key)
    };

    const renderFormItems = () => {
        return Object.entries(formItems).map((formItemEntry) => {
            const index = parseInt(formItemEntry[0])
            const formItem = formItemEntry[1]
            return <ListItem
                onClick={(e)=>{
                    
                    handleFormItemClick(e, index,formItem)
                }}
                key={formItem.key}
                button
                selected={selectedIndex === index}
            >
                <ListItemText primary={formItem.text} />
            </ListItem>
        })
    }
    return (
        <div>
            <List>
                {renderFormItems()}
            </List>
        </div>
    )
}
