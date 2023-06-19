import React from 'react';
import { FieldArray, Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { NotifierService, NotifierV1ConfigParams, NotifierV1Type, NotifierV1Types } from '../../../Services/NotifierService';
import { InputMultiSelect } from '../../../Components/Input/InputMultiSelect';
import { TagService, TagV1 } from '../../../Services/TagService';
import { InputSelect } from '../../../Components/Input/InputSelect';
import { faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';


export interface AddNotifierParamsStepProps
{
    type:       NotifierV1Type;
    onFinish:   ( params: NotifierV1ConfigParams ) => any;
}


interface AddNotifierParamsStepFormDataEmailSmtp
{
    host:       string;
    port:       string;
    user:       string;
    password:   string;
    from:       string;
    to:         Array<string>;
}


interface AddNotifierParamsStepFormData
{
    email_smtp:     AddNotifierParamsStepFormDataEmailSmtp;
}


interface AddNotifierParamsStepState
{
    initialFormData:    AddNotifierParamsStepFormData;
}


export class AddNotifierParamsStep extends React.Component<AddNotifierParamsStepProps, AddNotifierParamsStepState>
{
    constructor ( props: AddNotifierParamsStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                email_smtp: {
                    host:       '',
                    port:       '465',
                    user:       '',
                    password:   '',
                    from:       '',
                    to:         []
                }
            }
        };

        this._submit = this._submit.bind(this);
    }


    private async _submit ( values: AddNotifierParamsStepFormData ): Promise<void>
    {
        const params: NotifierV1ConfigParams = {};

        switch ( this.props.type )
        {
            case NotifierV1Type.EmailSmtp:
                params.email_smtp = {
                    host:       values.email_smtp.host.trim(),
                    port:       parseInt(values.email_smtp.port, 10),
                    user:       values.email_smtp.user.trim(),
                    password:   values.email_smtp.password.trim(),
                    from:       values.email_smtp.from.trim(),
                    to:         values.email_smtp.to.map( s => s.trim() ).filter( s => !!s ),
                };
                break;
            // TODO: Error on default
        }

        this.props.onFinish(params);
    }


    public render ( )
    {
        return (
            <div className='AddNotifierParamsStep'>
                <h1>Add a notifier</h1>

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    {({ values }) => (
                        <Form>
                            {this.props.type === NotifierV1Type.EmailSmtp ?
                                <>
                                    <InputText
                                        name='email_smtp.host'
                                        label='SMTP-Host'
                                        required={true}
                                    />

                                    <InputText
                                        name='email_smtp.port'
                                        label='SMTP-Port'
                                        required={true}
                                    />

                                    <InputText
                                        name='email_smtp.user'
                                        label='SMTP-Username'
                                        required={true}
                                    />

                                    <InputText
                                        type='password'
                                        name='email_smtp.password'
                                        label='SMTP-Password'
                                        required={true}
                                    />

                                    <InputText
                                        name='email_smtp.from'
                                        label='From'
                                        required={true}
                                    />

                                    <FieldArray name='email_smtp.to'>
                                        {( arrayHelpers ) => (
                                            <div className='SetupAddHostsStep-hosts'>
                                                <Button
                                                    type='button'
                                                    onClick={() => arrayHelpers.push({
                                                        labels: []
                                                    })}>
                                                    <FontAwesomeIcon icon={faPlus} />
                                                    Add a receiver
                                                </Button>

                                                <div className='SetupAddHostsStep-hosts-list'>
                                                    {values.email_smtp.to.length === 0 ?
                                                        <div className='SetupAddHostsStep-hosts-empty'>
                                                            No receivers added yet
                                                        </div>
                                                    : null}
                                                </div>

                                                {values.email_smtp.to.map( ( to, index ) => (
                                                    <div key={index} className='AdminAddEventPage-host'>
                                                        <InputText
                                                            name={`email_smtp.to.${index}`}
                                                            label='To'
                                                            required={true}
                                                        />

                                                        <Button
                                                            title='Delete receiver'
                                                            type='button'
                                                            onClick={() => arrayHelpers.remove(index)}>
                                                            <FontAwesomeIcon icon={faTimes} />
                                                        </Button>
                                                    </div>
                                                ))}
                                            </div>
                                        )}
                                    </FieldArray>
                                </>
                            : null}

                            <Button type='submit'>
                                Continue
                            </Button>
                        </Form>
                    )}
                </Formik>
            </div>
        );
    }
}
