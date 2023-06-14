import React from 'react';
import { faCheck } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import './SuccessBox.css';


export interface SuccessBoxProps
{
    message: string | null;
}


export class SuccessBox extends React.Component<SuccessBoxProps>
{
    public render ( )
    {
        if ( ! this.props.message )
        {
            return null;
        }

        return (
            <div className='SuccessBox'>
                <FontAwesomeIcon icon={faCheck} />
                <div className='SuccessBox-message'>
                    {this.props.message}
                </div>
            </div>
        );
    }
}
