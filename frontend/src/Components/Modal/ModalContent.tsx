import * as React from 'react';

import './Modal.css';


export interface ModalContentProps
{
    children?:  React.ReactNode | undefined;
}


export class ModalContent extends React.Component<ModalContentProps>
{
    public render ( )
    {
        return ( 
            <div className='ModalContent'>
                {this.props.children}
            </div>
        );
    }
}
