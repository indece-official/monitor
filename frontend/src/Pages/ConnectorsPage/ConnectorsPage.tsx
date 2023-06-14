import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { ConnectorService, ConnectorV1 } from '../../Services/ConnectorService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPen, faPlus, faTrash } from '@fortawesome/free-solid-svg-icons';
import { Formatter } from '../../utils/Formatter';
import { UserService, UserV1, isAdmin } from '../../Services/UserService';
import { List } from '../../Components/List/List';
import { ListEmpty } from '../../Components/List/ListEmpty';
import { ListItem } from '../../Components/List/ListItem';
import { ListItemHeader } from '../../Components/List/ListItemHeader';
import { ListItemHeaderField } from '../../Components/List/ListItemHeaderField';
import { ListItemHeaderAction } from '../../Components/List/ListItemHeaderAction';
import { ListItemBody } from '../../Components/List/ListItemBody';
import { HostService, HostV1 } from '../../Services/HostService';


export interface ConnectorsPageProps
{
}


interface ConnectorsPageState
{
    user:           UserV1 | null;
    connectors:     Array<ConnectorV1>;
    hostsMap:       Record<string, HostV1>;
    initialLoading: boolean;
    error:          Error | null;
}


export class ConnectorsPage extends React.Component<ConnectorsPageProps, ConnectorsPageState>
{
    private readonly _userService:      UserService;
    private readonly _connectorService: ConnectorService;
    private readonly _hostService:      HostService;
    private _intervalReload:            any | null;


    constructor ( props: ConnectorsPageProps )
    {
        super(props);

        this.state = {
            user:           null,
            connectors:     [],
            hostsMap:       {},
            initialLoading: true,
            error:          null
        };

        this._userService = UserService.getInstance();
        this._connectorService = ConnectorService.getInstance();
        this._hostService  = HostService.getInstance();
        this._intervalReload = null;
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                error:      null
            });

            const connectors = await this._connectorService.getConnectors();
            const hosts = await this._hostService.getHosts();
            const hostsMap: Record<string, HostV1> = {};
            for ( const host of hosts )
            {
                hostsMap[host.uid] = host;
            }

            this.setState({
                initialLoading: false,
                connectors,
                hostsMap
            });
        }
        catch ( err )
        {
            console.error(`Error loading connectors: ${(err as Error).message}`, err);

            this.setState({
                initialLoading: false,
                error:          err as Error
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

        this._intervalReload = setInterval( async ( ) =>
        {
            await this._load();
        }, 10000);
    }


    public componentWillUnmount ( ): void
    {
        if ( this._intervalReload )
        {
            clearInterval(this._intervalReload);
            this._intervalReload = null;
        }

        this._userService.isLoggedIn().unsubscribe(this);
    }


    public render ( )
    {
        return (
            <div className='ConnectorsPage'>
                <h1>Connectors</h1>

                <ErrorBox error={this.state.error} />

                {isAdmin(this.state.user) ?
                    <Button to={LinkUtils.make('connector', 'add')}>
                        <FontAwesomeIcon icon={faPlus} /> Add a connector
                    </Button>
                : null}

                <List>
                    {this.state.connectors.length === 0 && !this.state.initialLoading && !this.state.error ?
                        <ListEmpty>
                            No connectors found
                        </ListEmpty>
                    : null}

                    {this.state.connectors.map( ( connector ) => (
                        <ListItem key={connector.uid}>
                            <ListItemHeader>
                                <ListItemHeaderField
                                    to={LinkUtils.make('connector', connector.uid)}
                                    grow={true}
                                    text={`${connector.type} ${connector.version}`}
                                    subtext={this.state.hostsMap[connector.host_uid]?.name || '-'}
                                />

                                <ListItemHeaderField
                                    to={LinkUtils.make('connector', connector.uid)}
                                    text={Formatter.connectorStatus(connector)}
                                />

                                {isAdmin(this.state.user) ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('connector', connector.uid, 'edit')}
                                        icon={faPen}
                                    />
                                : null}
                                
                                {isAdmin(this.state.user) ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('connector', connector.uid, 'delete')}
                                        icon={faTrash}
                                    />
                                : null}
                            </ListItemHeader>
                           
                            {connector.error ?
                                <ListItemBody>
                                    <ErrorBox error={new Error(connector.error)} />
                                </ListItemBody>
                            : null}
                        </ListItem>
                    ))}
                </List>

                <Spinner active={this.state.initialLoading} />
            </div>
        );
    }
}
