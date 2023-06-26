import React from 'react';
import DayJS from 'dayjs';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { MaintenanceService, MaintenanceV1 } from '../../Services/MaintenanceService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPen, faTrash } from '@fortawesome/free-solid-svg-icons';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { LabelValueList } from '../../Components/LabelValueList/LabelValueList';
import { LabelValue } from '../../Components/LabelValueList/LabelValue';


export interface MaintenancePageRouteParams
{
    maintenanceUID:    string;
}


export interface MaintenancePageProps extends RouteComponentProps<MaintenancePageRouteParams>
{
}


interface MaintenancePageState
{
    maintenance:        MaintenanceV1 | null;
    loading:    boolean;
    error:      Error | null;
}


class $MaintenancePage extends React.Component<MaintenancePageProps, MaintenancePageState>
{
    private readonly _maintenanceService: MaintenanceService;


    constructor ( props: MaintenancePageProps )
    {
        super(props);

        this.state = {
            maintenance:       null,
            loading:    false,
            error:      null
        };

        this._maintenanceService  = MaintenanceService.getInstance();
    }


    private async _loadMaintenance ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const maintenance = await this._maintenanceService.getMaintenance(
                this.props.router.params.maintenanceUID
            );

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


    public async componentDidMount ( ): Promise<void>
    {
        await this._loadMaintenance();
    }


    public render ( )
    {
        return (
            <div className='MaintenancePage'>
                <h1>Maintenance</h1>

                <ErrorBox error={this.state.error} />

                {this.state.maintenance ?
                    <div className='MaintenancePage-actions'>
                        <Button to={LinkUtils.make('maintenance', this.state.maintenance.uid, 'edit')}>
                            <FontAwesomeIcon icon={faPen} /> Edit
                        </Button>
                    
                        <Button to={LinkUtils.make('maintenance', this.state.maintenance.uid, 'delete')}>
                            <FontAwesomeIcon icon={faTrash} /> Delete
                        </Button>
                    </div>
                : null}

                {this.state.maintenance ?
                    <LabelValueList>
                        <LabelValue
                            label='Title'
                            value={this.state.maintenance.title}
                        />

                        <LabelValue
                            label='Message'
                            value={this.state.maintenance.message || '-'}
                        />

                        <LabelValue
                            label='Start'
                            value={DayJS(this.state.maintenance.datetime_start).format()}
                        />
                    </LabelValueList>
                : null}

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}


export const MaintenancePage = withRouter($MaintenancePage);
