import React from 'react';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';


export interface AddConnectorSuccessStepProps
{
    connectorUID:   string;
}


export class AddConnectorSuccessStep extends React.Component<AddConnectorSuccessStepProps>
{
    public render ( )
    {
        return (
            <div className='AddConnectorSuccessStep'>
                <h1>Add a connector</h1>

                <SuccessBox
                    message={`The connector was successfully created and configured.`}
                />
            </div>
        );
    }
}
