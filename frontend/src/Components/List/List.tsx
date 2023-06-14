import React from 'react';

import './List.css';


export interface ListProps
{
    children?:  React.ReactNode;
}


export class List extends React.Component<ListProps>
{
    public render ( )
    {
        return (
            <div className='List'>
                {this.props.children}
            </div>
        );
    }
}
