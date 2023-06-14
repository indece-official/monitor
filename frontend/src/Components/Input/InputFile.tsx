import * as React from 'react';
import { Field, FieldProps } from 'formik';

import './Input.css';


export interface InputFileProps
{
    label:      string;
    name:       string;
    disabled?:  boolean;
    required?:  boolean;
}


export class InputFile extends React.Component<InputFileProps>
{
    public render ( )
    {
        return (
            <Field name={this.props.name}>
                {( fieldProps: FieldProps<File | null> ) => (
                    <label className={`InputFile ${fieldProps.field.value ? '' : 'empty'} ${this.props.disabled ? '' : 'disabled'} ${fieldProps.meta.error ? 'error' : ''}`}>
                        <div className='InputFile-label'>{this.props.label} {this.props.required ? '*' : ''}</div>

                        <input
                            {...fieldProps.field}
                            type='file'
                            value={undefined}
                            onChange={(event) => {
                                fieldProps.form.setFieldValue(fieldProps.field.name, event.currentTarget.files ? event.currentTarget.files[0] : null);
                            }} 
                        />

                        <div className='InputFile-value'>
                            {fieldProps.field.value ? fieldProps.field.value.name : ''}
                        </div>
                    </label>
                )}
            </Field>
        );
    }
}
