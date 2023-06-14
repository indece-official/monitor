import React from 'react';
import { AddUserStartStep } from './AddUserStartStep/AddUserStartStep';
import { AddUserSuccessStep } from './AddUserSuccessStep/AddUserSuccessStep';


export interface AddUserPageProps
{
}


enum AddUserStep
{
    Start   = 'start',
    Success = 'success'
}


interface AddUserPageState
{
    step:           AddUserStep;
    userUID:    string | null;
}


export class AddUserPage extends React.Component<AddUserPageProps, AddUserPageState>
{
    constructor ( props: AddUserPageProps )
    {
        super(props);

        this.state = {
            step:           AddUserStep.Start,
            userUID:    null
        };

        this._finishStart = this._finishStart.bind(this);
    }


    private _finishStart ( userUID: string ): void
    {
        this.setState({
            step:   AddUserStep.Success,
            userUID
        });
    }


    public render ( )
    {
        return (
            <div className='AddUserPage'>
                {this.state.step === AddUserStep.Start ?
                    <AddUserStartStep
                        onFinish={this._finishStart}
                    />
                : null}
                
                {this.state.step === AddUserStep.Success && this.state.userUID ?
                    <AddUserSuccessStep
                        userUID={this.state.userUID}
                    />
                : null}
            </div>
        );
    }
}
