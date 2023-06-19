import React from 'react';
import { AgentV1, AgentV1Status, AgentService } from '../../../Services/AgentService';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { Formatter } from '../../../utils/Formatter';


export interface AddAgentWaitRegisteredStepProps
{
    agentUID:   string;
    onFinish:   ( ) => any;
    onFailed:   ( ) => any;
}


interface AddAgentWaitRegisteredStepState
{
    agent:      AgentV1 | null;
    loading:    boolean;
    error:      Error | null;
}


export class AddAgentWaitRegisteredStep extends React.Component<AddAgentWaitRegisteredStepProps, AddAgentWaitRegisteredStepState>
{
    private readonly _agentService:  AgentService;
    private _intervalReload:    any | null;


    constructor ( props: AddAgentWaitRegisteredStepProps )
    {
        super(props);

        this.state = {
            agent:      null,
            loading:    true,
            error:      null
        };

        this._agentService = AgentService.getInstance();

        this._intervalReload = null;
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const agent = await this._agentService.getAgent(this.props.agentUID);

            this.setState({
                loading:    false,
                agent
            });

            if ( agent.status === AgentV1Status.Unregistered )
            {
                if ( this._intervalReload )
                {
                    clearInterval(this._intervalReload);
                    this._intervalReload = null;
                }

                this.props.onFinish();
            }
            else if ( agent.status === AgentV1Status.Error )
            {
                if ( this._intervalReload )
                {
                    clearInterval(this._intervalReload);
                    this._intervalReload = null;
                }

                this.props.onFailed();
            }
        }
        catch ( err )
        {
            console.error(`Error loading agent execution ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        this._intervalReload = setInterval( async ( ) => {
            await this._load();
        }, 5000);

        await this._load();
    }


    public componentWillUnmount ( ): void
    {
        if ( this._intervalReload )
        {
            clearInterval(this._intervalReload);
            this._intervalReload = null;
        }
    }

    public render ( )
    {
        return (
            <div className='AddAgentWaitRegisteredStep'>
                <h1>Add a agent</h1>

                <ErrorBox error={this.state.error} />

                {this.state.agent ?
                    <div>
                        <div>Waiting for agent to be configured ...</div>
                        <div>Status: {Formatter.agentStatus(this.state.agent)}</div>
                    </div>
                : null}

                <Spinner active={true} />
            </div>
        );
    }
}
