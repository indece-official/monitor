import React from 'react';
import { Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { AgentService } from '../../../Services/AgentService';
import { HostService, HostV1 } from '../../../Services/HostService';
import { InputSelect } from '../../../Components/Input/InputSelect';


export interface AddAgentStartStepProps
{
    hostUID?:   string;
    onFinish:   ( agentUID: string ) => any;
}


interface AddAgentStartStepFormData
{
    host_uid:   string;
}


interface AddAgentStartStepState
{
    initialFormData:    AddAgentStartStepFormData;
    hosts:              Array<HostV1>;
    loading:            boolean;
    error:              Error | null;
}


export class AddAgentStartStep extends React.Component<AddAgentStartStepProps, AddAgentStartStepState>
{
    private readonly _hostService:  HostService;
    private readonly _agentService: AgentService;


    constructor ( props: AddAgentStartStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                host_uid:   ''
            },
            hosts:      [],
            loading:    false,
            error:      null
        };

        this._hostService = HostService.getInstance();
        this._agentService = AgentService.getInstance();

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

            const hosts = await this._hostService.getHosts();

            this.setState({
                loading:    false,
                hosts
            });
        }
        catch ( err )
        {
            console.error(`Error loading hosts: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }

    private _download ( filename: string, content: string ): void
    {
        const element = document.createElement('a');
        element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(content));
        element.setAttribute('download', filename);
      
        element.style.display = 'none';
        document.body.appendChild(element);
      
        element.click();
      
        document.body.removeChild(element);
    }
    
    
    private async _submit ( values: AddAgentStartStepFormData ): Promise<void>
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

            const result = await this._agentService.addAgent({
                host_uid: values.host_uid.trim()
            });

            this.setState({
                loading:    false
            });

            this._download('agent.conf', result.config_file);

            this.props.onFinish(result.agent_uid);
        }
        catch ( err )
        {
            console.error(`Error adding agent: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        this.setState({
            initialFormData: {
                ...this.state.initialFormData,
                host_uid:   this.props.hostUID || ''
            }
        });

        await this._load();
    }


    public componentDidUpdate ( prevProps: AddAgentStartStepProps ): void
    {
        if ( this.props.hostUID !== prevProps.hostUID )
        {
            this.setState({
                initialFormData: {
                    ...this.state.initialFormData,
                    host_uid:   this.props.hostUID || ''
                }
            });
        }
    }


    public render ( )
    {
        return (
            <div className='AddAgentStartStep'>
                <h1>Add a agent</h1>

                <ErrorBox error={this.state.error} />

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    {({ values }) => (
                        <Form>
                            <InputSelect
                                name='host_uid'
                                label='Host'
                                options={this.state.hosts.map( ( host ) => ({
                                    label:  host.name,
                                    value:  host.uid
                                }))}
                                required={true}
                            />

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
