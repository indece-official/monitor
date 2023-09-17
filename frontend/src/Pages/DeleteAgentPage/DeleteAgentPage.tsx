import React from 'react';
import { AgentService, AgentV1 } from '../../Services/AgentService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { LinkUtils } from '../../utils/LinkUtils';
import { sleep } from 'ts-delay';
import { PageContent } from '../../Components/PageContent/PageContent';


export interface DeleteAgentPageRouteParams
{
    agentUID:    string;
}


export interface DeleteAgentPageProps extends RouteComponentProps<DeleteAgentPageRouteParams>
{
}


interface DeleteAgentPageState
{
    agent:   AgentV1 | null;
    loading:    boolean;
    error:      Error | null;
    success:    string | null;
}


class $DeleteAgentPage extends React.Component<DeleteAgentPageProps, DeleteAgentPageState>
{
    private readonly _agentService: AgentService;


    constructor ( props: DeleteAgentPageProps )
    {
        super(props);

        this.state = {
            agent:   null,
            loading:    false,
            error:      null,
            success:    null
        };

        this._agentService = AgentService.getInstance();

        this._delete = this._delete.bind(this);
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const agent = await this._agentService.getAgent(this.props.router.params.agentUID);

            this.setState({
                loading:    false,
                agent
            });
        }
        catch ( err )
        {
            console.error(`Error loading agent: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    private async _delete ( ): Promise<void>
    {
        try
        {
            if ( this.state.loading || !this.state.agent )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            await this._agentService.deleteAgent(this.state.agent.uid);

            this.setState({
                loading:    false,
                success:    'The agent was successfully deleted.'
            });

            await sleep(1000);

            this.props.router.navigate(LinkUtils.make('agents'));
        }
        catch ( err )
        {
            console.error(`Error deleting agent: ${(err as Error).message}`, err);

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
            <PageContent>
                <h1>Delete agent</h1>

                <ErrorBox error={this.state.error} />

                <div>Do you really want to delete agent {this.state.agent ? this.state.agent.type : '?'}?</div>

                <Button
                    onClick={this._delete}
                    disabled={this.state.loading}>
                    Delete
                </Button>

                <SuccessBox message={this.state.success} />

                <Spinner active={this.state.loading} />
            </PageContent>
        );
    }
}


export const DeleteAgentPage = withRouter($DeleteAgentPage);
