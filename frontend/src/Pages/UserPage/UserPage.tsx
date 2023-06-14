import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { UserService, UserV1, UserV1Source } from '../../Services/UserService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPen, faTrash } from '@fortawesome/free-solid-svg-icons';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { LabelValueList } from '../../Components/LabelValueList/LabelValueList';
import { LabelValue } from '../../Components/LabelValueList/LabelValue';
import { Formatter } from '../../utils/Formatter';


export interface UserPageRouteParams
{
    userUID:    string;
}


export interface UserPageProps extends RouteComponentProps<UserPageRouteParams>
{
}


interface UserPageState
{
    user:       UserV1 | null;
    loading:    boolean;
    error:      Error | null;
}


class $UserPage extends React.Component<UserPageProps, UserPageState>
{
    private readonly _userService: UserService;


    constructor ( props: UserPageProps )
    {
        super(props);

        this.state = {
            user:       null,
            loading:    false,
            error:      null
        };

        this._userService  = UserService.getInstance();
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const user = await this._userService.getUser(
                this.props.router.params.userUID
            );

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


    public async componentDidMount ( ): Promise<void>
    {
        await this._load();
    }


    public render ( )
    {
        return (
            <div className='UserPage'>
                <h1>User</h1>

                <ErrorBox error={this.state.error} />

                {this.state.user ?
                    <div className='UserPage-actions'>
                        {this.state.user.source === UserV1Source.Local ?
                            <Button to={LinkUtils.make('user', this.state.user.uid, 'edit')}>
                                <FontAwesomeIcon icon={faPen} /> Edit
                            </Button>
                        : null}
                    
                        <Button to={LinkUtils.make('user', this.state.user.uid, 'delete')}>
                            <FontAwesomeIcon icon={faTrash} /> Delete
                        </Button>
                    </div>
                : null}

                {this.state.user ?
                    <LabelValueList>
                        <LabelValue
                            label='Username'
                            value={this.state.user.username}
                        />
                        
                        <LabelValue
                            label='Name'
                            value={this.state.user.name || '-'}
                        />
                        
                        <LabelValue
                            label='Email'
                            value={this.state.user.email || '-'}
                        />
                       
                        <LabelValue
                            label='Source'
                            value={this.state.user.source}
                        />

                        <LabelValue
                            label='Roles'
                            value={this.state.user.roles.map(Formatter.userRole).join(', ') || '-'}
                        />
                    </LabelValueList>
                : null}

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}


export const UserPage = withRouter($UserPage);
