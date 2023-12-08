import React from 'react';
import DayJS from 'dayjs';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faBars, faCircleUser, faEnvelope, faGear, faGrip, faLandmark, faSignOut, faSpaghettiMonsterFlying, faTags, faTimes, faUserCircle } from '@fortawesome/free-solid-svg-icons';
import { UserService, UserV1, isAdmin } from '../../Services/UserService';
import { Link } from 'react-router-dom';
import { LinkUtils } from '../../utils/LinkUtils';
import { Environment } from '../../utils/Environment';

import './SideNav.css';


export interface SideNavProps
{
}


interface SideNavState
{
    expanded:   boolean;
    user:       UserV1 | null;
}


export class SideNav extends React.Component<SideNavProps, SideNavState>
{
    private readonly _userService:  UserService;


    constructor ( props: SideNavProps )
    {
        super(props);

        this.state = {
            expanded:   false,
            user:       null
        };

        this._userService = UserService.getInstance();

        this._logout = this._logout.bind(this);
        this._toggle = this._toggle.bind(this);
        this._close = this._close.bind(this);
    }


    private async _logout ( ): Promise<void>
    {
        try
        {
            await this._userService.logout();
        }
        catch ( err )
        {
            console.error(`Error loggin out user: ${(err as Error).message}`, err);
        }
    }


    private _toggle ( ): void
    {
        const expanded = this.state.expanded;

        this.setState({
            expanded: !expanded
        });

        setImmediate( ( ) =>
        {
            if ( expanded )
            {
                document.removeEventListener('click', this._close);
            }
            else
            {
                document.addEventListener('click', this._close, {
                    once: true
                });
            }
        });
    }
   
   
    private _close ( evt: any ): void
    {
        evt.preventDefault();

        this.setState({
            expanded: false
        });

        document.removeEventListener('click', this._close);
    }


    public componentDidMount ( ): void
    {
        this._userService.isLoggedIn().subscribe(this, ( user ) =>
        {
            this.setState({
                user
            });
        });

        const user = this._userService.isLoggedIn().get();
        this.setState({
            user
        });
    }


    public componentWillUnmount ( ): void
    {
        this._userService.isLoggedIn().unsubscribe(this);
    }


    public render ( )
    {
        if ( Environment.setup.enabled || (!this.state.user && !Environment.setup.enabled) )
        {
            return null;
        }

        return (
            <div className='SideNav'>
                <div className='SideNav-menubutton' onClickCapture={this._toggle}>
                    <FontAwesomeIcon icon={faBars} />
                </div>
                
                <div className={'SideNav-content' + (this.state.expanded ? ' expanded' : '')}>
                    <div className='SideNav-content-close' onClick={this._close}>
                        <FontAwesomeIcon icon={faTimes} />
                    </div>

                    <div className='SideNav-items'>
                        <Link
                            className='SideNav-item'
                            to={LinkUtils.make('')}>
                            <div className='SideNav-item-icon'>
                                <FontAwesomeIcon icon={faGrip} />
                            </div>
                            <div className='SideNav-item-label'>
                                Dashboard
                            </div>
                        </Link>

                        <Link
                            className='SideNav-item'
                            to={LinkUtils.make('hosts')}>
                            <div className='SideNav-item-icon'>
                                <FontAwesomeIcon icon={faLandmark} />
                            </div>
                            <div className='SideNav-item-label'>
                                Hosts
                            </div>
                        </Link>

                        <Link
                            className='SideNav-item'
                            to={LinkUtils.make('agents')}>
                            <div className='SideNav-item-icon'>
                                <FontAwesomeIcon icon={faSpaghettiMonsterFlying} />
                            </div>
                            <div className='SideNav-item-label'>
                                Agents
                            </div>
                        </Link>
                    
                        <Link
                            className='SideNav-item'
                            to={LinkUtils.make('tags')}>
                            <div className='SideNav-item-icon'>
                                <FontAwesomeIcon icon={faTags} />
                            </div>
                            <div className='SideNav-item-label'>
                                Tags
                            </div>
                        </Link>
                    
                        <Link
                            className='SideNav-item'
                            to={LinkUtils.make('notifiers')}>
                            <div className='SideNav-item-icon'>
                                <FontAwesomeIcon icon={faEnvelope} />
                            </div>
                            <div className='SideNav-item-label'>
                                Notifiers
                            </div>
                        </Link>

                        {isAdmin(this.state.user) ?
                            <Link
                                className='SideNav-item'
                                to={LinkUtils.make('users')}>
                                <div className='SideNav-item-icon'>
                                    <FontAwesomeIcon icon={faCircleUser} />
                                </div>
                                <div className='SideNav-item-label'>
                                    Users
                                </div>
                            </Link>
                        : null}
                        
                        {isAdmin(this.state.user) && !Environment.setup.enabled?
                            <Link
                                className='SideNav-item'
                                to={LinkUtils.make('config')}>
                                <div className='SideNav-item-icon'>
                                    <FontAwesomeIcon icon={faGear} />
                                </div>
                                <div className='SideNav-item-label'>
                                    Config
                                </div>
                            </Link>
                        : null}
                    </div>

                    {this.state.user ?
                        <div className='SideNav-user'>
                            <Link
                                className='SideNav-user-image'
                                to={LinkUtils.make('user', this.state.user.uid)}>
                                <FontAwesomeIcon icon={faUserCircle} />
                            </Link>

                            <Link
                                className='SideNav-user-name'
                                to={LinkUtils.make('user', this.state.user.uid)}>
                                {this.state.user.name}
                            </Link>
                            
                            <div className='SideNav-user-actions'>
                                <div className='SideNav-user-action'>
                                    <FontAwesomeIcon icon={faSignOut} onClick={this._logout} />
                                </div>
                            </div>
                        </div>
                    : null}

                    <div className='SideNav-copyright'>
                        Copyright {DayJS().format('YYYY')} indece UG (haftungsbeschränkt)
                    </div>

                    <div className='SideNav-version'>
                        indece Monitor {Environment.server.version}
                    </div>
                </div>
            </div>
        );
    }
}
