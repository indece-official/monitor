import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { AgentService, AgentV1 } from '../../Services/AgentService';
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


export interface AgentsPageProps
{
}


interface AgentsPageState
{
    user:           UserV1 | null;
    agents:         Array<AgentV1>;
    hostsMap:       Record<string, HostV1>;
    initialLoading: boolean;
    error:          Error | null;
}


export class AgentsPage extends React.Component<AgentsPageProps, AgentsPageState>
{
    private readonly _userService:  UserService;
    private readonly _agentService: AgentService;
    private readonly _hostService:  HostService;
    private _intervalReload:        any | null;


    constructor ( props: AgentsPageProps )
    {
        super(props);

        this.state = {
            user:           null,
            agents:         [],
            hostsMap:       {},
            initialLoading: true,
            error:          null
        };

        this._userService = UserService.getInstance();
        this._agentService = AgentService.getInstance();
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

            const agents = await this._agentService.getAgents();
            const hosts = await this._hostService.getHosts();
            const hostsMap: Record<string, HostV1> = {};
            for ( const host of hosts )
            {
                hostsMap[host.uid] = host;
            }

            this.setState({
                initialLoading: false,
                agents,
                hostsMap
            });
        }
        catch ( err )
        {
            console.error(`Error loading agents: ${(err as Error).message}`, err);

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
            <div className='AgentsPage'>
                <h1>Agents</h1>

                <ErrorBox error={this.state.error} />

                {isAdmin(this.state.user) ?
                    <Button to={LinkUtils.make('agent', 'add')}>
                        <FontAwesomeIcon icon={faPlus} /> Add a agent
                    </Button>
                : null}

                <List>
                    {this.state.agents.length === 0 && !this.state.initialLoading && !this.state.error ?
                        <ListEmpty>
                            No agents found
                        </ListEmpty>
                    : null}

                    {this.state.agents.map( ( agent ) => (
                        <ListItem key={agent.uid}>
                            <ListItemHeader>
                                <ListItemHeaderField
                                    to={LinkUtils.make('agent', agent.uid)}
                                    grow={true}
                                    text={`${agent.type} ${agent.version}`}
                                    subtext={this.state.hostsMap[agent.host_uid]?.name || '-'}
                                />

                                <ListItemHeaderField
                                    to={LinkUtils.make('agent', agent.uid)}
                                    text={Formatter.agentStatus(agent)}
                                />

                                {isAdmin(this.state.user) ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('agent', agent.uid, 'edit')}
                                        icon={faPen}
                                    />
                                : null}
                                
                                {isAdmin(this.state.user) ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('agent', agent.uid, 'delete')}
                                        icon={faTrash}
                                    />
                                : null}
                            </ListItemHeader>
                           
                            {agent.error ?
                                <ListItemBody>
                                    <ErrorBox error={new Error(agent.error)} />
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
