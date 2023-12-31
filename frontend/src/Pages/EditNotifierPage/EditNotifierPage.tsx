import React from 'react';
import { NotifierService, NotifierV1, NotifierV1ConfigParams, NotifierV1ConfigParamsHttpMethod, NotifierV1ConfigParamsHttpMethods, NotifierV1Filter, NotifierV1Type } from '../../Services/NotifierService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { FieldArray, Form, Formik } from 'formik';
import { InputText } from '../../Components/Input/InputText';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';
import { InputMultiSelect } from '../../Components/Input/InputMultiSelect';
import { InputCheckbox } from '../../Components/Input/InputCheckbox';
import { TagService, TagV1 } from '../../Services/TagService';
import { sleep } from 'ts-delay';
import { PageContent } from '../../Components/PageContent/PageContent';
import { InputTextarea } from '../../Components/Input/InputTextarea';
import { InputSelect } from '../../Components/Input/InputSelect';


export interface EditNotifierPageRouteParams
{
    notifierUID:    string;
}


export interface EditNotifierPageProps extends RouteComponentProps<EditNotifierPageRouteParams>
{
}


interface EditNotifierPageFormDataParamsEmailSmtp
{
    host:       string;
    port:       string;
    user:       string;
    password:   string;
    from:       string;
    to:         Array<string>;
}


interface EditNotifierParamsStepFormDataHttpHeader
{
    name:   string;
    value:  string;
}


interface EditNotifierParamsStepFormDataHttp
{
    url:        string;
    method:     string;
    headers:    Array<EditNotifierParamsStepFormDataHttpHeader>;
    body:       string;
}


interface EditNotifierParamsStepFormDataMicrosoftTeams
{
    webhook_url:    string;
}



interface EditNotifierPageFormDataParams
{
    email_smtp:         EditNotifierPageFormDataParamsEmailSmtp;
    http:               EditNotifierParamsStepFormDataHttp;
    microsoft_teams:    EditNotifierParamsStepFormDataMicrosoftTeams;
}


interface EditNotifierPageFormDataFilter
{
    tag_uids:       Array<string>;
    critical:       boolean;
    warning:        boolean;
    unknown:        boolean;
    decline:        boolean;
    min_duration:   string;
}


interface EditNotifierPageFormData
{
    name:       string;
    params:     EditNotifierPageFormDataParams;
    filters:    Array<EditNotifierPageFormDataFilter>;
}


interface EditNotifierPageState
{
    initialFormData:    EditNotifierPageFormData;
    notifier:           NotifierV1 | null;
    tags:               Array<TagV1>
    loading:            boolean;
    error:              Error | null;
    success:            string | null;
}


class $EditNotifierPage extends React.Component<EditNotifierPageProps, EditNotifierPageState>
{
    private readonly _notifierService: NotifierService;
    private readonly _tagService: TagService;


    constructor ( props: EditNotifierPageProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                name:       '',
                params:     {
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
                        method:     '',
                        headers:    [],
                        body:       ''
                    },
                    microsoft_teams: {
                        webhook_url:    ''
                    }
                },
                filters:    []
            },
            notifier:       null,
            tags:           [],
            loading:        false,
            error:          null,
            success:        null
        };

        this._notifierService = NotifierService.getInstance();
        this._tagService = TagService.getInstance();

        this._submit = this._submit.bind(this);
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const tags = await this._tagService.getTags();
            const notifier = await this._notifierService.getNotifier(this.props.router.params.notifierUID);

            this.setState({
                loading:    false,
                notifier,
                tags,
                initialFormData: {
                    name: notifier.name,
                    params: {
                        email_smtp: {
                            host:       notifier.config?.params.email_smtp?.host || '',
                            port:       '' + (notifier.config?.params.email_smtp?.port || ''),
                            user:       notifier.config?.params.email_smtp?.user || '',
                            password:   notifier.config?.params.email_smtp?.password || '',
                            from:       notifier.config?.params.email_smtp?.from || '',
                            to:         notifier.config?.params.email_smtp?.to || [],
                        },
                        http: {
                            url:        notifier.config?.params.http?.url || '',
                            method:     notifier.config?.params.http?.method || '',
                            headers:    (notifier.config?.params.http?.headers || []).map( ( o ) =>({
                                name:   o.name,
                                value:  o.value
                            })),
                            body:       notifier.config?.params.http?.body || '',
                        },
                        microsoft_teams: {
                            webhook_url:        notifier.config?.params.microsoft_teams?.webhook_url || '',
                        }
                    },
                    filters: (notifier.config?.filters || []).map( ( filter ) => ({
                        tag_uids:       filter.tag_uids,
                        critical:       filter.critical,
                        warning:        filter.warning,
                        unknown:        filter.unknown,
                        decline:        filter.decline,
                        min_duration:   filter.min_duration
                    }))
                }
            });
        }
        catch ( err )
        {
            console.error(`Error loading notifier: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    private async _submit ( values: EditNotifierPageFormData ): Promise<void>
    {
        try
        {
            if ( this.state.loading || !this.state.notifier )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            const params: NotifierV1ConfigParams = {};

            switch ( this.state.notifier.type )
            {
                case NotifierV1Type.EmailSmtp:
                    params.email_smtp = {
                        host:       values.params.email_smtp.host.trim(),
                        port:       parseInt(values.params.email_smtp.port, 10),
                        user:       values.params.email_smtp.user.trim(),
                        password:   values.params.email_smtp.password.trim(),
                        from:       values.params.email_smtp.from.trim(),
                        to:         values.params.email_smtp.to.map( s => s.trim() ).filter( s => !!s ),
                    };
                    break;
                case NotifierV1Type.Http:
                    params.http = {
                        url:        values.params.http.url.trim(),
                        method:     values.params.http.method as NotifierV1ConfigParamsHttpMethod,
                        headers:    values.params.http.headers.map( o => ({name: o.name.trim(), value: o.value.trim()})).filter( o => !!o.name && !!o.value ),
                        body:       values.params.http.body.trim() || null
                    };
                    break;
                case NotifierV1Type.MicrosoftTeams:
                    params.microsoft_teams = {
                        webhook_url:        values.params.microsoft_teams.webhook_url.trim()
                    };
                    break;
                // TODO: Error on default
            }

            const filters: Array<NotifierV1Filter> = values.filters.map( ( filterValues ) => ({
                tag_uids:       filterValues.tag_uids.filter( s => !!s ),
                critical:       filterValues.critical,
                warning:        filterValues.warning,
                unknown:        filterValues.unknown,
                decline:        filterValues.decline,
                min_duration:   filterValues.min_duration.trim()
            }));

            await this._notifierService.updateNotifier(
                this.state.notifier.uid,
                {
                    name:       values.name.trim(),
                    config: {
                        params,
                        filters
                    }
                }
            );

            this.setState({
                loading:    false,
                success:    'The notifier was successfully updated.'
            });

            await sleep(1000);

            this.props.router.navigate(-1);
        }
        catch ( err )
        {
            console.error(`Error updating notifier: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._load();
    }


    public render ( )
    {
        return (
            <PageContent>
                <h1>Edit notifier</h1>

                <ErrorBox error={this.state.error} />

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    {({ values }) => (
                        <Form>
                            <InputText
                                name='name'
                                label='Name'
                                required={true}
                            />

                            {this.state.notifier?.type === NotifierV1Type.EmailSmtp ?
                                <>
                                    <InputText
                                        name='params.email_smtp.host'
                                        label='SMTP-Host'
                                        required={true}
                                    />

                                    <InputText
                                        name='params.email_smtp.port'
                                        label='SMTP-Port'
                                        required={true}
                                    />

                                    <InputText
                                        name='params.email_smtp.user'
                                        label='SMTP-Username'
                                        required={true}
                                    />

                                    <InputText
                                        type='password'
                                        name='params.email_smtp.password'
                                        label='SMTP-Password'
                                        required={true}
                                    />

                                    <InputText
                                        name='params.email_smtp.from'
                                        label='From'
                                        required={true}
                                    />

                                    <FieldArray name='params.email_smtp.to'>
                                        {( arrayHelpers ) => (
                                            <div className='SetupAddHostsStep-hosts'>
                                                <Button
                                                    type='button'
                                                    onClick={() => arrayHelpers.push('')}>
                                                    <FontAwesomeIcon icon={faPlus} />
                                                    Add a receiver
                                                </Button>

                                                <div className='SetupAddHostsStep-hosts-list'>
                                                    {values.params.email_smtp.to.length === 0 ?
                                                        <div className='SetupAddHostsStep-hosts-empty'>
                                                            No receivers added yet
                                                        </div>
                                                    : null}
                                                </div>

                                                {values.params.email_smtp.to.map( ( to, index ) => (
                                                    <div key={index} className='AdminAddEventPage-host'>
                                                        <InputText
                                                            name={`params.email_smtp.to.${index}`}
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

                            {this.state.notifier?.type === NotifierV1Type.Http ?
                                <>
                                    <InputText
                                        name='params.http.url'
                                        label='URL'
                                        required={true}
                                    />

                                    <InputSelect
                                        name='params.http.method'
                                        label='Method'
                                        required={true}
                                        options={NotifierV1ConfigParamsHttpMethods.map( ( method ) => ({
                                            label:  method,
                                            value:  method
                                        }))}
                                    />

                                    <FieldArray name='params.http.headers'>
                                        {( arrayHelpers ) => (
                                            <div className='SetupAddHostsStep-hosts'>
                                                <Button
                                                    type='button'
                                                    onClick={() => arrayHelpers.push({name: '', value: ''})}>
                                                    <FontAwesomeIcon icon={faPlus} />
                                                    Add a header
                                                </Button>

                                                <div className='SetupAddHostsStep-hosts-list'>
                                                    {values.params.http.headers.length === 0 ?
                                                        <div className='SetupAddHostsStep-hosts-empty'>
                                                            No headers added yet
                                                        </div>
                                                    : null}
                                                </div>

                                                {values.params.http.headers.map( ( header, index ) => (
                                                    <div key={index} className='AdminAddEventPage-host'>
                                                        <InputText
                                                            name={`params.http.headers.${index}.name`}
                                                            label='Name'
                                                            required={true}
                                                        />

                                                        <InputText
                                                            name={`params.http.headers.${index}.value`}
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
                                        name='params.http.body'
                                        label='Request-Body'
                                    />
                                </>
                            : null}

                            {this.state.notifier?.type === NotifierV1Type.MicrosoftTeams ?
                                <>
                                    <InputText
                                        name='params.microsoft_teams.webhook_url'
                                        label='Webhook-URL'
                                        required={true}
                                    />
                                </>
                            : null}

                            <FieldArray name='filters'>
                                {( arrayHelpers ) => (
                                    <div className='SetupAddHostsStep-hosts'>
                                        <Button
                                            type='button'
                                            onClick={() => arrayHelpers.push({
                                                tag_uids:       [],
                                                critical:       true,
                                                warning:        true,
                                                unknown:        true,
                                                decline:        true,
                                                min_duration:   '5m'
                                            })}>
                                            <FontAwesomeIcon icon={faPlus} />
                                            Add a filter
                                        </Button>

                                        <div className='SetupAddHostsStep-hosts-list'>
                                            {values.filters.length === 0 ?
                                                <div className='SetupAddHostsStep-hosts-empty'>
                                                    No filters added yet
                                                </div>
                                            : null}
                                        </div>

                                        {values.filters.map( ( filter, index ) => (
                                            <div key={index} className='AdminAddEventPage-host'>
                                                <InputMultiSelect
                                                    name={`filters.${index}.tag_uids`}
                                                    label='Tags'
                                                    options={this.state.tags.map( ( tag ) => ({
                                                        label:  tag.name,
                                                        value:  tag.uid
                                                    }))}
                                                />

                                                <InputCheckbox
                                                    name={`filters.${index}.critical`}
                                                    label='Notify on critical'
                                                />

                                                <InputCheckbox
                                                    name={`filters.${index}.warning`}
                                                    label='Notify on warning'
                                                />

                                                <InputCheckbox
                                                    name={`filters.${index}.unknown`}
                                                    label='Notify on unknown'
                                                />

                                                <InputCheckbox
                                                    name={`filters.${index}.decline`}
                                                    label='Notify on decline'
                                                />

                                                <InputText
                                                    name={`filters.${index}.min_duration`}
                                                    label='Min. duration'
                                                    required={true}
                                                />

                                                <Button
                                                    title='Delete filter'
                                                    type='button'
                                                    onClick={() => arrayHelpers.remove(index)}>
                                                    <FontAwesomeIcon icon={faTimes} />
                                                </Button>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </FieldArray>

                            <Button
                                type='submit'
                                disabled={this.state.loading}>
                                Save
                            </Button>
                        </Form>
                    )}
                </Formik>

                <SuccessBox message={this.state.success} />

                <Spinner active={this.state.loading} />
            </PageContent>
        );
    }
}


export const EditNotifierPage = withRouter($EditNotifierPage);
