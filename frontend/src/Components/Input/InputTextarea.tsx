import { Field, FieldProps } from 'formik';
import * as React from 'react';

import './Input.css';


export interface InputTextareaProps
{
    label:      string;
    name:       string;
    disabled?:  boolean;
    required?:  boolean;
}


interface InputTextareaState
{
    focussed:   boolean;
}


export class InputTextarea extends React.Component<InputTextareaProps, InputTextareaState>
{
    constructor ( props: InputTextareaProps )
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
                    <label className={`InputTextarea ${fieldProps.field.value ? '' : 'empty'} ${this.props.disabled ? '' : 'disabled'} ${this.state.focussed ? 'focussed' : ''} ${fieldProps.meta.error ? 'error' : ''}`}>
                        <div className='InputTextarea-label'>{this.props.label} {this.props.required ? '*' : ''}</div>

                        <textarea
                            {...fieldProps.field}
                            placeholder=''
                            onFocus={this._focus}
                            onBlur={this._blur}>{fieldProps.field.value}</textarea>

                        {fieldProps.meta.error ?
                            <div className='InputTextarea-error'>
                                {fieldProps.meta.error}
                            </div>
                        : null}
                    </label>
                )}
            </Field>
        );
    }
}
