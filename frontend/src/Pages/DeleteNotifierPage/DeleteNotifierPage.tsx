import React from 'react';
import { NotifierService, NotifierV1 } from '../../Services/NotifierService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { sleep } from 'ts-delay';
import { LinkUtils } from '../../utils/LinkUtils';
import { PageContent } from '../../Components/PageContent/PageContent';


export interface DeleteNotifierPageRouteParams
{
    notifierUID:    string;
}


export interface DeleteNotifierPageProps extends RouteComponentProps<DeleteNotifierPageRouteParams>
{
}


interface DeleteNotifierPageState
{
    notifier:   NotifierV1 | null;
    loading:    boolean;
    error:      Error | null;
    success:    string | null;
}


class $DeleteNotifierPage extends React.Component<DeleteNotifierPageProps, DeleteNotifierPageState>
{
    private readonly _notifierService: NotifierService;


    constructor ( props: DeleteNotifierPageProps )
    {
        super(props);

        this.state = {
            notifier:   null,
            loading:    false,
            error:      null,
            success:    null
        };

        this._notifierService = NotifierService.getInstance();

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

            const notifier = await this._notifierService.getNotifier(this.props.router.params.notifierUID);

            this.setState({
                loading:    false,
                notifier
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


    private async _delete ( ): Promise<void>
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

            await this._notifierService.deleteNotifier(this.state.notifier.uid);

            this.setState({
                loading:    false,
                success:    'The notifier was successfully deleted.'
            });

            await sleep(1000);

            this.props.router.navigate(LinkUtils.make('notifiers'));
        }
        catch ( err )
        {
            console.error(`Error deleting notifier: ${(err as Error).message}`, err);

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
                <h1>Delete notifier</h1>

                <ErrorBox error={this.state.error} />

                <div>Do you really want to delete notifier {this.state.notifier ? this.state.notifier.name : '?'}?</div>

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


export const DeleteNotifierPage = withRouter($DeleteNotifierPage);
