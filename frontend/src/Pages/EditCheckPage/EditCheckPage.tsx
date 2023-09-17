import React from 'react';
import { CheckService, CheckV1 } from '../../Services/CheckService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Form, Formik } from 'formik';
import { InputText } from '../../Components/Input/InputText';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { CheckerService, CheckerV1 } from '../../Services/CheckerService';
import { HostService, HostV1 } from '../../Services/HostService';
import { AgentService, AgentV1 } from '../../Services/AgentService';
import { sleep } from 'ts-delay';


export interface EditCheckPageRouteParams
{
    hostUID:    string;
    checkUID:    string;
}


export interface EditCheckPageProps extends RouteComponentProps<EditCheckPageRouteParams>
{
}


interface EditCheckPageFormData
{
    name:           string;
    schedule:       string;
    params:         Record<string, string>;
}


interface EditCheckPageState
{
    initialFormData:    EditCheckPageFormData;
    check:              CheckV1 | null;
    checker:            CheckerV1 | null;
    agent:              AgentV1 | null;
    host:               HostV1 | null;
    loading:            boolean;
    error:              Error | null;
    success:            string | null;
}


class $EditCheckPage extends React.Component<EditCheckPageProps, EditCheckPageState>
{
    private readonly _checkService: CheckService;
    private readonly _checkerService: CheckerService;
    private readonly _agentService: AgentService;
    private readonly _hostService: HostService;


    constructor ( props: EditCheckPageProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                name:           '',
                schedule:       '',
                params:         {}
            },
            check:      null,
            checker:    null,
            agent:      null,
            host:       null,
            loading:    false,
            error:      null,
            success:    null
        };

        this._checkService = CheckService.getInstance();
        this._checkerService = CheckerService.getInstance();
        this._agentService = AgentService.getInstance();
        this._hostService = HostService.getInstance();

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

            const check = await this._checkService.getCheck(this.props.router.params.checkUID);
            const checker = await this._checkerService.getChecker(check.checker_uid);
            const agent = await this._agentService.getAgent(checker.agent_uid);
            const host = await this._hostService.getHost(agent.host_uid);

            const params: Record<string, string> = {};
            for ( const param of check.params )
            {
                params[param.name] = param.value;
            }

            this.setState({
                loading:    false,
                check,
                checker,
                agent,
                host,
                initialFormData: {
                    name:       check.name || '',
                    schedule:   check.schedule || '',
                    params
                }
            });
        }
        catch ( err )
        {
            console.error(`Error loading check: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    private async _submit ( values: EditCheckPageFormData ): Promise<void>
    {
        try
        {
            if ( this.state.loading || !this.state.check )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            await this._checkService.updateCheck(
                this.state.check.uid,
                {
                    name:       values.name.trim(),
                    schedule:   values.schedule.trim() || null,
                    params:     Object.entries(values.params).map( o => ({
                        name:       o[0].trim(),
                        value:      o[1].trim()
                    }))
                }
            );

            this.setState({
                loading:    false,
                success:    'The check was successfully updated.'
            });

            await sleep(1000);

            this.props.router.navigate(-1);
        }
        catch ( err )
        {
            console.error(`Error updating check: ${(err as Error).message}`, err);

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
            <div className='EditCheckPage'>
                <h1>Edit check</h1>

                <ErrorBox error={this.state.error} />

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
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

                        <InputText
                            name='schedule'
                            label='Schedule'
                        />

                        {this.state.checker ? this.state.checker.capabilities.params.map( ( param ) => (
                            <InputText
                                name={`params.${param.name}`}
                                label={param.label}
                                required={param.required}
                            />
                        )) : null}

                        <Button
                            type='submit'
                            disabled={this.state.loading}>
                            Save
                        </Button>
                    </Form>
                </Formik>

                <SuccessBox message={this.state.success} />

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}


export const EditCheckPage = withRouter($EditCheckPage);
