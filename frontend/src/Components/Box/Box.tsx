import React from 'react';

import './Box.css';


export interface BoxProps
{
    className?: string;
    children?:  React.ReactNode | undefined;
}


export class Box extends React.Component<BoxProps>
{
    public render ( )
    {
        return (
            <div className={`Box ${this.props.className || ''}`}>
                {this.props.children}
            </div>
        );
    }
}
