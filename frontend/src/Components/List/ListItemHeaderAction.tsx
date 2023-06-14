import React from 'react';
import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Link, To } from 'react-router-dom';

import './List.css';


export interface ListItemHeaderActionProps
{
    to?:        To;
    icon:       IconProp;
    title?:     string;
    active?:    boolean;
    onClick?:   ( evt: React.MouseEvent ) => any;
}


export class ListItemHeaderAction extends React.Component<ListItemHeaderActionProps>
{
    public render ( )
    {
        if ( this.props.to )
        {
            return (
                <Link
                    className={`ListItemHeaderAction ${this.props.active ? 'active' : ''} clickable`}
                    onClick={this.props.onClick}
                    title={this.props.title}
                    to={this.props.to}>
                    <div className='ListItemHeaderAction-inner'>
                        <FontAwesomeIcon icon={this.props.icon} fixedWidth={true} />
                    </div>
                </Link>
            );
        }

        return (
            <div
                className={`ListItemHeaderAction ${this.props.active ? 'active' : ''} ${this.props.onClick ? 'clickable' : ''}`}
                onClick={this.props.onClick}
                title={this.props.title}>
                <div className='ListItemHeaderAction-inner'>
                    <FontAwesomeIcon icon={this.props.icon} fixedWidth={true} />
                </div>
            </div>
        );
    }
}
