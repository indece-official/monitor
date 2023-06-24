import React from 'react';
import { Field, FieldProps } from 'formik';
import { InputSelectOption } from './InputSelect';

import './Input.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTimes } from '@fortawesome/free-solid-svg-icons';


export interface InputMultiSelectProps
{
    label:      string;
    name:       string;
    disabled?:  boolean;
    required?:  boolean;
    options:    Array<InputSelectOption>;
}


interface InputMultiSelectState
{
    focussed:   boolean;
}


export class InputMultiSelect extends React.Component<InputMultiSelectProps, InputMultiSelectState>
{
    constructor ( props: InputMultiSelectProps )
    {
        super(props);

        this.state = {
            focussed:   false
        };

        this._validate = this._validate.bind(this);
        this._focus = this._focus.bind(this);
        this._blur = this._blur.bind(this);
    }


    private _validate ( value: Array<string> ): string | null
    {
        value = (value || []).map( s => s.trim() ).filter( s => !!s );

        if ( this.props.required && (!value || value.length === 0) )
        {
            return `Please select ${this.props.label}`;
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
                {( fieldProps: FieldProps<Array<string>> ) => (
                    <label className={`InputMultiSelect ${fieldProps.field.value ? '' : 'empty'} ${this.props.disabled ? '' : 'disabled'} ${this.state.focussed ? 'focussed' : ''} ${fieldProps.meta.error ? 'error' : ''}`}>
                        <div className='InputMultiSelect-field'>
                            <div className='InputMultiSelect-label'>{this.props.label} {this.props.required ? '*' : ''}</div>
                            
                            {/*<select
                                {...fieldProps.field}
                                placeholder=''
                                multiple={true}
                                onFocus={this._focus}
                                onBlur={this._blur}>
                                <option value='' key=''></option>

                                {this.props.options.map( ( option ) => (
                                    <option value={option.value} key={option.value}>{option.label}</option>
                                ))}
                            </select>*/}
                            <div className='InputMultiSelect-input'>
                                {(fieldProps.field.value || []).map( ( value, i ) => (
                                    <div className='InputMultiSelect-input-value' key={i}>
                                        <div className='InputMultiSelect-input-value-label'>
                                            {this.props.options.find( o => o.value === value )?.label || value }
                                        </div>

                                        <div
                                            className='InputMultiSelect-input-value-action'
                                            onClick={ ( ) => fieldProps.form.setFieldValue(fieldProps.field.name, (fieldProps.field.value || []).filter( s => s !== value))}>
                                            <FontAwesomeIcon icon={faTimes} />
                                        </div>
                                    </div>
                                ))}

                                <div className='InputMultiSelect-input-select'>
                                    <input
                                        type='text'
                                    />

                                    <div className='InputMultiSelect-dropdown'>
                                        {this.props.options.filter( o => !(fieldProps.field.value || []).includes(o.value) ).map( ( option ) => (
                                            <div
                                                className='InputMultiSelect-dropdown-option'
                                                key={option.value}
                                                onClick={ ( ) => fieldProps.form.setFieldValue(fieldProps.field.name, [...(fieldProps.field.value || []), option.value])}>
                                                {option.label}
                                            </div>
                                        ))}

                                        {this.props.options.filter( o => !(fieldProps.field.value || []).includes(o.value) ).length === 0 ?
                                            <div className='InputMultiSelect-dropdown-empty'>
                                                No options found
                                            </div>
                                        : null}
                                    </div>
                                </div>
                            </div>
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
