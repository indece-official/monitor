import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { ConnectorService, ConnectorV1 } from '../../Services/ConnectorService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPen, faTrash } from '@fortawesome/free-solid-svg-icons';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { LabelValueList } from '../../Components/LabelValueList/LabelValueList';
import { LabelValue } from '../../Components/LabelValueList/LabelValue';
import { Formatter } from '../../utils/Formatter';
import { UserService, UserV1 } from '../../Services/UserService';
import { HostService, HostV1 } from '../../Services/HostService';


export interface ConnectorPageRouteParams
{
    connectorUID:    string;
}


export interface ConnectorPageProps extends RouteComponentProps<ConnectorPageRouteParams>
{
}


interface ConnectorPageState
{
    user:           UserV1 | null;
    connector:      ConnectorV1 | null;
    host:           HostV1 | null;
    initialLoading: boolean;
    error:          Error | null;
}


class $ConnectorPage extends React.Component<ConnectorPageProps, ConnectorPageState>
{
    private readonly _userService:      UserService;
    private readonly _connectorService: ConnectorService;
    private readonly _hostService:      HostService;
    private _intervalReload:            any | null;


    constructor ( props: ConnectorPageProps )
    {
        super(props);

        this.state = {
            user:           null,
            connector:      null,
            host:           null,
            initialLoading: true,
            error:          null
        };

        this._userService = UserService.getInstance();
        this._connectorService = ConnectorService.getInstance();
        this._hostService = HostService.getInstance();
        this._intervalReload = null;
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                error:  null
            });

            const connector = await this._connectorService.getConnector(
                this.props.router.params.connectorUID
            );
            const host = await this._hostService.getHost(connector.host_uid);

            this.setState({
                initialLoading: false,
                connector,
                host
            });
        }
        catch ( err )
        {
            console.error(`Error loading connector: ${(err as Error).message}`, err);

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
            <div className='ConnectorPage'>
                <h1>Connector</h1>

                <ErrorBox error={this.state.error} />

                {this.state.connector ?
                    <div className='ConnectorPage-actions'>
                        <Button to={LinkUtils.make('connector', this.state.connector.uid, 'edit')}>
                            <FontAwesomeIcon icon={faPen} /> Edit
                        </Button>
                    
                        <Button to={LinkUtils.make('connector', this.state.connector.uid, 'delete')}>
                            <FontAwesomeIcon icon={faTrash} /> Delete
                        </Button>
                    </div>
                : null}

                {this.state.connector ?
                    <LabelValueList>
                        <LabelValue
                            label='Host'
                            value={this.state.host ? this.state.host.name : '-'}
                        />
                        
                        <LabelValue
                            label='Type'
                            value={this.state.connector.type || '-'}
                        />
                       
                        <LabelValue
                            label='Version'
                            value={this.state.connector.version || '-'}
                        />
                        
                        <LabelValue
                            label='Status'
                            value={Formatter.connectorStatus(this.state.connector)}
                        />

                        {this.state.connector.error ?
                            <LabelValue
                                label='Fehler'
                                value={this.state.connector.error}
                            />
                        : null}
                    </LabelValueList>
                : null}

                <Spinner active={this.state.initialLoading} />
            </div>
        );
    }
}


export const ConnectorPage = withRouter($ConnectorPage);
