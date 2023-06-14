import { faTimes } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import * as React from 'react';

import './Modal.css';


export interface ModalProps
{
    onClose?:   ( ) => any;
    closable?:  boolean;
    className?: string;
    children?:  React.ReactNode | undefined;
}


export class Modal extends React.Component<ModalProps>
{
    constructor ( props: ModalProps )
    {
        super(props);

        this._keyUp = this._keyUp.bind(this);
    }


    private _keyUp ( evt: KeyboardEvent ): void
    {
        if ( evt.key === 'Escape' && this.props.onClose && this.props.closable !== false )
        {
            this.props.onClose();
        }
    }


    public componentDidMount ( ): void
    {
        window.addEventListener('keyup', this._keyUp);
    }


    public componentWillUnmount ( ): void
    {
        window.removeEventListener('keyup', this._keyUp);
    }


    public render ( )
    {
        return ( 
            <div className={`Modal ${this.props.className || ''}`}>
                <div className='Modal-backdrop' onClick={this.props.onClose} />

                <div className='Modal-inner'>
                    {this.props.children}

                    <div className='Modal-close' onClick={this.props.onClose}>
                        <FontAwesomeIcon icon={faTimes} />
                    </div>
                </div>
            </div>
        );
    }
}
