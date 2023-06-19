import React from 'react';
import { Button } from '../../../Components/Button/Button';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';


export interface AddAgentConnectStepProps
{
    agentUID:   string;
    onFinish:   ( ) => any;
}


export class AddAgentConnectStep extends React.Component<AddAgentConnectStepProps>
{
    public render ( )
    {
        return (
            <div className='AddAgentConnectStep'>
                <h1>Add a agent</h1>

                <SuccessBox
                    message={`The agent was successfully created.`}
                />

                <Button
                    onClick={this.props.onFinish}>
                    Next
                </Button>
            </div>
        );
    }
}
