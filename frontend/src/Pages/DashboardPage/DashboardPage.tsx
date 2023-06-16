import React from 'react';
import { UserService, UserV1 } from '../../Services/UserService';
import { CheckStatusV1Status } from '../../Services/CheckService';
import { HostService, HostV1 } from '../../Services/HostService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { Tag } from '../../Components/Tag/Tag';
import { Formatter } from '../../utils/Formatter';
import { Link } from 'react-router-dom';
import { LinkUtils } from '../../utils/LinkUtils';

import './DashboardPage.css';


export interface DashboardPageProps
{
}


interface HostStats
{
    status:     CheckStatusV1Status;
}


interface DashboardPageState
{
    user:           UserV1 | null;
    hosts:          Array<HostV1>;
    hostStatuses:   Record<string, HostStats>;
    loading:        boolean;
    error:          Error | null;
}


export class DashboardPage extends React.Component<DashboardPageProps, DashboardPageState>
{
    private readonly _userService:  UserService;
    private readonly _hostService:  HostService;
    private _intervalReloadChecks:  any | null;


    constructor ( props: DashboardPageProps )
    {
        super(props);

        this.state = {
            user:           null,
            hosts:          [],
            hostStatuses:   {},
            loading:        true,
            error:          null
        };

        this._userService = UserService.getInstance();
        this._hostService = HostService.getInstance();

        this._intervalReloadChecks = null;
    }


    private async _loadHosts ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const hosts = await this._hostService.getHosts();
            const hostStatuses: Record<string, HostStats> = {};
            for ( const host of hosts )
            {
                hostStatuses[host.uid] = {
                    status: CheckStatusV1Status.Unknown
                };

                if ( host.status.count_unknown > 0 )
                {
                    hostStatuses[host.uid].status = CheckStatusV1Status.Unknown;
                }
                else if ( host.status.count_critical > 0 )
                {
                    hostStatuses[host.uid].status = CheckStatusV1Status.Critical;
                }
                else if ( host.status.count_warning > 0 )
                {
                    hostStatuses[host.uid].status = CheckStatusV1Status.Warning;
                }
                else if ( host.status.count_ok > 0 )
                {
                    hostStatuses[host.uid].status = CheckStatusV1Status.Ok;
                }
                else
                {
                    hostStatuses[host.uid].status = CheckStatusV1Status.Unknown;
                }
            }

            this.setState({
                loading:    false,
                hosts,
                hostStatuses
            });
        }
        catch ( err )
        {
            console.error(`Error loading hosts: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }
   

    public async componentDidMount ( ): Promise<void>
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

        await this._loadHosts();

        this._intervalReloadChecks = setInterval( async ( ) =>
        {
            await this._loadHosts();
        }, 10000);
    }


    public componentWillUnmount ( ): void
    {
        this._userService.isLoggedIn().unsubscribe(this);

        if ( this._intervalReloadChecks )
        {
            clearInterval(this._intervalReloadChecks);
            this._intervalReloadChecks = null;
        }
    }


    public render ( )
    {
        return (
            <div className='DashboardPage'>
                <h1>Welcome</h1>

                <ErrorBox error={this.state.error} />

                <div className='DashboardPage-hosts'>
                    {this.state.hosts.map( ( host ) => (
                        <Link
                            key={host.uid}
                            className={'DashboardPage-host status-' + this.state.hostStatuses[host.uid]?.status.toLowerCase() || 'unknown'}
                            to={LinkUtils.make('host', host.uid)}>
                            <div className='DashboardPage-host-status'>
                                <div>{Formatter.checkStatus(this.state.hostStatuses[host.uid]?.status || CheckStatusV1Status.Unknown)}</div>
                            </div>

                            <div className='DashboardPage-host-details'>
                                <div className='DashboardPage-host-name' title={'Host ' + host.name}>
                                    {host.name}
                                </div>
                                
                                <div className='DashboardPage-host-tags'>
                                    {host.tags.map( ( tag ) => (
                                        <Tag key={tag.name} color={tag.color}>{tag.name}</Tag>
                                    ))}
                                </div>

                                <div className='DashboardPage-host-message'>
                                    {Formatter.checkStatus(CheckStatusV1Status.Critical)}: {host.status.count_critical}, 
                                    {Formatter.checkStatus(CheckStatusV1Status.Warning)}: {host.status.count_warning}, 
                                    {Formatter.checkStatus(CheckStatusV1Status.Ok)}: {host.status.count_ok}, 
                                    {Formatter.checkStatus(CheckStatusV1Status.Unknown)}: {host.status.count_unknown}
                                </div>
                            </div>
                        </Link>
                    ))}
                </div>

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}
