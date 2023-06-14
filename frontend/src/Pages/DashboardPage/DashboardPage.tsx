import React from 'react';
import { UserService, UserV1 } from '../../Services/UserService';
import { CheckService, CheckStatusV1Status, CheckV1 } from '../../Services/CheckService';
import { HostService, HostV1 } from '../../Services/HostService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { Tag } from '../../Components/Tag/Tag';

import './DashboardPage.css';
import { Formatter } from '../../utils/Formatter';
import { Link } from 'react-router-dom';
import { LinkUtils } from '../../utils/LinkUtils';


export interface DashboardPageProps
{
}


interface HostStats
{
    status:         CheckStatusV1Status;
    count_critical: number;
    count_warning:  number;
    count_ok:       number;
    count_unknown:  number;
}


interface DashboardPageState
{
    user:           UserV1 | null;
    hosts:          Array<HostV1>;
    checks:         Array<CheckV1>;
    hostStatuses:   Record<string, HostStats>;
    loading:        boolean;
    error:          Error | null;
}


export class DashboardPage extends React.Component<DashboardPageProps, DashboardPageState>
{
    private readonly _userService:  UserService;
    private readonly _hostService:  HostService;
    private readonly _checkService: CheckService;
    private _intervalReloadChecks:  any | null;


    constructor ( props: DashboardPageProps )
    {
        super(props);

        this.state = {
            user:           null,
            hosts:          [],
            checks:         [],
            hostStatuses:   {},
            loading:        true,
            error:          null
        };

        this._userService = UserService.getInstance();
        this._hostService = HostService.getInstance();
        this._checkService = CheckService.getInstance();

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

            this.setState({
                loading:    false,
                hosts
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
   
   
    private async _loadChecks ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const checks = await this._checkService.getChecks();
            const hostStatuses: Record<string, HostStats> = {};
            for ( const check of checks )
            {
                if ( ! hostStatuses[check.host_uid] )
                {
                    hostStatuses[check.host_uid] = {
                        status:         CheckStatusV1Status.Unknown,
                        count_critical: 0,
                        count_warning:  0,
                        count_ok:       0,
                        count_unknown:  0
                    };
                }

                switch ( check.status?.status )
                {
                    case CheckStatusV1Status.Critical:
                        hostStatuses[check.host_uid].count_critical++;
                        break;
                    case CheckStatusV1Status.Warning:
                        hostStatuses[check.host_uid].count_warning++;
                        break;
                    case CheckStatusV1Status.Ok:
                        hostStatuses[check.host_uid].count_ok++;
                        break;
                    case CheckStatusV1Status.Unknown:
                    default:
                        hostStatuses[check.host_uid].count_unknown++;
                        break;
                }
            }

            for ( const hostUID in hostStatuses )
            {
                if ( hostStatuses[hostUID].count_unknown > 0 )
                {
                    hostStatuses[hostUID].status = CheckStatusV1Status.Unknown;
                }
                else if ( hostStatuses[hostUID].count_critical > 0 )
                {
                    hostStatuses[hostUID].status = CheckStatusV1Status.Critical;
                }
                else if ( hostStatuses[hostUID].count_warning > 0 )
                {
                    hostStatuses[hostUID].status = CheckStatusV1Status.Warning;
                }
                else if ( hostStatuses[hostUID].count_ok > 0 )
                {
                    hostStatuses[hostUID].status = CheckStatusV1Status.Ok;
                }
                else
                {
                    hostStatuses[hostUID].status = CheckStatusV1Status.Unknown;
                }
            }

            this.setState({
                loading:    false,
                checks,
                hostStatuses
            });
        }
        catch ( err )
        {
            console.error(`Error loading checks: ${(err as Error).message}`, err);

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
        await this._loadChecks();

        this._intervalReloadChecks = setInterval( async ( ) =>
        {
            await this._loadChecks();
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
                                    {Formatter.checkStatus(CheckStatusV1Status.Critical)}: {this.state.hostStatuses[host.uid]?.count_critical ?? '-'}, 
                                    {Formatter.checkStatus(CheckStatusV1Status.Warning)}: {this.state.hostStatuses[host.uid]?.count_warning ?? '-'}, 
                                    {Formatter.checkStatus(CheckStatusV1Status.Ok)}: {this.state.hostStatuses[host.uid]?.count_ok ?? '-'}, 
                                    {Formatter.checkStatus(CheckStatusV1Status.Unknown)}: {this.state.hostStatuses[host.uid]?.count_unknown ?? '-'}
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
