import * as React from 'react';
import { faCheck } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import './Input.css';


export interface InputCheckboxProps
{
    label:      string;
    name:       string;
    value:      boolean;
    disabled?:  boolean;
    error?:     string | boolean;
    onChange?:  ( evt: React.ChangeEvent<HTMLInputElement> ) => any;
    children?:  React.ReactNode | undefined;
}


export class InputCheckbox extends React.Component<InputCheckboxProps>
{
    public render ( )
    {
        return (
            <label className={`InputCheckbox ${this.props.disabled ? 'disabled' : ''} ${this.props.error ? 'error' : ''}`}>
                <input
                    type='checkbox'
                    name={this.props.name}
                    value={1}
                    disabled={this.props.disabled}
                    checked={this.props.value}
                    onChange={this.props.onChange}
                />

                <div className='InputCheckbox-mark'>
                    <FontAwesomeIcon icon={faCheck} />
                </div>

                <div className='InputCheckbox-label'>
                    {this.props.label} {this.props.children}
                </div>
            </label>
        );
    }
}
