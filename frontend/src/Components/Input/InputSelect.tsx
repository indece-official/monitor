import React from 'react';
import { Field, FieldProps } from 'formik';

import './Input.css';


export interface InputSelectOption
{
    label: string;
    value: string;
}


export interface InputSelectProps
{
    label:      string;
    name:       string;
    disabled?:  boolean;
    required?:  boolean;
    options:    Array<InputSelectOption>;
}


interface InputSelectState
{
    focussed:   boolean;
}


export class InputSelect extends React.Component<InputSelectProps, InputSelectState>
{
    constructor ( props: InputSelectProps )
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
            return `Bitte ${this.props.label} wählen`;
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
                    <label className={`InputSelect ${fieldProps.field.value ? '' : 'empty'} ${this.props.disabled ? '' : 'disabled'} ${this.state.focussed ? 'focussed' : ''} ${fieldProps.meta.error ? 'error' : ''}`}>
                        <div className='InputText-field'>
                            <div className='InputSelect-label'>{this.props.label} {this.props.required ? '*' : ''}</div>
                            
                            <select
                                {...fieldProps.field}
                                placeholder=''
                                onFocus={this._focus}
                                onBlur={this._blur}>
                                <option value='' key=''></option>

                                {this.props.options.map( ( option ) => (
                                    <option value={option.value} key={option.value}>{option.label}</option>
                                ))}
                            </select>
                        </div>

                        {fieldProps.meta.error ?
                            <div className='InputSelect-error'>
                                {fieldProps.meta.error}
                            </div>
                        : null}
                    </label>
                )}
            </Field>
        );
    }
}
