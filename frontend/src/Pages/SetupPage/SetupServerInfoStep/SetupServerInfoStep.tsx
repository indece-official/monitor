import React from 'react';
import { Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { ConfigPropertyV1Key, ConfigService } from '../../../Services/ConfigService';


export interface SetupServerInfoStepProps
{
    onFinish:   ( ) => any;
}


interface SetupServerInfoStepFormData
{
    agent_host: string;
    agent_port: string;
}


interface SetupServerInfoStepState
{
    initialFormData:    SetupServerInfoStepFormData;
    loading:            boolean;
    error:              Error | null;
}


export class SetupServerInfoStep extends React.Component<SetupServerInfoStepProps, SetupServerInfoStepState>
{
    private readonly _configService: ConfigService;


    constructor ( props: SetupServerInfoStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                agent_host: window.location.hostname,
                agent_port: '9440'
            },
            loading:    false,
            error:      null
        };

        this._configService = ConfigService.getInstance();

        this._submit = this._submit.bind(this);
    }


    private async _submit ( values: SetupServerInfoStepFormData ): Promise<void>
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

            await this._configService.setConfigProperty(
                ConfigPropertyV1Key.AgentHost,
                values.agent_host.trim()
            );
            
            await this._configService.setConfigProperty(
                ConfigPropertyV1Key.AgentPort,
                values.agent_port.trim()
            );

            this.setState({
                loading:    false
            });

            this.props.onFinish();
        }
        catch ( err )
        {
            console.error(`Error storing config: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public render ( )
    {
        return (
            <div className='SetupServerInfoStep'>
                <h1>Setup - Server informations</h1>

                <ErrorBox error={this.state.error} />

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    {({ values }) => (
                        <Form>
                            <InputText
                                name='agent_host'
                                label='Hostname for agents to connect to'
                                required={true}
                            />
                           
                            <InputText
                                name='agent_port'
                                label='Port for agents to connect to'
                                required={true}
                            />

                            <Button
                                type='submit'
                                disabled={this.state.loading}>
                                Save
                            </Button>
                        </Form>
                    )}
                </Formik>

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}
