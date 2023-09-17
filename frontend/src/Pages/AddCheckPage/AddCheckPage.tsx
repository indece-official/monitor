import React from 'react';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { AddCheckStartStep } from './AddCheckStartStep/AddCheckStartStep';
import { AddCheckSuccessStep } from './AddCheckSuccessStep/AddCheckSuccessStep';
import { PageContent } from '../../Components/PageContent/PageContent';



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
        this._addOther = this._addOther.bind(this);
    }


    private _finishStart ( checkUID: string ): void
    {
        this.setState({
            step:   AddCheckStep.Success,
            checkUID
        });
    }


    private _addOther ( ): void
    {
        this.setState({
            step:       AddCheckStep.Start,
            checkUID:   null
        });
    }


    public render ( )
    {
        return (
            <PageContent>
                {this.state.step === AddCheckStep.Start ?
                    <AddCheckStartStep
                        onFinish={this._finishStart}
                        hostUID={this.props.router.params.hostUID}
                    />
                : null}
                
                {this.state.step === AddCheckStep.Success && this.state.checkUID ?
                    <AddCheckSuccessStep
                        hostUID={this.props.router.params.hostUID}
                        checkUID={this.state.checkUID}
                        onAddOther={this._addOther}
                    />
                : null}
            </PageContent>
        );
    }
}


export const AddCheckPage = withRouter($AddCheckPage);
