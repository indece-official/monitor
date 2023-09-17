import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Spinner } from '../../Components/Spinner/Spinner';
import { TagService, TagV1 } from '../../Services/TagService';
import { Button } from '../../Components/Button/Button';
import { LinkUtils } from '../../utils/LinkUtils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPen, faTrash } from '@fortawesome/free-solid-svg-icons';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { LabelValueList } from '../../Components/LabelValueList/LabelValueList';
import { LabelValue } from '../../Components/LabelValueList/LabelValue';
import { PageContent } from '../../Components/PageContent/PageContent';


export interface TagPageRouteParams
{
    tagUID:    string;
}


export interface TagPageProps extends RouteComponentProps<TagPageRouteParams>
{
}


interface TagPageState
{
    tag:        TagV1 | null;
    loading:    boolean;
    error:      Error | null;
}


class $TagPage extends React.Component<TagPageProps, TagPageState>
{
    private readonly _tagService: TagService;


    constructor ( props: TagPageProps )
    {
        super(props);

        this.state = {
            tag:       null,
            loading:    false,
            error:      null
        };

        this._tagService  = TagService.getInstance();
    }


    private async _loadTag ( ): Promise<void>
    {
        try
        {
            this.setState({
                loading:    true,
                error:      null
            });

            const tag = await this._tagService.getTag(
                this.props.router.params.tagUID
            );

            this.setState({
                loading:    false,
                tag
            });
        }
        catch ( err )
        {
            console.error(`Error loading tag: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._loadTag();
    }


    public render ( )
    {
        return (
            <PageContent>
                <h1>Tag</h1>

                <ErrorBox error={this.state.error} />

                {this.state.tag ?
                    <div className='TagPage-actions'>
                        <Button to={LinkUtils.make('tag', this.state.tag.uid, 'edit')}>
                            <FontAwesomeIcon icon={faPen} /> Edit
                        </Button>
                    
                        <Button to={LinkUtils.make('tag', this.state.tag.uid, 'delete')}>
                            <FontAwesomeIcon icon={faTrash} /> Delete
                        </Button>
                    </div>
                : null}

                {this.state.tag ?
                    <LabelValueList>
                        <LabelValue
                            label='Name'
                            value={this.state.tag.name}
                        />

                        <LabelValue
                            label='Color'
                            value={this.state.tag.color}
                        />
                    </LabelValueList>
                : null}

                <Spinner active={this.state.loading} />
            </PageContent>
        );
    }
}


export const TagPage = withRouter($TagPage);
