import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { AgentService, AgentV1 } from '../../Services/AgentService';
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


export interface AgentPageRouteParams
{
    agentUID:    string;
}


export interface AgentPageProps extends RouteComponentProps<AgentPageRouteParams>
{
}


interface AgentPageState
{
    user:           UserV1 | null;
    agent:          AgentV1 | null;
    host:           HostV1 | null;
    initialLoading: boolean;
    error:          Error | null;
}


class $AgentPage extends React.Component<AgentPageProps, AgentPageState>
{
    private readonly _userService:  UserService;
    private readonly _agentService: AgentService;
    private readonly _hostService:  HostService;
    private _intervalReload:        any | null;


    constructor ( props: AgentPageProps )
    {
        super(props);

        this.state = {
            user:           null,
            agent:          null,
            host:           null,
            initialLoading: true,
            error:          null
        };

        this._userService = UserService.getInstance();
        this._agentService = AgentService.getInstance();
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

            const agent = await this._agentService.getAgent(
                this.props.router.params.agentUID
            );
            const host = await this._hostService.getHost(agent.host_uid);

            this.setState({
                initialLoading: false,
                agent,
                host
            });
        }
        catch ( err )
        {
            console.error(`Error loading agent: ${(err as Error).message}`, err);

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
            <div className='AgentPage'>
                <h1>Agent</h1>

                <ErrorBox error={this.state.error} />

                {this.state.agent ?
                    <div className='AgentPage-actions'>
                        <Button to={LinkUtils.make('agent', this.state.agent.uid, 'edit')}>
                            <FontAwesomeIcon icon={faPen} /> Edit
                        </Button>
                    
                        <Button to={LinkUtils.make('agent', this.state.agent.uid, 'delete')}>
                            <FontAwesomeIcon icon={faTrash} /> Delete
                        </Button>
                    </div>
                : null}

                {this.state.agent ?
                    <LabelValueList>
                        <LabelValue
                            label='Host'
                            value={this.state.host ? this.state.host.name : '-'}
                        />
                        
                        <LabelValue
                            label='Type'
                            value={this.state.agent.type || '-'}
                        />
                       
                        <LabelValue
                            label='Version'
                            value={this.state.agent.version || '-'}
                        />
                        
                        <LabelValue
                            label='Status'
                            value={Formatter.agentStatus(this.state.agent)}
                        />

                        {this.state.agent.error ?
                            <LabelValue
                                label='Fehler'
                                value={this.state.agent.error}
                            />
                        : null}
                    </LabelValueList>
                : null}

                <Spinner active={this.state.initialLoading} />
            </div>
        );
    }
}


export const AgentPage = withRouter($AgentPage);
