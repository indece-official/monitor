import React from 'react';
import { CheckService, CheckV1 } from '../../Services/CheckService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';


export interface DeleteCheckPageRouteParams
{
    hostUID:    string;
    checkUID:   string;
}


export interface DeleteCheckPageProps extends RouteComponentProps<DeleteCheckPageRouteParams>
{
}


interface DeleteCheckPageState
{
    check:      CheckV1 | null;
    loading:    boolean;
    error:      Error | null;
    success:    string | null;
}


class $DeleteCheckPage extends React.Component<DeleteCheckPageProps, DeleteCheckPageState>
{
    private readonly _checkService: CheckService;


    constructor ( props: DeleteCheckPageProps )
    {
        super(props);

        this.state = {
            check:      null,
            loading:    false,
            error:      null,
            success:    null
        };

        this._checkService = CheckService.getInstance();

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

            const check = await this._checkService.getCheck(this.props.router.params.checkUID);

            this.setState({
                loading:    false,
                check
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


    private async _delete ( ): Promise<void>
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

            await this._checkService.deleteCheck(this.state.check.uid);

            this.setState({
                loading:    false,
                success:    'The check was successfully deleted.'
            });
        }
        catch ( err )
        {
            console.error(`Error deleting check: ${(err as Error).message}`, err);

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
                <h1>Delete check</h1>

                <ErrorBox error={this.state.error} />

                <div>Do you really want to delete check {this.state.check ? this.state.check.name : '?'}?</div>

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


export const DeleteCheckPage = withRouter($DeleteCheckPage);
