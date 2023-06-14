import React from 'react';
import { ConnectorV1, ConnectorV1Status, ConnectorService } from '../../../Services/ConnectorService';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { Formatter } from '../../../utils/Formatter';


export interface AddConnectorWaitRegisteredStepProps
{
    connectorUID:   string;
    onFinish:       ( ) => any;
    onFailed:       ( ) => any;
}


interface AddConnectorWaitRegisteredStepState
{
    connector:  ConnectorV1 | null;
    loading:    boolean;
    error:      Error | null;
}


export class AddConnectorWaitRegisteredStep extends React.Component<AddConnectorWaitRegisteredStepProps, AddConnectorWaitRegisteredStepState>
{
    private readonly _connectorService:  ConnectorService;
    private _intervalReload:    any | null;


    constructor ( props: AddConnectorWaitRegisteredStepProps )
    {
        super(props);

        this.state = {
            connector:  null,
            loading:    true,
            error:      null
        };

        this._connectorService = ConnectorService.getInstance();

        this._intervalReload = null;
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

            if ( connector.status === ConnectorV1Status.Unregistered )
            {
                if ( this._intervalReload )
                {
                    clearInterval(this._intervalReload);
                    this._intervalReload = null;
                }

                this.props.onFinish();
            }
            else if ( connector.status === ConnectorV1Status.Error )
            {
                if ( this._intervalReload )
                {
                    clearInterval(this._intervalReload);
                    this._intervalReload = null;
                }

                this.props.onFailed();
            }
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
        this._intervalReload = setInterval( async ( ) => {
            await this._load();
        }, 5000);

        await this._load();
    }


    public componentWillUnmount ( ): void
    {
        if ( this._intervalReload )
        {
            clearInterval(this._intervalReload);
            this._intervalReload = null;
        }
    }

    public render ( )
    {
        return (
            <div className='AddConnectorWaitRegisteredStep'>
                <h1>Add a connector</h1>

                <ErrorBox error={this.state.error} />

                {this.state.connector ?
                    <div>
                        <div>Waiting for connector to be configured ...</div>
                        <div>Status: {Formatter.connectorStatus(this.state.connector)}</div>
                    </div>
                : null}

                <Spinner active={true} />
            </div>
        );
    }
}
