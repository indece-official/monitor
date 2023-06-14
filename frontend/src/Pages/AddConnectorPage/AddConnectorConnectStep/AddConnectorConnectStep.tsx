import React from 'react';
import { Button } from '../../../Components/Button/Button';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';


export interface AddConnectorConnectStepProps
{
    connectorUID:   string;
    onFinish:       ( ) => any;
}


export class AddConnectorConnectStep extends React.Component<AddConnectorConnectStepProps>
{
    public render ( )
    {
        return (
            <div className='AddConnectorConnectStep'>
                <h1>Add a connector</h1>

                <SuccessBox
                    message={`The connector was successfully created.`}
                />

                <Button
                    onClick={this.props.onFinish}>
                    Next
                </Button>
            </div>
        );
    }
}
