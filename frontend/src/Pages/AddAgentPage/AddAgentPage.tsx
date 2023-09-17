import React from 'react';
import { AddAgentStartStep } from './AddAgentStartStep/AddAgentStartStep';
import { AddAgentSuccessStep } from './AddAgentSuccessStep/AddAgentSuccessStep';
import { AddAgentWaitRegisteredStep } from './AddAgentWaitRegisteredStep/AddAgentWaitRegisteredStep';
import { AddAgentFailedStep } from './AddAgentFailedStep/AddAgentFailedStep';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { PageContent } from '../../Components/PageContent/PageContent';


export interface AddAgentPageRouteParams
{
    hostUID?:   string;
}


export interface AddAgentPageProps extends RouteComponentProps<AddAgentPageRouteParams>
{
}


enum AddAgentStep
{
    Start           = 'start',
    WaitRegistered  = 'wait_registered',
    Success         = 'success',
    Failed          = 'failed'
}


interface AddAgentPageState
{
    step:       AddAgentStep;
    agentUID:   string | null;
}


class $AddAgentPage extends React.Component<AddAgentPageProps, AddAgentPageState>
{
    constructor ( props: AddAgentPageProps )
    {
        super(props);

        this.state = {
            step:           AddAgentStep.Start,
            agentUID:   null
        };

        this._finishStart = this._finishStart.bind(this);
        this._finishWaitRegistered = this._finishWaitRegistered.bind(this);
        this._fail = this._fail.bind(this);
    }


    private _finishStart ( agentUID: string ): void
    {
        this.setState({
            step:   AddAgentStep.WaitRegistered,
            agentUID
        });
    }

   
    private _finishWaitRegistered ( ): void
    {
        this.setState({
            step:   AddAgentStep.Success
        });
    }
    
    
    private _fail ( ): void
    {
        this.setState({
            step:   AddAgentStep.Failed
        });
    }


    public render ( )
    {
        return (
            <PageContent>
                {this.state.step === AddAgentStep.Start ?
                    <AddAgentStartStep
                        hostUID={this.props.router.params.hostUID}
                        onFinish={this._finishStart}
                    />
                : null}

                {this.state.step === AddAgentStep.WaitRegistered && this.state.agentUID ?
                    <AddAgentWaitRegisteredStep
                        agentUID={this.state.agentUID}
                        onFinish={this._finishWaitRegistered}
                        onFailed={this._fail}
                    />
                : null}
    
                {this.state.step === AddAgentStep.Success && this.state.agentUID ?
                    <AddAgentSuccessStep
                        agentUID={this.state.agentUID}
                    />
                : null}
                
                {this.state.step === AddAgentStep.Failed && this.state.agentUID ?
                    <AddAgentFailedStep
                        agentUID={this.state.agentUID}
                    />
                : null}
            </PageContent>
        );
    }
}


export const AddAgentPage = withRouter($AddAgentPage);
