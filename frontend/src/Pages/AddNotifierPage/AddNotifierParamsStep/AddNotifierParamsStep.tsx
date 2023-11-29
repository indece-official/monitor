import React from 'react';
import { FieldArray, Form, Formik } from 'formik';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { NotifierV1ConfigParams, NotifierV1ConfigParamsHttpMethod, NotifierV1ConfigParamsHttpMethods, NotifierV1Type } from '../../../Services/NotifierService';
import { faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { InputSelect } from '../../../Components/Input/InputSelect';
import { InputTextarea } from '../../../Components/Input/InputTextarea';


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


interface AddNotifierParamsStepFormDataHttpHeader
{
    name:   string;
    value:  string;
}


interface AddNotifierParamsStepFormDataHttp
{
    url:        string;
    method:     string;
    headers:    Array<AddNotifierParamsStepFormDataHttpHeader>;
    body:       string;
}


interface AddNotifierParamsStepFormDataMicrosoftTeams
{
    webhook_url:    string;
}


interface AddNotifierParamsStepFormData
{
    email_smtp:         AddNotifierParamsStepFormDataEmailSmtp;
    http:               AddNotifierParamsStepFormDataHttp;
    microsoft_teams:    AddNotifierParamsStepFormDataMicrosoftTeams;
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
                },
                http: {
                    url:        '',
                    method:     NotifierV1ConfigParamsHttpMethod.Get,
                    headers:    [],
                    body:       ''
                },
                microsoft_teams: {
                    webhook_url:    ''
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
            case NotifierV1Type.Http:
                params.http = {
                    url:        values.http.url.trim(),
                    method:     values.http.method as NotifierV1ConfigParamsHttpMethod,
                    headers:    values.http.headers.map( o => ({name: o.name.trim(), value: o.value.trim()})).filter( o => !!o.name && !!o.value ),
                    body:       values.http.body.trim() || null
                };
                break;
            case NotifierV1Type.MicrosoftTeams:
                params.microsoft_teams = {
                    webhook_url:        values.microsoft_teams.webhook_url.trim()
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
                                                    onClick={() => arrayHelpers.push('')}>
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

                            {this.props.type === NotifierV1Type.Http ?
                                <>
                                    <InputText
                                        name='http.url'
                                        label='URL'
                                        required={true}
                                    />

                                    <InputSelect
                                        name='http.method'
                                        label='Method'
                                        required={true}
                                        options={NotifierV1ConfigParamsHttpMethods.map( ( method ) => ({
                                            label:  method,
                                            value:  method
                                        }))}
                                    />

                                    <FieldArray name='http.headers'>
                                        {( arrayHelpers ) => (
                                            <div className='SetupAddHostsStep-hosts'>
                                                <Button
                                                    type='button'
                                                    onClick={() => arrayHelpers.push({name: '', value: ''})}>
                                                    <FontAwesomeIcon icon={faPlus} />
                                                    Add a header
                                                </Button>

                                                <div className='SetupAddHostsStep-hosts-list'>
                                                    {values.http.headers.length === 0 ?
                                                        <div className='SetupAddHostsStep-hosts-empty'>
                                                            No headers added yet
                                                        </div>
                                                    : null}
                                                </div>

                                                {values.http.headers.map( ( header, index ) => (
                                                    <div key={index} className='AdminAddEventPage-host'>
                                                        <InputText
                                                            name={`http.headers.${index}.name`}
                                                            label='Name'
                                                            required={true}
                                                        />

                                                        <InputText
                                                            name={`http.headers.${index}.value`}
                                                            label='Value'
                                                            required={true}
                                                        />

                                                        <Button
                                                            title='Delete header'
                                                            type='button'
                                                            onClick={() => arrayHelpers.remove(index)}>
                                                            <FontAwesomeIcon icon={faTimes} />
                                                        </Button>
                                                    </div>
                                                ))}
                                            </div>
                                        )}
                                    </FieldArray>

                                    <InputTextarea
                                        name='http.body'
                                        label='Request-Body'
                                    />
                                </>
                            : null}

                            {this.props.type === NotifierV1Type.MicrosoftTeams ?
                                <>
                                    <InputText
                                        name='microsoft_teams.webhook_url'
                                        label='Webhook-URL'
                                        required={true}
                                    />
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
