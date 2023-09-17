import React from 'react';
import { SuccessBox } from '../../../Components/SuccessBox/SuccessBox';
import { Button } from '../../../Components/Button/Button';
import { LinkUtils } from '../../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus } from '@fortawesome/free-solid-svg-icons';


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

                <Button to={LinkUtils.make('host', this.props.hostUID, 'agent', 'add')}>
                    <FontAwesomeIcon icon={faPlus} /> Add an agent
                </Button>
            </div>
        );
    }
}
