import React from 'react';

import './List.css';


export interface ListEmptyProps
{
    children?:  React.ReactNode;
}


export class ListEmpty extends React.Component<ListEmptyProps>
{
    public render ( )
    {
        return (
            <div className='ListEmpty'>
                {this.props.children}
            </div>
        );
    }
}
