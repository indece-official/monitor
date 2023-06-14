import React from 'react';
import { Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { ConnectorService } from '../../../Services/ConnectorService';
import { HostService, HostV1 } from '../../../Services/HostService';
import { InputSelect } from '../../../Components/Input/InputSelect';


export interface AddConnectorStartStepProps
{
    onFinish:   ( connectorUID: string ) => any;
}


interface AddConnectorStartStepFormData
{
    host_uid:   string;
}


interface AddConnectorStartStepState
{
    initialFormData:    AddConnectorStartStepFormData;
    hosts:              Array<HostV1>;
    loading:            boolean;
    error:              Error | null;
}


export class AddConnectorStartStep extends React.Component<AddConnectorStartStepProps, AddConnectorStartStepState>
{
    private readonly _hostService:      HostService;
    private readonly _connectorService: ConnectorService;


    constructor ( props: AddConnectorStartStepProps )
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
        this._connectorService = ConnectorService.getInstance();

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
    
    
    private async _submit ( values: AddConnectorStartStepFormData ): Promise<void>
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

            const result = await this._connectorService.addConnector({
                host_uid: values.host_uid.trim()
            });

            this.setState({
                loading:    false
            });

            this._download('connector.conf', result.config_file);

            this.props.onFinish(result.connector_uid);
        }
        catch ( err )
        {
            console.error(`Error adding connector: ${(err as Error).message}`, err);

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
            <div className='AddConnectorStartStep'>
                <h1>Add a connector</h1>

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
