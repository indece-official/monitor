import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { HostService, HostV1 } from '../../Services/HostService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPen, faPlus, faTrash } from '@fortawesome/free-solid-svg-icons';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { LabelValueList } from '../../Components/LabelValueList/LabelValueList';
import { LabelValue } from '../../Components/LabelValueList/LabelValue';
import { Tag } from '../../Components/Tag/Tag';
import { CheckService, CheckStatusV1Status, CheckV1 } from '../../Services/CheckService';
import { List } from '../../Components/List/List';
import { ListItem } from '../../Components/List/ListItem';
import { ListItemHeader } from '../../Components/List/ListItemHeader';
import { ListItemHeaderField } from '../../Components/List/ListItemHeaderField';
import { isAdmin, UserService, UserV1 } from '../../Services/UserService';
import { ListItemHeaderAction } from '../../Components/List/ListItemHeaderAction';
import { Formatter } from '../../utils/Formatter';


export interface HostPageRouteParams
{
    hostUID:    string;
}


export interface HostPageProps extends RouteComponentProps<HostPageRouteParams>
{
}


interface HostPageState
{
    host:       HostV1 | null;
    user:       UserV1 | null;
    checks:     Array<CheckV1>;
    loading:    boolean;
    error:      Error | null;
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
            host:       null,
            user:       null,
            checks:     [],
            loading:    false,
            error:      null
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
            <div className='HostPage'>
                <h1>Host</h1>

                <ErrorBox error={this.state.error} />

                {this.state.host && isAdmin(this.state.user) ?
                    <div className='HostPage-actions'>
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

                {this.state.host ?
                    <LabelValueList>
                        <LabelValue
                            label='Name'
                            value={this.state.host.name || '-'}
                        />

                        <LabelValue label='Tags'>
                            {this.state.host.tags.map( ( tag ) => (
                                <Tag key={tag.name} color={tag.color}>{tag.name}</Tag>
                            ))}
                        </LabelValue>
                    </LabelValueList>
                : null}

                <br />

                <List>
                    {this.state.checks.map( ( check ) => (
                        <ListItem key={check.uid}>
                            <ListItemHeader>
                                <ListItemHeaderField
                                    text={Formatter.checkStatus(check.status?.status || CheckStatusV1Status.Unknown)}
                                />
                                
                                <ListItemHeaderField
                                    text={check.name}
                                    subtext={check.status?.message}
                                    grow={true}
                                />

                                {isAdmin(this.state.user) && check.custom && this.state.host ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('host', this.state.host.uid,'check', check.uid, 'edit')}
                                        icon={faPen}
                                    />
                                : null}
                                
                                {isAdmin(this.state.user) && check.custom && this.state.host ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('host', this.state.host.uid,'host', check.uid, 'delete')}
                                        icon={faTrash}
                                    />
                                : null}
                            </ListItemHeader>
                        </ListItem>
                    ))}
                </List>

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}


export const HostPage = withRouter($HostPage);