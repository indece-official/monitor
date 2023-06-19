import React from 'react';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';


export interface AddAgentSuccessStepProps
{
    agentUID:   string;
}


export class AddAgentSuccessStep extends React.Component<AddAgentSuccessStepProps>
{
    public render ( )
    {
        return (
            <div className='AddAgentSuccessStep'>
                <h1>Add a agent</h1>

                <SuccessBox
                    message={`The agent was successfully created and configured.`}
                />
            </div>
        );
    }
}
