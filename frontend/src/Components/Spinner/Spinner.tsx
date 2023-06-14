import * as React from 'react';

import './Spinner.css';


export enum SpinnerColor
{
    Default     = 'default',
    White       = 'white'
}


export interface SpinnerProps
{
    active?:    boolean;
    color?:     SpinnerColor;
}


export class Spinner extends React.Component<SpinnerProps>
{
    public render ( )
    {
        if ( this.props.active === false )
        {
            return null;
        }

        return (
            <div className={`Spinner color-${this.props.color || SpinnerColor.Default}`}>
                <div />
                <div />
            </div>
        );
    }
}
