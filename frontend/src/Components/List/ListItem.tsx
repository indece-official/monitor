import React from 'react';

import './List.css';


export enum ListItemColor
{
    Red     = 'red',
    Green   = 'green'
}


export interface ListItemProps
{
    color?:     ListItemColor;
    children?:  React.ReactNode;
}


export class ListItem extends React.Component<ListItemProps>
{
    public render ( )
    {
        return (
            <div className={`ListItem ${this.props.color ? 'color-' + this.props.color : ''}`}>
                {this.props.children}
            </div>
        );
    }
}
