import * as React from 'react';

import './LabelValueList.css';


export interface LabelValueListProps
{
    children?:  React.ReactNode | null;
}


export class LabelValueList extends React.Component<LabelValueListProps>
{
    public render ( )
    {
        return (
            <div className='LabelValueList'>
                {this.props.children}
            </div>
        );
    }
}
