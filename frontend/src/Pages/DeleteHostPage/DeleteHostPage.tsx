import React from 'react';
import { HostService, HostV1 } from '../../Services/HostService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { sleep } from 'ts-delay';
import { LinkUtils } from '../../utils/LinkUtils';
import { PageContent } from '../../Components/PageContent/PageContent';


export interface DeleteHostPageRouteParams
{
    hostUID:    string;
}


export interface DeleteHostPageProps extends RouteComponentProps<DeleteHostPageRouteParams>
{
}


interface DeleteHostPageState
{
    host:       HostV1 | null;
    loading:    boolean;
    error:      Error | null;
    success:    string | null;
}


class $DeleteHostPage extends React.Component<DeleteHostPageProps, DeleteHostPageState>
{
    private readonly _hostService: HostService;


    constructor ( props: DeleteHostPageProps )
    {
        super(props);

        this.state = {
            host:       null,
            loading:    false,
            error:      null,
            success:    null
        };

        this._hostService = HostService.getInstance();

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

            const host = await this._hostService.getHost(this.props.router.params.hostUID);

            this.setState({
                loading:    false,
                host
            });
        }
        catch ( err )
        {
            console.error(`Error loading host: ${(err as Error).message}`, err);

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
            if ( this.state.loading || !this.state.host )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            await this._hostService.deleteHost(this.state.host.uid);

            this.setState({
                loading:    false,
                success:    'The host was successfully deleted.'
            });

            await sleep(1000);

            this.props.router.navigate(LinkUtils.make('hosts'));
        }
        catch ( err )
        {
            console.error(`Error deleting host: ${(err as Error).message}`, err);

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
                <h1>Delete host</h1>

                <ErrorBox error={this.state.error} />

                <div>Do you really want to delete host {this.state.host ? this.state.host.name : '?'}?</div>

                <Button
                    onClick={this._delete}
                    disabled={this.state.loading}>
                    Delete
                </Button>

                <SuccessBox message={this.state.success} />

                <Spinner active={this.state.loading} />
            </PageContent>
        );
    }
}


export const DeleteHostPage = withRouter($DeleteHostPage);
