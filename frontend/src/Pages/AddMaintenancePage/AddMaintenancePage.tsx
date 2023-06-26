import React from 'react';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { AddMaintenanceStartStep } from './AddMaintenanceStartStep/AddMaintenanceStartStep';
import { AddMaintenanceSuccessStep } from './AddMaintenanceSuccessStep/AddMaintenanceSuccessStep';



export interface AddMaintenancePageRouteParams
{
    hostUID:    string;
}


export interface AddMaintenancePageProps extends RouteComponentProps<AddMaintenancePageRouteParams>
{
}


enum AddMaintenanceStep
{
    Start   = 'start',
    Success = 'success'
}


interface AddMaintenancePageState
{
    step:           AddMaintenanceStep;
    maintenanceUID: string | null;
}


class $AddMaintenancePage extends React.Component<AddMaintenancePageProps, AddMaintenancePageState>
{
    constructor ( props: AddMaintenancePageProps )
    {
        super(props);

        this.state = {
            step:           AddMaintenanceStep.Start,
            maintenanceUID: null
        };

        this._finishStart = this._finishStart.bind(this);
    }


    private _finishStart ( maintenanceUID: string ): void
    {
        this.setState({
            step:   AddMaintenanceStep.Success,
            maintenanceUID
        });
    }


    public render ( )
    {
        return (
            <div className='AddMaintenancePage'>
                {this.state.step === AddMaintenanceStep.Start ?
                    <AddMaintenanceStartStep
                        onFinish={this._finishStart}
                        hostUID={this.props.router.params.hostUID}
                    />
                : null}
                
                {this.state.step === AddMaintenanceStep.Success && this.state.maintenanceUID ?
                    <AddMaintenanceSuccessStep
                        maintenanceUID={this.state.maintenanceUID}
                    />
                : null}
            </div>
        );
    }
}


export const AddMaintenancePage = withRouter($AddMaintenancePage);
