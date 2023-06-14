import * as React from 'react';
import { ErrorTranslator } from '../../utils/ErrorTranslator';

import './ErrorBox.css';
import { faWarning } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';


export interface ErrorBoxProps
{
    error:      Error | null;
    className?: string;
}


export class ErrorBox extends React.Component<ErrorBoxProps>
{
    public render ( )
    {
        if ( ! this.props.error )
        {
            return null;
        }

        return (
            <div className={`ErrorBox ${this.props.className || ''}`}>
                <FontAwesomeIcon icon={faWarning} />

                <div className='ErrorBox-text'>
                    {ErrorTranslator.translate(this.props.error)}
                </div>
            </div>
        );
    }
}
