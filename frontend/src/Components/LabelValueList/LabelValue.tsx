import * as React from 'react';

import './LabelValueList.css';


export interface LabelValueProps
{
    label:      string;
    value?:     string;
    children?:  React.ReactNode;
}


export class LabelValue extends React.Component<LabelValueProps>
{
    public render ( )
    {
        return (
            <div className='LabelValue'>
                <div className='LabelValue-label'>
                    {this.props.label}
                </div>

                <div className='LabelValue-value'>
                    {this.props.value ? this.props.value : this.props.children}
                </div>
            </div>
        );
    }
}
