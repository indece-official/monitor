import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { HostService, HostV1 } from '../../Services/HostService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPen, faPlus, faTrash } from '@fortawesome/free-solid-svg-icons';
import { List } from '../../Components/List/List';
import { ListEmpty } from '../../Components/List/ListEmpty';
import { ListItem } from '../../Components/List/ListItem';
import { ListItemHeaderField } from '../../Components/List/ListItemHeaderField';
import { ListItemHeader } from '../../Components/List/ListItemHeader';
import { ListItemHeaderAction } from '../../Components/List/ListItemHeaderAction';
import { isAdmin, UserService, UserV1 } from '../../Services/UserService';


export interface HostsPageProps
{
}


interface HostsPageState
{
    hosts:      Array<HostV1>;
    user:       UserV1 | null;
    loading:    boolean;
    error:      Error | null;
}


export class HostsPage extends React.Component<HostsPageProps, HostsPageState>
{
    private readonly _hostService: HostService;
    private readonly _userService: UserService;


    constructor ( props: HostsPageProps )
    {
        super(props);

        this.state = {
            hosts:      [],
            user:       null,
            loading:    false,
            error:      null
        };

        this._hostService  = HostService.getInstance();
        this._userService  = UserService.getInstance();
    }


    private async _load ( ): Promise<void>
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
        
        await this._load();
    }


    public componentWillUnmount ( ): void
    {
        this._userService.isLoggedIn().unsubscribe(this);
    }


    public render ( )
    {
        return (
            <div className='HostsPage'>
                <h1>Hosts</h1>

                <ErrorBox error={this.state.error} />

                {isAdmin(this.state.user) ?
                    <Button to={LinkUtils.make('host', 'add')}>
                        <FontAwesomeIcon icon={faPlus} /> Add a host
                    </Button>
                : null}

                <List>
                    {this.state.hosts.length === 0 && !this.state.loading && !this.state.error ?
                        <ListEmpty>
                            No hosts found
                        </ListEmpty>
                    : null}

                    {this.state.hosts.map( ( host ) => (
                        <ListItem key={host.uid}>
                            <ListItemHeader>
                                <ListItemHeaderField
                                    to={LinkUtils.make('host', host.uid)}
                                    grow={true}
                                    text={host.name}
                                />

                                {isAdmin(this.state.user) ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('host', host.uid, 'edit')}
                                        icon={faPen}
                                    />
                                : null}
                                
                                {isAdmin(this.state.user) ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('host', host.uid, 'delete')}
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
