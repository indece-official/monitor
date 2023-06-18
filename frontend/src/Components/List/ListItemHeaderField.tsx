import React from 'react';
import { Link, To } from 'react-router-dom';

import './List.css';


export interface ListItemHeaderFieldProps
{
    className?: string;
    grow?:      boolean;
    text:       string;
    subtext?:   string;
    to?:        To;
    onClick?:   ( evt: React.MouseEvent ) => any;
}


export class ListItemHeaderField extends React.Component<ListItemHeaderFieldProps>
{
    public render ( )
    {
        if ( this.props.to )
        {
            return (
                <Link
                    className={`ListItemHeaderField ${this.props.grow ? 'grow': ''} clickable ${this.props.className || ''}`}
                    to={this.props.to}
                    onClick={this.props.onClick}>
                    <div className='ListItemHeaderField-text'>
                        {this.props.text}
                    </div>
    
                    {this.props.subtext ?
                        <div className='ListItemHeaderField-subtext'>
                            {this.props.subtext}
                        </div>
                    : null}
                </Link>
            );
        }

        return (
            <div
                className={`ListItemHeaderField ${this.props.grow ? 'grow': ''} ${this.props.onClick ? 'clickable': ''} ${this.props.className || ''}`}
                onClick={this.props.onClick}>
                <div className='ListItemHeaderField-text'>
                    {this.props.text}
                </div>

                {this.props.subtext ?
                    <div className='ListItemHeaderField-subtext'>
                        {this.props.subtext}
                    </div>
                : null}
            </div>
        );
    }
}
