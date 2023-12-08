import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { HostService, HostV1 } from '../../Services/HostService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPen, faPlus, faRefresh, faTrash } from '@fortawesome/free-solid-svg-icons';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { Tag } from '../../Components/Tag/Tag';
import { CheckService, CheckStatusV1Status, CheckV1 } from '../../Services/CheckService';
import { List } from '../../Components/List/List';
import { ListItem } from '../../Components/List/ListItem';
import { ListItemHeader } from '../../Components/List/ListItemHeader';
import { ListItemHeaderField } from '../../Components/List/ListItemHeaderField';
import { isAdmin, UserService, UserV1 } from '../../Services/UserService';
import { ListItemHeaderAction } from '../../Components/List/ListItemHeaderAction';
import { Formatter } from '../../utils/Formatter';
import { ListItemBody } from '../../Components/List/ListItemBody';
import { PageContent } from '../../Components/PageContent/PageContent';

import './HostPage.css';


export interface HostPageRouteParams
{
    hostUID:    string;
}


export interface HostPageProps extends RouteComponentProps<HostPageRouteParams>
{
}


interface HostPageState
{
    host:               HostV1 | null;
    user:               UserV1 | null;
    checks:             Array<CheckV1>;
    expandedCheckUID:   string | null;
    loading:            boolean;
    error:              Error | null;
}


class $HostPage extends React.Component<HostPageProps, HostPageState>
{
    private readonly _hostService: HostService;
    private readonly _userService: UserService;
    private readonly _checkService: CheckService;
    private _intervalReloadChecks:  any | null;


    constructor ( props: HostPageProps )
    {
        super(props);

        this.state = {
            host:               null,
            user:               null,
            checks:             [],
            expandedCheckUID:   null,
            loading:            false,
            error:              null
        };

        this._hostService  = HostService.getInstance();
        this._userService  = UserService.getInstance();
        this._checkService = CheckService.getInstance();

        this._intervalReloadChecks = null;
    }


    private async _loadHost ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const host = await this._hostService.getHost(
                this.props.router.params.hostUID
            );

            this.setState({
                loading:    false,
                host
            });
        }
        catch ( err )
        {
            console.error(`Error loading host: ${(err as Error).message}`, err);

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

            const checks = await this._checkService.getHostChecks(this.props.router.params.hostUID);

            this.setState({
                loading:    false,
                checks
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


    private async _executeCheck ( checkUID: string ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            await this._checkService.executeCheck(checkUID);
        }
        catch ( err )
        {
            console.error(`Error executing checks: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    private _toggleCheck ( checkUID: string ): void
    {
        if ( checkUID === this.state.expandedCheckUID )
        {
            this.setState({
                expandedCheckUID:   null
            });
        }
        else
        {
            this.setState({
                expandedCheckUID: checkUID
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

        await this._loadHost();
        await this._loadChecks();

        this._intervalReloadChecks = setInterval( async ( ) =>
        {
            await this._loadChecks();
        }, 10000);
    }


    public componentWillUnmount ( ): void
    {
        if ( this._intervalReloadChecks )
        {
            clearInterval(this._intervalReloadChecks);
            this._intervalReloadChecks = null;
        }

        this._userService.isLoggedIn().unsubscribe(this);
    }


    public render ( )
    {
        return (
            <PageContent className='HostPage'>
                <h1>{this.state.host?.name || '-'}</h1>
                
                {this.state.host && this.state.host.tags.length > 0 ?
                    <div className='HostPage-tags'>
                        {this.state.host.tags.map( ( tag ) => (
                            <Tag key={tag.name} color={tag.color}>{tag.name}</Tag>
                        ))}
                    </div>
                : null}

                <ErrorBox error={this.state.error} />

                {this.state.host && isAdmin(this.state.user) ?
                    <div className='HostPage-actions'>
                        <Button to={LinkUtils.make('host', this.state.host.uid, 'agent', 'add')}>
                            <FontAwesomeIcon icon={faPlus} /> Add an agent
                        </Button>

                        <Button to={LinkUtils.make('host', this.state.host.uid, 'check', 'add')}>
                            <FontAwesomeIcon icon={faPlus} /> Add a check
                        </Button>
                        
                        <Button to={LinkUtils.make('host', this.state.host.uid, 'edit')}>
                            <FontAwesomeIcon icon={faPen} /> Edit
                        </Button>
                    
                        <Button to={LinkUtils.make('host', this.state.host.uid, 'delete')}>
                            <FontAwesomeIcon icon={faTrash} /> Delete
                        </Button>
                    </div>
                : null}

                <br />

                <List>
                    {this.state.checks.map( ( check ) => (
                        <ListItem key={check.uid}>
                            <ListItemHeader>
                                <ListItemHeaderField
                                    onClick={ ( ) => this._toggleCheck(check.uid) }
                                    className={`HostPage-check-status status-${(check.status?.status || CheckStatusV1Status.Unknown).toLowerCase()}`}
                                    text={Formatter.checkStatus(check.status?.status || CheckStatusV1Status.Unknown)}
                                />
                                
                                <ListItemHeaderField
                                    onClick={ ( ) => this._toggleCheck(check.uid) }
                                    text={check.name}
                                    subtext={check.status?.message}
                                    grow={true}
                                />

                                {isAdmin(this.state.user) ?
                                    <ListItemHeaderAction
                                        onClick={ ( ) => this._executeCheck(check.uid)}
                                        icon={faRefresh}
                                    />
                                : null}
                                
                                {isAdmin(this.state.user) && check.custom && this.state.host ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('host', this.state.host.uid,'check', check.uid, 'edit')}
                                        icon={faPen}
                                    />
                                : null}
                                
                                {isAdmin(this.state.user) && this.state.host ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('host', this.state.host.uid,'check', check.uid, 'delete')}
                                        icon={faTrash}
                                    />
                                : null}
                            </ListItemHeader>

                            {this.state.expandedCheckUID === check.uid ?
                                <ListItemBody>
                                    {check.status ?
                                        <div>
                                            <div>Last checked: {check.status.datetime_created}</div>
                                            {Object.entries(check.status.data).map( ( entry, i ) => (
                                                <div key={i}>{entry[0]}: {entry[1]}</div>
                                            ))}
                                        </div>
                                    :
                                        <div>
                                            No check results captured yet
                                        </div>
                                    }
                                </ListItemBody>
                            : null}
                        </ListItem>
                    ))}
                </List>

                <Spinner active={this.state.loading} />
            </PageContent>
        );
    }
}


export const HostPage = withRouter($HostPage);
