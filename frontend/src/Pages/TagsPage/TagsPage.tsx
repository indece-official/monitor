import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { TagService, TagV1 } from '../../Services/TagService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faTrash } from '@fortawesome/free-solid-svg-icons';
import { List } from '../../Components/List/List';
import { ListEmpty } from '../../Components/List/ListEmpty';
import { ListItem } from '../../Components/List/ListItem';
import { ListItemHeaderField } from '../../Components/List/ListItemHeaderField';
import { ListItemHeader } from '../../Components/List/ListItemHeader';
import { ListItemHeaderAction } from '../../Components/List/ListItemHeaderAction';
import { Tag } from '../../Components/Tag/Tag';


export interface TagsPageProps
{
}


interface TagsPageState
{
    tags:       Array<TagV1>;
    loading:    boolean;
    error:      Error | null;
}


export class TagsPage extends React.Component<TagsPageProps, TagsPageState>
{
    private readonly _tagService: TagService;


    constructor ( props: TagsPageProps )
    {
        super(props);

        this.state = {
            tags:      [],
            loading:    false,
            error:      null
        };

        this._tagService  = TagService.getInstance();
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const tags = await this._tagService.getTags();

            this.setState({
                loading:    false,
                tags
            });
        }
        catch ( err )
        {
            console.error(`Error loading tags: ${(err as Error).message}`, err);

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
            <div className='TagsPage'>
                <h1>Tags</h1>

                <ErrorBox error={this.state.error} />

                <Button to={LinkUtils.make('tag', 'add')}>
                    <FontAwesomeIcon icon={faPlus} /> Add a tag
                </Button>

                <List>
                    {this.state.tags.length === 0 && !this.state.loading && !this.state.error ?
                        <ListEmpty>
                            No tags found
                        </ListEmpty>
                    : null}

                    {this.state.tags.map( ( tag ) => (
                        <ListItem key={tag.uid}>
                            <ListItemHeader>
                                <ListItemHeaderField
                                    to={LinkUtils.make('tag', tag.uid)}
                                    grow={true}>
                                    <Tag color={tag.color}>
                                        {tag.name}
                                    </Tag>
                                </ListItemHeaderField>

                                <ListItemHeaderAction
                                    to={LinkUtils.make('tag', tag.uid, 'delete')}
                                    icon={faTrash}
                                />
                            </ListItemHeader>
                        </ListItem>
                    ))}
                </List>

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}
