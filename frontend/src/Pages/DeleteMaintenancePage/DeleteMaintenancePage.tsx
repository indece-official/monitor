import React from 'react';
import { MaintenanceService, MaintenanceV1 } from '../../Services/MaintenanceService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';


export interface DeleteMaintenancePageRouteParams
{
    maintenanceUID:    string;
}


export interface DeleteMaintenancePageProps extends RouteComponentProps<DeleteMaintenancePageRouteParams>
{
}


interface DeleteMaintenancePageState
{
    maintenance:    MaintenanceV1 | null;
    loading:        boolean;
    error:          Error | null;
    success:        string | null;
}


class $DeleteMaintenancePage extends React.Component<DeleteMaintenancePageProps, DeleteMaintenancePageState>
{
    private readonly _maintenanceService: MaintenanceService;


    constructor ( props: DeleteMaintenancePageProps )
    {
        super(props);

        this.state = {
            maintenance:    null,
            loading:        false,
            error:          null,
            success:        null
        };

        this._maintenanceService = MaintenanceService.getInstance();

        this._delete = this._delete.bind(this);
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const maintenance = await this._maintenanceService.getMaintenance(this.props.router.params.maintenanceUID);

            this.setState({
                loading:    false,
                maintenance
            });
        }
        catch ( err )
        {
            console.error(`Error loading maintenance: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    private async _delete ( ): Promise<void>
    {
        try
        {
            if ( this.state.loading || !this.state.maintenance )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            await this._maintenanceService.deleteMaintenance(this.state.maintenance.uid);

            this.setState({
                loading:    false,
                success:    'The maintenance was successfully deleted.'
            });
        }
        catch ( err )
        {
            console.error(`Error deleting maintenance: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._load();
    }


    public render ( )
    {
        return (
            <div className='AddMaintenanceStartStep'>
                <h1>Delete maintenance</h1>

                <ErrorBox error={this.state.error} />

                <div>Do you really want to delete maintenance {this.state.maintenance ? this.state.maintenance.title : '?'}?</div>

                <Button
                    onClick={this._delete}
                    disabled={this.state.loading}>
                    Delete
                </Button>

                <SuccessBox message={this.state.success} />

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}


export const DeleteMaintenancePage = withRouter($DeleteMaintenancePage);
