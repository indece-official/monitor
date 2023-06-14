import React from 'react';
import { Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { CheckService } from '../../../Services/CheckService';
import { HostService, HostV1 } from '../../../Services/HostService';
import { CheckerService, CheckerV1 } from '../../../Services/CheckerService';
import { InputSelect } from '../../../Components/Input/InputSelect';
import { ConnectorService } from '../../../Services/ConnectorService';


export interface AddCheckStartStepProps
{
    hostUID:    string;
    onFinish:   ( checkUID: string ) => any;
}


interface AddCheckStartStepFormData
{
    name:           string;
    checker_uid:    string;
    schedule:       string;
    params:         Record<string, string>;
}


interface AddCheckStartStepState
{
    initialFormData:    AddCheckStartStepFormData;
    host:               HostV1 | null;
    checkers:           Array<CheckerV1>;
    loading:            boolean;
    error:              Error | null;
}


export class AddCheckStartStep extends React.Component<AddCheckStartStepProps, AddCheckStartStepState>
{
    private readonly _hostService: HostService;
    private readonly _connectorService: ConnectorService;
    private readonly _checkerService: CheckerService;
    private readonly _checkService: CheckService;


    constructor ( props: AddCheckStartStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                name:           '',
                checker_uid:    '',
                schedule:       '',
                params:         {}
            },
            host:       null,
            checkers:   [],
            loading:    false,
            error:      null
        };

        this._hostService = HostService.getInstance();
        this._connectorService = ConnectorService.getInstance();
        this._checkerService = CheckerService.getInstance();
        this._checkService = CheckService.getInstance();

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

            const host = await this._hostService.getHost(this.props.hostUID);
            let checkers = await this._checkerService.getCheckers();
            let connectors = await this._connectorService.getConnectors();

            connectors = connectors.filter( ( connector ) => connector.host_uid === host.uid );
            checkers = checkers.filter( ( checker ) => checker.custom_checks && connectors.find( connector => connector.type === checker.connector_type ) );

            this.setState({
                loading:    false,
                host,
                checkers
            });
        }
        catch ( err )
        {
            console.error(`Error loading hosts and checkers: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    private async _submit ( values: AddCheckStartStepFormData ): Promise<void>
    {
        try
        {
            if ( this.state.loading )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            const checkUID = await this._checkService.addCheck({
                name:           values.name.trim(),
                host_uid:       this.props.hostUID,
                checker_uid:    values.checker_uid.trim(),
                schedule:       values.schedule.trim() || null,
                params:         Object.entries(values.params).map( o => ({
                    name:           o[0].trim(),
                    value:          o[1].trim()
                }))
            });

            this.setState({
                loading:    false
            });

            this.props.onFinish(checkUID);
        }
        catch ( err )
        {
            console.error(`Error adding check: ${(err as Error).message}`, err);

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
            <div className='AddCheckStartStep'>
                <h1>Add a check</h1>

                <ErrorBox error={this.state.error} />

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    {({ values }) => (
                        <Form>
                            <div>
                                Host {this.state.host?.name || '-'}
                            </div>

                            <br />

                            <InputText
                                name='name'
                                label='Name'
                                required={true}
                            />

                            <InputSelect
                                name='checker_uid'
                                label='Checker'
                                options={this.state.checkers.map( ( checker ) => ({
                                    label:  checker.type,
                                    value:  checker.uid
                                }))}
                                required={true}
                            />

                            <InputText
                                name='schedule'
                                label='Schedule'
                            />

                            {values.checker_uid ?
                                this.state.checkers.find( o => o.uid === values.checker_uid )?.capabilities.params.map( ( param ) => (
                                    <InputText
                                        name={`params.${param.name}`}
                                        label={param.label}
                                        required={param.required}
                                    />
                                ))
                            : null}

                            <Button
                                type='submit'
                                disabled={this.state.loading}>
                                Create
                            </Button>
                        </Form>
                    )}
                </Formik>

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}
