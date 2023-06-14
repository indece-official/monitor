import React from 'react';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { AddCheckStartStep } from './AddCheckStartStep/AddCheckStartStep';
import { AddCheckSuccessStep } from './AddCheckSuccessStep/AddCheckSuccessStep';



export interface AddCheckPageRouteParams
{
    hostUID:    string;
}


export interface AddCheckPageProps extends RouteComponentProps<AddCheckPageRouteParams>
{
}


enum AddCheckStep
{
    Start   = 'start',
    Success = 'success'
}


interface AddCheckPageState
{
    step:       AddCheckStep;
    checkUID:   string | null;
}


class $AddCheckPage extends React.Component<AddCheckPageProps, AddCheckPageState>
{
    constructor ( props: AddCheckPageProps )
    {
        super(props);

        this.state = {
            step:           AddCheckStep.Start,
            checkUID:    null
        };

        this._finishStart = this._finishStart.bind(this);
    }


    private _finishStart ( checkUID: string ): void
    {
        this.setState({
            step:   AddCheckStep.Success,
            checkUID
        });
    }


    public render ( )
    {
        return (
            <div className='AddCheckPage'>
                {this.state.step === AddCheckStep.Start ?
                    <AddCheckStartStep
                        onFinish={this._finishStart}
                        hostUID={this.props.router.params.hostUID}
                    />
                : null}
                
                {this.state.step === AddCheckStep.Success && this.state.checkUID ?
                    <AddCheckSuccessStep
                        checkUID={this.state.checkUID}
                    />
                : null}
            </div>
        );
    }
}


export const AddCheckPage = withRouter($AddCheckPage);
