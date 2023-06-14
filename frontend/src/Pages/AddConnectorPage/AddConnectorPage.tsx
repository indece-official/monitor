import React from 'react';
import { AddConnectorStartStep } from './AddConnectorStartStep/AddConnectorStartStep';
import { AddConnectorSuccessStep } from './AddConnectorSuccessStep/AddConnectorSuccessStep';
import { AddConnectorConnectStep } from './AddConnectorConnectStep/AddConnectorConnectStep';
import { AddConnectorWaitRegisteredStep } from './AddConnectorWaitRegisteredStep/AddConnectorWaitRegisteredStep';
import { AddConnectorFailedStep } from './AddConnectorFailedStep/AddConnectorFailedStep';


export interface AddConnectorPageProps
{
}


enum AddConnectorStep
{
    Start           = 'start',
    Connect         = 'connect',
    WaitRegistered  = 'wait_registered',
    Success         = 'success',
    Failed          = 'failed'
}


interface AddConnectorPageState
{
    step:           AddConnectorStep;
    connectorUID:   string | null;
}


export class AddConnectorPage extends React.Component<AddConnectorPageProps, AddConnectorPageState>
{
    constructor ( props: AddConnectorPageProps )
    {
        super(props);

        this.state = {
            step:           AddConnectorStep.Start,
            connectorUID:   null
        };

        this._finishStart = this._finishStart.bind(this);
        this._finishConnect = this._finishConnect.bind(this);
        this._finishWaitRegistered = this._finishWaitRegistered.bind(this);
        this._fail = this._fail.bind(this);
    }


    private _finishStart ( connectorUID: string ): void
    {
        this.setState({
            step:   AddConnectorStep.Connect,
            connectorUID
        });
    }
    
    
    private _finishConnect ( ): void
    {
        this.setState({
            step:   AddConnectorStep.WaitRegistered
        });
    }
   
   
    private _finishWaitRegistered ( ): void
    {
        this.setState({
            step:   AddConnectorStep.Success
        });
    }
    
    
    private _fail ( ): void
    {
        this.setState({
            step:   AddConnectorStep.Failed
        });
    }


    public render ( )
    {
        return (
            <div className='AddConnectorPage'>
                {this.state.step === AddConnectorStep.Start ?
                    <AddConnectorStartStep
                        onFinish={this._finishStart}
                    />
                : null}
                
                {this.state.step === AddConnectorStep.Connect && this.state.connectorUID ?
                    <AddConnectorConnectStep
                        connectorUID={this.state.connectorUID}
                        onFinish={this._finishConnect}
                    />
                : null}
                
                {this.state.step === AddConnectorStep.WaitRegistered && this.state.connectorUID ?
                    <AddConnectorWaitRegisteredStep
                        connectorUID={this.state.connectorUID}
                        onFinish={this._finishWaitRegistered}
                        onFailed={this._fail}
                    />
                : null}
    
                {this.state.step === AddConnectorStep.Success && this.state.connectorUID ?
                    <AddConnectorSuccessStep
                        connectorUID={this.state.connectorUID}
                    />
                : null}
                
                {this.state.step === AddConnectorStep.Failed && this.state.connectorUID ?
                    <AddConnectorFailedStep
                        connectorUID={this.state.connectorUID}
                    />
                : null}
            </div>
        );
    }
}
