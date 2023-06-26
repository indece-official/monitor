import React from 'react';
import DayJS from 'dayjs';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faBars, faSignOut, faTimes, faUserCircle } from '@fortawesome/free-solid-svg-icons';
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
                        {Environment.setup.enabled ?
                            <div className='SideNav-item'>
                                <Link to={LinkUtils.make('setup')}>Setup</Link>
                            </div>
                        : null}
                        
                        {!this.state.user && !Environment.setup.enabled ?
                            <div className='SideNav-item'>
                                <Link to={LinkUtils.make('login')}>Login</Link>
                            </div>
                        : null}
                        
                        {this.state.user && !Environment.setup.enabled?
                            <div className='SideNav-item'>
                                <Link to={LinkUtils.make('dashboard')}>Dashboard</Link>
                            </div>
                        : null}
                        
                        {this.state.user && !Environment.setup.enabled?
                            <div className='SideNav-item'>
                                <Link to={LinkUtils.make('hosts')}>Hosts</Link>
                            </div>
                        : null}

                        {this.state.user && !Environment.setup.enabled?
                            <div className='SideNav-item'>
                                <Link to={LinkUtils.make('agents')}>Agents</Link>
                            </div>
                        : null}
                        
                        {this.state.user && !Environment.setup.enabled?
                            <div className='SideNav-item'>
                                <Link to={LinkUtils.make('tags')}>Tags</Link>
                            </div>
                        : null}
                        
                        {this.state.user && !Environment.setup.enabled?
                            <div className='SideNav-item'>
                                <Link to={LinkUtils.make('notifiers')}>Notifiers</Link>
                            </div>
                        : null}

                        {this.state.user && !Environment.setup.enabled?
                            <div className='SideNav-item'>
                                <Link to={LinkUtils.make('maintenances')}>Maintenances</Link>
                            </div>
                        : null}
                        
                        {isAdmin(this.state.user) && !Environment.setup.enabled?
                            <div className='SideNav-item'>
                                <Link to={LinkUtils.make('users')}>Users</Link>
                            </div>
                        : null}
                        
                        {isAdmin(this.state.user) && !Environment.setup.enabled?
                            <div className='SideNav-item'>
                                <Link to={LinkUtils.make('config')}>Config</Link>
                            </div>
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
                                <FontAwesomeIcon icon={faSignOut} onClick={this._logout} />
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
