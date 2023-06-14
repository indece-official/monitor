import { Field, FieldProps } from 'formik';
import * as React from 'react';

import './Input.css';


export interface InputDateTimeProps
{
    label:      string;
    name:       string;
    disabled?:  boolean;
    required?:  boolean;
}


interface InputDateTimeState
{
    focussed:   boolean;
}


export class InputDateTime extends React.Component<InputDateTimeProps, InputDateTimeState>
{
    constructor ( props: InputDateTimeProps )
    {
        super(props);

        this.state = {
            focussed:   false
        };

        this._focus = this._focus.bind(this);
        this._blur = this._blur.bind(this);
    }


    private _focus ( ): void
    {
        this.setState({
            focussed: true
        });
    }
    
    
    private _blur ( ): void
    {
        this.setState({
            focussed: false
        });
    }


    public render ( )
    {
        return (
            <Field name={this.props.name}>
                {( fieldProps: FieldProps<string> ) => (
                    <label className={`InputDateTime ${fieldProps.field.value ? '' : 'empty'} ${this.props.disabled ? '' : 'disabled'} ${this.state.focussed ? 'focussed' : ''} ${fieldProps.meta.error ? 'error' : ''}`}>
                        <div className='InputDateTime-label'>{this.props.label} {this.props.required ? '*' : ''}</div>

                        <input
                            {...fieldProps.field}
                            type='datetime-local'
                            placeholder=''
                            onFocus={this._focus}
                            onBlur={this._blur}
                        />

                        {fieldProps.meta.error ?
                            <div className='InputDateTime-error'>
                                {fieldProps.meta.error}
                            </div>
                        : null}
                    </label>
                )}
            </Field>
        );
    }
}
