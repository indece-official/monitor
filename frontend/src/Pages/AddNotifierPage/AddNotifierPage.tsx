import React from 'react';
import { AddNotifierNameTypeStep } from './AddNotifierNameTypeStep/AddNotifierNameTypeStep';
import { AddNotifierSuccessStep } from './AddNotifierSuccessStep/AddNotifierSuccessStep';
import { AddNotifierV1Request, NotifierV1ConfigParams, NotifierV1Filter, NotifierV1Type } from '../../Services/NotifierService';
import { AddNotifierCreateStep } from './AddNotifierCreateStep/AddNotifierCreateStep';
import { AddNotifierFiltersStep } from './AddNotifierFiltersStep/AddNotifierFiltersStep';
import { AddNotifierParamsStep } from './AddNotifierParamsStep/AddNotifierParamsStep';


export interface AddNotifierPageProps
{
}


enum AddNotifierStep
{
    NameType    = 'name_type',
    Params      = 'params',
    Filters     = 'filters',
    Create      = 'create',
    Success     = 'success'
}


interface AddNotifierPageState
{
    step:           AddNotifierStep;
    addNotifer:     Partial<AddNotifierV1Request>;
    notifierUID:    string | null;
}


export class AddNotifierPage extends React.Component<AddNotifierPageProps, AddNotifierPageState>
{
    constructor ( props: AddNotifierPageProps )
    {
        super(props);

        this.state = {
            step:           AddNotifierStep.NameType,
            addNotifer:     {},
            notifierUID:    null
        };

        this._finishNameType = this._finishNameType.bind(this);
        this._finishParams = this._finishParams.bind(this);
        this._finishFilters = this._finishFilters.bind(this);
        this._finishCreate = this._finishCreate.bind(this);
    }


    private _finishNameType ( name: string, type: NotifierV1Type ): void
    {
        this.setState({
            step:           AddNotifierStep.Params,
            addNotifer: {
                name,
                type
            }
        });
    }


    private _finishParams ( params: NotifierV1ConfigParams ): void
    {
        this.setState({
            step:           AddNotifierStep.Filters,
            addNotifer: {
                ...this.state.addNotifer,
                config: {
                    filters: this.state.addNotifer.config?.filters || [],
                    params
                }
            }
        });
    }


    private _finishFilters ( filters: Array<NotifierV1Filter> ): void
    {
        this.setState({
            step:           AddNotifierStep.Create,
            addNotifer: {
                ...this.state.addNotifer,
                config: {
                    params: this.state.addNotifer.config?.params || {},
                    filters
                }
            }
        });
    }


    private _finishCreate ( notifierUID: string ): void
    {
        this.setState({
            step:           AddNotifierStep.Success,
            notifierUID
        });
    }


    public render ( )
    {
        return (
            <div className='AddNotifierPage'>
                {this.state.step === AddNotifierStep.NameType ?
                    <AddNotifierNameTypeStep
                        onFinish={this._finishNameType}
                    />
                : null}

                {this.state.step === AddNotifierStep.Params && this.state.addNotifer.type ?
                    <AddNotifierParamsStep
                        type={this.state.addNotifer.type}
                        onFinish={this._finishParams}
                    />
                : null}
                
                {this.state.step === AddNotifierStep.Filters ?
                    <AddNotifierFiltersStep
                        onFinish={this._finishFilters}
                    />
                : null}

                {this.state.step === AddNotifierStep.Create ?
                    <AddNotifierCreateStep
                        addNotifier={this.state.addNotifer as AddNotifierV1Request}
                        onFinish={this._finishCreate}
                    />
                : null}
                
                {this.state.step === AddNotifierStep.Success && this.state.notifierUID ?
                    <AddNotifierSuccessStep
                        notifierUID={this.state.notifierUID}
                    />
                : null}
            </div>
        );
    }
}
