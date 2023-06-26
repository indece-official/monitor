import React from 'react';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';


export interface AddMaintenanceSuccessStepProps
{
    maintenanceUID:    string;
}


export class AddMaintenanceSuccessStep extends React.Component<AddMaintenanceSuccessStepProps>
{
    public render ( )
    {
        return (
            <div className='AddMaintenanceSuccessStep'>
                <h1>New maintenance</h1>

                <SuccessBox
                    message={`The new maintenance was successfully created`}
                />
            </div>
        );
    }
}
