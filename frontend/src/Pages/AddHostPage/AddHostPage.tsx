import React from 'react';
import { AddHostStartStep } from './AddHostStartStep/AddHostStartStep';
import { AddHostSuccessStep } from './AddHostSuccessStep/AddHostSuccessStep';
import { PageContent } from '../../Components/PageContent/PageContent';


export interface AddHostPageProps
{
}


enum AddHostStep
{
    Start   = 'start',
    Success = 'success'
}


interface AddHostPageState
{
    step:           AddHostStep;
    hostUID:    string | null;
}


export class AddHostPage extends React.Component<AddHostPageProps, AddHostPageState>
{
    constructor ( props: AddHostPageProps )
    {
        super(props);

        this.state = {
            step:           AddHostStep.Start,
            hostUID:    null
        };

        this._finishStart = this._finishStart.bind(this);
    }


    private _finishStart ( hostUID: string ): void
    {
        this.setState({
            step:   AddHostStep.Success,
            hostUID
        });
    }


    public render ( )
    {
        return (
            <PageContent>
                {this.state.step === AddHostStep.Start ?
                    <AddHostStartStep
                        onFinish={this._finishStart}
                    />
                : null}
                
                {this.state.step === AddHostStep.Success && this.state.hostUID ?
                    <AddHostSuccessStep
                        hostUID={this.state.hostUID}
                    />
                : null}
            </PageContent>
        );
    }
}
