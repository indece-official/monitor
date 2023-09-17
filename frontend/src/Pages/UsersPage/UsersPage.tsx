import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { UserService, UserV1 } from '../../Services/UserService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPen, faPlus, faTrash } from '@fortawesome/free-solid-svg-icons';
import { List } from '../../Components/List/List';
import { ListEmpty } from '../../Components/List/ListEmpty';
import { ListItem } from '../../Components/List/ListItem';
import { ListItemHeader } from '../../Components/List/ListItemHeader';
import { ListItemHeaderField } from '../../Components/List/ListItemHeaderField';
import { ListItemHeaderAction } from '../../Components/List/ListItemHeaderAction';
import { Formatter } from '../../utils/Formatter';
import { PageContent } from '../../Components/PageContent/PageContent';


export interface UsersPageProps
{
}


interface UsersPageState
{
    users:      Array<UserV1>;
    loading:    boolean;
    error:      Error | null;
}


export class UsersPage extends React.Component<UsersPageProps, UsersPageState>
{
    private readonly _userService: UserService;


    constructor ( props: UsersPageProps )
    {
        super(props);

        this.state = {
            users:      [],
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

            const users = await this._userService.getUsers();

            this.setState({
                loading:    false,
                users
            });
        }
        catch ( err )
        {
            console.error(`Error loading users: ${(err as Error).message}`, err);

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
                <h1>Users</h1>

                <ErrorBox error={this.state.error} />

                <Button to={LinkUtils.make('user', 'add')}>
                    <FontAwesomeIcon icon={faPlus} /> Add an user
                </Button>

                <List>
                    {this.state.users.length === 0 && !this.state.loading && !this.state.error ?
                        <ListEmpty>
                            No users found.
                        </ListEmpty>
                    : null}

                    {this.state.users.map( ( user ) => (
                        <ListItem key={user.uid}>
                            <ListItemHeader>
                                <ListItemHeaderField
                                    grow={true}
                                    text={user.username}
                                    subtext={`${user.name} | ${user.email}`}
                                    to={LinkUtils.make('user', user.uid)}
                                />

                                <ListItemHeaderField
                                    text={user.roles.map(Formatter.userRole).join(', ')}
                                    subtext={user.source}
                                    to={LinkUtils.make('user', user.uid)}
                                />

                                <ListItemHeaderAction
                                    to={LinkUtils.make('user', user.uid, 'edit')}
                                    icon={faPen}
                                />
                                
                                <ListItemHeaderAction
                                    to={LinkUtils.make('user', user.uid, 'delete')}
                                    icon={faTrash}
                                />
                            </ListItemHeader>
                        </ListItem>
                    ))}
                </List>

                <Spinner active={this.state.loading} />
            </PageContent>
        );
    }
}
