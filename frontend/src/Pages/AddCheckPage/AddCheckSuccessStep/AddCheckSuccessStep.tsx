import React from 'react';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';


export interface AddCheckSuccessStepProps
{
    checkUID:    string;
}


export class AddCheckSuccessStep extends React.Component<AddCheckSuccessStepProps>
{
    public render ( )
    {
        return (
            <div className='AddCheckSuccessStep'>
                <h1>New check</h1>

                <SuccessBox
                    message={`The new check was successfully created`}
                />
            </div>
        );
    }
}
