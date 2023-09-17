import { Field, FieldProps } from 'formik';
import * as React from 'react';

import './Input.css';


export interface InputTextProps
{
    label:      string;
    name:       string;
    type?:      'text' | 'password' | 'email' | 'color';
    disabled?:  boolean;
    required?:  boolean;
}


interface InputTextState
{
    focussed:   boolean;
}


export class InputText extends React.Component<InputTextProps, InputTextState>
{
    constructor ( props: InputTextProps )
    {
        super(props);

        this.state = {
            focussed:   false
        };

        this._validate = this._validate.bind(this);
        this._focus = this._focus.bind(this);
        this._blur = this._blur.bind(this);
    }


    private _validate ( value: string ): string | null
    {
        value = (value || '').trim();

        if ( this.props.required && !value )
        {
            return `Bitte ${this.props.label} eingeben`;
        }

        return null;
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
            <Field
                name={this.props.name}
                validate={this._validate}>
                {( fieldProps: FieldProps<string> ) => (
                    <label className={`InputText ${fieldProps.field.value ? '' : 'empty'} ${this.props.disabled ? '' : 'disabled'} ${this.state.focussed ? 'focussed' : ''} ${fieldProps.meta.error ? 'error' : ''}`}>
                        <div className='InputText-field'>
                            <div className='InputText-label'>{this.props.label} {this.props.required ? '*' : ''}</div>

                            <input
                                {...fieldProps.field}
                                type={this.props.type || 'text'}
                                placeholder=''
                                onFocus={this._focus}
                                onBlur={this._blur}
                            />
                        </div>

                        {fieldProps.meta.error ?
                            <div className='InputText-error'>
                                {fieldProps.meta.error}
                            </div>
                        : null}
                    </label>
                )}
            </Field>
        );
    }
}
