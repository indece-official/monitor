import React from 'react';

import './List.css';


export interface ListItemHeaderProps
{
    children?:  React.ReactNode;
}


export class ListItemHeader extends React.Component<ListItemHeaderProps>
{
    public render ( )
    {
        return (
            <div className='ListItemHeader'>
                {this.props.children}
            </div>
        );
    }
}
