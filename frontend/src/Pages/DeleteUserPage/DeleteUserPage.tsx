import React from 'react';
import { UserService, UserV1 } from '../../Services/UserService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';


export interface DeleteUserPageRouteParams
{
    userUID:    string;
}


export interface DeleteUserPageProps extends RouteComponentProps<DeleteUserPageRouteParams>
{
}


interface DeleteUserPageState
{
    user:       UserV1 | null;
    loading:    boolean;
    error:      Error | null;
    success:    string | null;
}


class $DeleteUserPage extends React.Component<DeleteUserPageProps, DeleteUserPageState>
{
    private readonly _userService: UserService;


    constructor ( props: DeleteUserPageProps )
    {
        super(props);

        this.state = {
            user:   null,
            loading:    false,
            error:      null,
            success:    null
        };

        this._userService = UserService.getInstance();

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

            const user = await this._userService.getUser(this.props.router.params.userUID);

            this.setState({
                loading:    false,
                user
            });
        }
        catch ( err )
        {
            console.error(`Error loading user: ${(err as Error).message}`, err);

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
            if ( this.state.loading || !this.state.user )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            await this._userService.deleteUser(this.state.user.uid);

            this.setState({
                loading:    false,
                success:    'The user was successfully deleted.'
            });
        }
        catch ( err )
        {
            console.error(`Error deleting user: ${(err as Error).message}`, err);

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
            <div className='AddUserStartStep'>
                <h1>Delete user</h1>

                <ErrorBox error={this.state.error} />

                <div>Do you really want to delete user {this.state.user ? this.state.user.username : '?'}?</div>

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


export const DeleteUserPage = withRouter($DeleteUserPage);
