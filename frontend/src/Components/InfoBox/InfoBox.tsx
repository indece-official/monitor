import React from 'react';
import { faInfoCircle } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import './InfoBox.css';


export interface InfoBoxProps
{
    message: string | null;
}


export class InfoBox extends React.Component<InfoBoxProps>
{
    public render ( )
    {
        if ( ! this.props.message )
        {
            return null;
        }

        return (
            <div className='InfoBox'>
                <FontAwesomeIcon icon={faInfoCircle} />

                <div className='InfoBox-message'>
                    {this.props.message}
                </div>
            </div>
        );
    }
}
