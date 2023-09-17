import React from 'react';
import { TagService, TagV1 } from '../../Services/TagService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { sleep } from 'ts-delay';
import { LinkUtils } from '../../utils/LinkUtils';
import { PageContent } from '../../Components/PageContent/PageContent';


export interface DeleteTagPageRouteParams
{
    tagUID:    string;
}


export interface DeleteTagPageProps extends RouteComponentProps<DeleteTagPageRouteParams>
{
}


interface DeleteTagPageState
{
    tag:       TagV1 | null;
    loading:    boolean;
    error:      Error | null;
    success:    string | null;
}


class $DeleteTagPage extends React.Component<DeleteTagPageProps, DeleteTagPageState>
{
    private readonly _tagService: TagService;


    constructor ( props: DeleteTagPageProps )
    {
        super(props);

        this.state = {
            tag:   null,
            loading:    false,
            error:      null,
            success:    null
        };

        this._tagService = TagService.getInstance();

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

            const tag = await this._tagService.getTag(this.props.router.params.tagUID);

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


    private async _delete ( ): Promise<void>
    {
        try
        {
            if ( this.state.loading || !this.state.tag )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            await this._tagService.deleteTag(this.state.tag.uid);

            this.setState({
                loading:    false,
                success:    'The tag was successfully deleted.'
            });

            await sleep(1000);

            this.props.router.navigate(LinkUtils.make('tags'));
        }
        catch ( err )
        {
            console.error(`Error deleting tag: ${(err as Error).message}`, err);

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
                <h1>Delete tag</h1>

                <ErrorBox error={this.state.error} />

                <div>Do you really want to delete tag {this.state.tag ? this.state.tag.name : '?'}?</div>

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


export const DeleteTagPage = withRouter($DeleteTagPage);
