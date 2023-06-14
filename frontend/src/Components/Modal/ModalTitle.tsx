import * as React from 'react';

import './Modal.css';


export interface ModalTitleProps
{
    children?:  React.ReactNode | undefined;
}


export class ModalTitle extends React.Component<ModalTitleProps>
{
    public render ( )
    {
        return ( 
            <div className='ModalTitle'>
                {this.props.children}
            </div>
        );
    }
}
