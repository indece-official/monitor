import * as React from 'react';
import { faCheck } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Field, FieldProps } from 'formik';

import './Input.css';


export interface InputCheckboxProps
{
    label:      string;
    name:       string;
    disabled?:  boolean;
    children?:  React.ReactNode | undefined;
}


export class InputCheckbox extends React.Component<InputCheckboxProps>
{
    public render ( )
    {
        return (
            <Field name={this.props.name}>
                {( fieldProps: FieldProps<boolean> ) => (
                    <label className={`InputCheckbox ${this.props.disabled ? 'disabled' : ''} ${fieldProps.meta.error ? 'error' : ''}`}>
                        <div className='InputCheckbox-field'>
                            <input
                                {...fieldProps.field}
                                type='checkbox'
                                value={1}
                                disabled={this.props.disabled}
                                checked={fieldProps.field.value}
                            />

                            <div className='InputCheckbox-mark'>
                                <FontAwesomeIcon icon={faCheck} />
                            </div>

                            <div className='InputCheckbox-label'>
                                {this.props.label} {this.props.children}
                            </div>
                        </div>

                        {fieldProps.meta.error ?
                            <div className='InputCheckbox-error'>
                                {fieldProps.meta.error}
                            </div>
                        : null}
                    </label>
                )}
            </Field>
        );
    }
}
