import React from 'react';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';
import { Button } from '../../../Components/Button/Button';
import { LinkUtils } from '../../../utils/LinkUtils';


export interface AddCheckSuccessStepProps
{
    hostUID:    string;
    checkUID:   string;
    onAddOther: ( ) => any;
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

                <Button onClick={this.props.onAddOther}>
                    Add an other check
                </Button>

                <Button to={LinkUtils.make('host', this.props.hostUID)}>
                    Show host
                </Button>
            </div>
        );
    }
}
