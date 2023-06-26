import React from 'react';
import DayJS from 'dayjs';
import { MaintenanceV1 } from '../../Services/MaintenanceService';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faWarning } from '@fortawesome/free-solid-svg-icons';


export interface MaintenanceBoxProps
{
    maintenance:    MaintenanceV1;
}


export class MaintenanceBox extends React.Component<MaintenanceBoxProps>
{
    public render ( )
    {
        return (
            <div className='MaintenanceBox'>
                <FontAwesomeIcon icon={faWarning} />

                <div className='MaintenanceBox-title'>
                    {this.props.maintenance.title}
                </div>

                <div className='MaintenanceBox-message'>
                    {this.props.maintenance.message}
                </div>

                <div className='MaintenanceBox-duration'>
                    {DayJS(this.props.maintenance.datetime_start).format()}
                   
                    {this.props.maintenance.datetime_finish ?
                        ' - ' + DayJS(this.props.maintenance.datetime_finish).format()
                    : null}
                </div>
            </div>
        );
    }
}
