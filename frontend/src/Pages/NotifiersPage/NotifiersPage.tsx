import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { NotifierService, NotifierV1 } from '../../Services/NotifierService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPen, faPlus, faTrash } from '@fortawesome/free-solid-svg-icons';
import { UserService, UserV1, isAdmin } from '../../Services/UserService';
import { List } from '../../Components/List/List';
import { ListEmpty } from '../../Components/List/ListEmpty';
import { ListItem } from '../../Components/List/ListItem';
import { ListItemHeader } from '../../Components/List/ListItemHeader';
import { ListItemHeaderField } from '../../Components/List/ListItemHeaderField';
import { ListItemHeaderAction } from '../../Components/List/ListItemHeaderAction';
import { PageContent } from '../../Components/PageContent/PageContent';


export interface NotifiersPageProps
{
}


interface NotifiersPageState
{
    user:       UserV1 | null;
    notifiers:  Array<NotifierV1>;
    loading:    boolean;
    error:      Error | null;
}


export class NotifiersPage extends React.Component<NotifiersPageProps, NotifiersPageState>
{
    private readonly _userService:      UserService;
    private readonly _notifierService:  NotifierService;


    constructor ( props: NotifiersPageProps )
    {
        super(props);

        this.state = {
            user:       null,
            notifiers:  [],
            loading:    true,
            error:      null
        };

        this._userService = UserService.getInstance();
        this._notifierService = NotifierService.getInstance();
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                error:      null
            });

            const notifiers = await this._notifierService.getNotifiers();

            this.setState({
                loading: false,
                notifiers
            });
        }
        catch ( err )
        {
            console.error(`Error loading notifiers: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        this._userService.isLoggedIn().subscribe(this, ( user ) =>
        {
            this.setState({
                user
            });
        });

        const user = this._userService.isLoggedIn().get();
        this.setState({
            user
        });

        await this._load();
    }


    public componentWillUnmount ( ): void
    {
        this._userService.isLoggedIn().unsubscribe(this);
    }


    public render ( )
    {
        return (
            <PageContent>
                <h1>Notifiers</h1>

                <ErrorBox error={this.state.error} />

                {isAdmin(this.state.user) ?
                    <Button to={LinkUtils.make('notifier', 'add')}>
                        <FontAwesomeIcon icon={faPlus} /> Add a notifier
                    </Button>
                : null}

                <List>
                    {this.state.notifiers.length === 0 && !this.state.loading && !this.state.error ?
                        <ListEmpty>
                            No notifiers found
                        </ListEmpty>
                    : null}

                    {this.state.notifiers.map( ( notifier ) => (
                        <ListItem key={notifier.uid}>
                            <ListItemHeader>
                                <ListItemHeaderField
                                    grow={true}
                                    text={notifier.name}
                                    subtext={notifier.type}
                                />

                                {isAdmin(this.state.user) ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('notifier', notifier.uid, 'edit')}
                                        icon={faPen}
                                    />
                                : null}
                                
                                {isAdmin(this.state.user) ?
                                    <ListItemHeaderAction
                                        to={LinkUtils.make('notifier', notifier.uid, 'delete')}
                                        icon={faTrash}
                                    />
                                : null}
                            </ListItemHeader>
                        </ListItem>
                    ))}
                </List>

                <Spinner active={this.state.loading} />
            </PageContent>
        );
    }
}
