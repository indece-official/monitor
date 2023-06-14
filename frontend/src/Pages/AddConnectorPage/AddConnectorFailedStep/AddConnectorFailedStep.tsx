import React from 'react';
import { ConnectorV1, ConnectorService } from '../../../Services/ConnectorService';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';


export interface AddConnectorFailedStepProps
{
    connectorUID:   string;
}


interface AddConnectorFailedStepState
{
    connector:  ConnectorV1 | null;
    loading:    boolean;
    error:      Error | null;
}


export class AddConnectorFailedStep extends React.Component<AddConnectorFailedStepProps, AddConnectorFailedStepState>
{
    private readonly _connectorService:  ConnectorService;


    constructor ( props: AddConnectorFailedStepProps )
    {
        super(props);

        this.state = {
            connector:  null,
            loading:            true,
            error:              null
        };

        this._connectorService = ConnectorService.getInstance();
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const connector = await this._connectorService.getConnector(this.props.connectorUID);

            this.setState({
                loading:    false,
                connector
            });
        }
        catch ( err )
        {
            console.error(`Error loading connector execution ${(err as Error).message}`, err);

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
            <div className='AddConnectorFailedStep'>
                <h1>Add a connector</h1>

                <ErrorBox error={this.state.error} />

                <Spinner active={this.state.loading} />

                {this.state.connector && this.state.connector.error ?
                    <ErrorBox error={new Error(this.state.connector.error)} />
                : null}
            </div>
        );
    }
}
