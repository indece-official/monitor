import * as React from 'react';

import './Tag.css';


export interface TagProps
{
    color:      string;
    children?:  React.ReactNode | null;
}


export class Tag extends React.Component<TagProps>
{
    public render ( )
    {
        return (
            <div className="Tag" style={{backgroundColor: this.props.color}}>
                {this.props.children}
            </div>
        );
    }
}
