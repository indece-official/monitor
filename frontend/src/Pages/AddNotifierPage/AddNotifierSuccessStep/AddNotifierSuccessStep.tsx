import React from 'react';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';


export interface AddNotifierSuccessStepProps
{
    notifierUID:    string;
}


export class AddNotifierSuccessStep extends React.Component<AddNotifierSuccessStepProps>
{
    public render ( )
    {
        return (
            <div className='AddNotifierSuccessStep'>
                <h1>Add a notifier</h1>

                <SuccessBox
                    message={`The new notifier was successfully created`}
                />
            </div>
        );
    }
}
