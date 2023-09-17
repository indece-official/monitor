import React from 'react';

import './PageContent.css';


export enum PageContentSize
{
    Small   = 'small',
    Large   = 'large',
    Full    = 'full'
}


export interface PageContentProps
{
    className?: string;
    size?:      PageContentSize;
    centered?:  boolean;
    children?:  React.ReactNode | null;
}


export class PageContent extends React.Component<PageContentProps>
{
    public render ( )
    {
        return (
            <div className={`PageContent ${this.props.size || PageContentSize.Large} ${this.props.className || ''} ${this.props.centered ? 'centered' : ''}`}>
                <div className='PageContent-inner'>
                    {this.props.children}
                </div>
            </div>
        );
    }
}
