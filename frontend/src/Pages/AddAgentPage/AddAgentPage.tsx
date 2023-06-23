import React from 'react';
import { AddAgentStartStep } from './AddAgentStartStep/AddAgentStartStep';
import { AddAgentSuccessStep } from './AddAgentSuccessStep/AddAgentSuccessStep';
import { AddAgentConnectStep } from './AddAgentConnectStep/AddAgentConnectStep';
import { AddAgentWaitRegisteredStep } from './AddAgentWaitRegisteredStep/AddAgentWaitRegisteredStep';
import { AddAgentFailedStep } from './AddAgentFailedStep/AddAgentFailedStep';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';


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
    Connect         = 'connect',
    WaitRegistered  = 'wait_registered',
    Success         = 'success',
    Failed          = 'failed'
}


interface AddAgentPageState
{
    step:           AddAgentStep;
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
        this._finishConnect = this._finishConnect.bind(this);
        this._finishWaitRegistered = this._finishWaitRegistered.bind(this);
        this._fail = this._fail.bind(this);
    }


    private _finishStart ( agentUID: string ): void
    {
        this.setState({
            step:   AddAgentStep.Connect,
            agentUID
        });
    }
    
    
    private _finishConnect ( ): void
    {
        this.setState({
            step:   AddAgentStep.WaitRegistered
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
            <div className='AddAgentPage'>
                {this.state.step === AddAgentStep.Start ?
                    <AddAgentStartStep
                        hostUID={this.props.router.params.hostUID}
                        onFinish={this._finishStart}
                    />
                : null}
                
                {this.state.step === AddAgentStep.Connect && this.state.agentUID ?
                    <AddAgentConnectStep
                        agentUID={this.state.agentUID}
                        onFinish={this._finishConnect}
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
            </div>
        );
    }
}


export const AddAgentPage = withRouter($AddAgentPage);
