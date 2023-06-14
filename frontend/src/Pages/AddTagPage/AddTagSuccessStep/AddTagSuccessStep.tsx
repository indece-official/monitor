import React from 'react';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';


export interface AddTagSuccessStepProps
{
    tagUID:    string;
}


export class AddTagSuccessStep extends React.Component<AddTagSuccessStepProps>
{
    public render ( )
    {
        return (
            <div className='AddTagSuccessStep'>
                <h1>Add a tag</h1>

                <SuccessBox
                    message={`The new tag was successfully created`}
                />
            </div>
        );
    }
}
