import React from 'react';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';


export interface AddHostSuccessStepProps
{
    hostUID:    string;
}


export class AddHostSuccessStep extends React.Component<AddHostSuccessStepProps>
{
    public render ( )
    {
        return (
            <div className='AddHostSuccessStep'>
                <h1>Add a host</h1>

                <SuccessBox
                    message={`The new host was successfully created`}
                />
            </div>
        );
    }
}
