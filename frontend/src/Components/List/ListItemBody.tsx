import React from 'react';

import './List.css';


export interface ListItemBodyProps
{
    children?:  React.ReactNode;
}


export class ListItemBody extends React.Component<ListItemBodyProps>
{
    public render ( )
    {
        return (
            <div className='ListItemBody'>
                {this.props.children}
            </div>
        );
    }
}
