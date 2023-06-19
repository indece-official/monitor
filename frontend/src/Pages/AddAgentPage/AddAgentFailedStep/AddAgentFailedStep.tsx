import React from 'react';
import { AgentV1, AgentService } from '../../../Services/AgentService';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';


export interface AddAgentFailedStepProps
{
    agentUID:   string;
}


interface AddAgentFailedStepState
{
    agent:      AgentV1 | null;
    loading:    boolean;
    error:      Error | null;
}


export class AddAgentFailedStep extends React.Component<AddAgentFailedStepProps, AddAgentFailedStepState>
{
    private readonly _agentService:  AgentService;


    constructor ( props: AddAgentFailedStepProps )
    {
        super(props);

        this.state = {
            agent:      null,
            loading:    true,
            error:      null
        };

        this._agentService = AgentService.getInstance();
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
        await this._load();
    }


    public render ( )
    {
        return (
            <div className='AddAgentFailedStep'>
                <h1>Add a agent</h1>

                <ErrorBox error={this.state.error} />

                <Spinner active={this.state.loading} />

                {this.state.agent && this.state.agent.error ?
                    <ErrorBox error={new Error(this.state.agent.error)} />
                : null}
            </div>
        );
    }
}
