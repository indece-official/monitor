import React from 'react';
import { HostService, HostV1 } from '../../Services/HostService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Form, Formik } from 'formik';
import { InputText } from '../../Components/Input/InputText';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { TagService, TagV1 } from '../../Services/TagService';
import { InputMultiSelect } from '../../Components/Input/InputMultiSelect';


export interface EditHostPageRouteParams
{
    hostUID:    string;
}


export interface EditHostPageProps extends RouteComponentProps<EditHostPageRouteParams>
{
}


interface EditHostPageFormData
{
    name:       string;
    tag_uids:   Array<string>;
}


interface EditHostPageState
{
    initialFormData:    EditHostPageFormData;
    host:               HostV1 | null;
    tags:               Array<TagV1>
    loading:            boolean;
    error:              Error | null;
    success:            string | null;
}


class $EditHostPage extends React.Component<EditHostPageProps, EditHostPageState>
{
    private readonly _hostService: HostService;
    private readonly _tagService: TagService;


    constructor ( props: EditHostPageProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                name:       '',
                tag_uids:   []
            },
            host:       null,
            tags:       [],
            loading:    false,
            error:      null,
            success:    null
        };

        this._hostService = HostService.getInstance();
        this._tagService = TagService.getInstance();

        this._submit = this._submit.bind(this);
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
            const host = await this._hostService.getHost(this.props.router.params.hostUID);

            this.setState({
                loading:    false,
                host,
                initialFormData: {
                    name:       host.name || '',
                    tag_uids:   host.tags.map( ( tag ) => tag.uid )
                },
                tags
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


    private async _submit ( values: EditHostPageFormData ): Promise<void>
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

            await this._hostService.updateHost(
                this.state.host.uid,
                {
                    name:       values.name.trim(),
                    tag_uids:   values.tag_uids.map( s => s.trim() )
                }
            );

            this.setState({
                loading:    false,
                success:    'The host was successfully updated.'
            });
        }
        catch ( err )
        {
            console.error(`Error updating host: ${(err as Error).message}`, err);

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
            <div className='AddHostStartStep'>
                <h1>Edit host</h1>

                <ErrorBox error={this.state.error} />

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    <Form>
                        <InputText
                            name='name'
                            label='Name'
                            required={true}
                        />

                        <InputMultiSelect
                            name='tag_uids'
                            label='Tags'
                            options={this.state.tags.map( ( tag ) => ({
                                label:  tag.name,
                                value:  tag.uid
                            }))}
                        />

                        <Button
                            type='submit'
                            disabled={this.state.loading}>
                            Save
                        </Button>
                    </Form>
                </Formik>

                <SuccessBox message={this.state.success} />

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}


export const EditHostPage = withRouter($EditHostPage);
