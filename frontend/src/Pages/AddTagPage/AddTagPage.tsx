import React from 'react';
import { AddTagStartStep } from './AddTagStartStep/AddTagStartStep';
import { AddTagSuccessStep } from './AddTagSuccessStep/AddTagSuccessStep';


export interface AddTagPageProps
{
}


enum AddTagStep
{
    Start   = 'start',
    Success = 'success'
}


interface AddTagPageState
{
    step:           AddTagStep;
    tagUID:    string | null;
}


export class AddTagPage extends React.Component<AddTagPageProps, AddTagPageState>
{
    constructor ( props: AddTagPageProps )
    {
        super(props);

        this.state = {
            step:           AddTagStep.Start,
            tagUID:    null
        };

        this._finishStart = this._finishStart.bind(this);
    }


    private _finishStart ( tagUID: string ): void
    {
        this.setState({
            step:   AddTagStep.Success,
            tagUID
        });
    }


    public render ( )
    {
        return (
            <div className='AddTagPage'>
                {this.state.step === AddTagStep.Start ?
                    <AddTagStartStep
                        onFinish={this._finishStart}
                    />
                : null}
                
                {this.state.step === AddTagStep.Success && this.state.tagUID ?
                    <AddTagSuccessStep
                        tagUID={this.state.tagUID}
                    />
                : null}
            </div>
        );
    }
}
