import React from 'react';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Form, Formik } from 'formik';
import { InputText } from '../../Components/Input/InputText';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { TagService, TagV1 } from '../../Services/TagService';


export interface EditTagPageRouteParams
{
    tagUID:    string;
}


export interface EditTagPageProps extends RouteComponentProps<EditTagPageRouteParams>
{
}


interface EditTagPageFormData
{
    name:   string;
    color:  string;
}


interface EditTagPageState
{
    initialFormData:    EditTagPageFormData;
    tag:                TagV1 | null;
    loading:            boolean;
    error:              Error | null;
    success:            string | null;
}


class $EditTagPage extends React.Component<EditTagPageProps, EditTagPageState>
{
    private readonly _tagService: TagService;


    constructor ( props: EditTagPageProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                name:       '',
                color:      ''
            },
            tag:        null,
            loading:    false,
            error:      null,
            success:    null
        };

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

            const tag = await this._tagService.getTag(this.props.router.params.tagUID);

            this.setState({
                loading:    false,
                tag,
                initialFormData: {
                    name:       tag.name,
                    color:      tag.color
                }
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


    private async _submit ( values: EditTagPageFormData ): Promise<void>
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

            await this._tagService.updateTag(
                this.state.tag.uid,
                {
                    name:   values.name.trim(),
                    color:  values.color.trim()
                }
            );

            this.setState({
                loading:    false,
                success:    'The tag was successfully updated.'
            });
        }
        catch ( err )
        {
            console.error(`Error updating tag: ${(err as Error).message}`, err);

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
            <div className='AddTagStartStep'>
                <h1>Edit tag</h1>

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
                        
                        <InputText
                            name='color'
                            label='Color'
                            required={true}
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


export const EditTagPage = withRouter($EditTagPage);
