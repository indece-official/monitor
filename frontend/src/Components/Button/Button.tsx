import * as React from 'react';
import { Link, To } from 'react-router-dom';

import './Button.css';


export enum ButtonColor
{
    White   = 'white',
    Primary = 'primary',
    Yellow  = 'yellow'
}


export interface ButtonProps
{
    onClick?:       ( ) => any;
    active?:        boolean;
    className?:     string;
    type?:          'submit' | 'button' | 'reset';
    color?:         ButtonColor;
    disabled?:      boolean;
    title?:         string;
    href?:          string;
    to?:            To;
    target?:        React.HTMLAttributeAnchorTarget;
    children?:      React.ReactNode | undefined;
}


export class Button extends React.Component<ButtonProps>
{
    public render ( )
    {
        if ( this.props.to )
        {
            return (
                <Link
                    className={`Button color-${this.props.color || ButtonColor.Primary} ${this.props.active ? 'active' : ''} ${this.props.disabled ? 'disabled' : ''} ${this.props.className || ''}`}
                    type={this.props.type}
                    title={this.props.title}
                    to={this.props.to}
                    target={this.props.target}
                    onClick={this.props.onClick}>
                    {this.props.children}
                </Link>
            );
        }

        if ( this.props.href )
        {
            return (
                <a
                    className={`Button color-${this.props.color || ButtonColor.Primary} ${this.props.active ? 'active' : ''} ${this.props.disabled ? 'disabled' : ''} ${this.props.className || ''}`}
                    type={this.props.type}
                    title={this.props.title}
                    href={this.props.href}
                    target={this.props.target}
                    onClick={this.props.onClick}>
                    {this.props.children}
                </a>
            );
        }

        return (
            <button
                className={`Button color-${this.props.color || ButtonColor.Primary} ${this.props.active ? 'active' : ''} ${this.props.disabled ? 'disabled' : ''} ${this.props.className || ''}`}
                type={this.props.type}
                title={this.props.title}
                onClick={this.props.onClick}>
                {this.props.children}
            </button>
        );
    }
}
