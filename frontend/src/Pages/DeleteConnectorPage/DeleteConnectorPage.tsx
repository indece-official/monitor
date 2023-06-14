import React from 'react';
import { ConnectorService, ConnectorV1 } from '../../Services/ConnectorService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';


export interface DeleteConnectorPageRouteParams
{
    connectorUID:    string;
}


export interface DeleteConnectorPageProps extends RouteComponentProps<DeleteConnectorPageRouteParams>
{
}


interface DeleteConnectorPageState
{
    connector:   ConnectorV1 | null;
    loading:    boolean;
    error:      Error | null;
    success:    string | null;
}


class $DeleteConnectorPage extends React.Component<DeleteConnectorPageProps, DeleteConnectorPageState>
{
    private readonly _connectorService: ConnectorService;


    constructor ( props: DeleteConnectorPageProps )
    {
        super(props);

        this.state = {
            connector:   null,
            loading:    false,
            error:      null,
            success:    null
        };

        this._connectorService = ConnectorService.getInstance();

        this._delete = this._delete.bind(this);
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const connector = await this._connectorService.getConnector(this.props.router.params.connectorUID);

            this.setState({
                loading:    false,
                connector
            });
        }
        catch ( err )
        {
            console.error(`Error loading connector: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    private async _delete ( ): Promise<void>
    {
        try
        {
            if ( this.state.loading || !this.state.connector )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            await this._connectorService.deleteConnector(this.state.connector.uid);

            this.setState({
                loading:    false,
                success:    'The connector was successfully deleted.'
            });
        }
        catch ( err )
        {
            console.error(`Error deleting connector: ${(err as Error).message}`, err);

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
                <h1>Delete connector</h1>

                <ErrorBox error={this.state.error} />

                <div>Do you really want to delete connector {this.state.connector ? this.state.connector.type : '?'}?</div>

                <Button
                    onClick={this._delete}
                    disabled={this.state.loading}>
                    Delete
                </Button>

                <SuccessBox message={this.state.success} />

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}


export const DeleteConnectorPage = withRouter($DeleteConnectorPage);
