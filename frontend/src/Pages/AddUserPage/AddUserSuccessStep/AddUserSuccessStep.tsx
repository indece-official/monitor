import React from 'react';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';


export interface AddUserSuccessStepProps
{
    userUID:    string;
}


export class AddUserSuccessStep extends React.Component<AddUserSuccessStepProps>
{
    public render ( )
    {
        return (
            <div className='AddUserSuccessStep'>
                <h1>Add an user</h1>

                <SuccessBox
                    message={`The user was successfully created`}
                />
            </div>
        );
    }
}
