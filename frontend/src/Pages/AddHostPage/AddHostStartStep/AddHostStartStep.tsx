import React from 'react';
import { Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { HostService } from '../../../Services/HostService';
import { InputMultiSelect } from '../../../Components/Input/InputMultiSelect';
import { TagService, TagV1 } from '../../../Services/TagService';


export interface AddHostStartStepProps
{
    onFinish:   ( hostUID: string ) => any;
}


interface AddHostStartStepFormData
{
    name:       string;
    tag_uids:   Array<string>;
}


interface AddHostStartStepState
{
    initialFormData:    AddHostStartStepFormData;
    tags:               Array<TagV1>
    loading:            boolean;
    error:              Error | null;
}


export class AddHostStartStep extends React.Component<AddHostStartStepProps, AddHostStartStepState>
{
    private readonly _hostService: HostService;
    private readonly _tagService: TagService;


    constructor ( props: AddHostStartStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                name:       '',
                tag_uids:   []
            },
            tags:       [],
            loading:    false,
            error:      null
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


    private async _submit ( values: AddHostStartStepFormData ): Promise<void>
    {
        try
        {
            if ( this.state.loading )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            const hostUID = await this._hostService.addHost({
                name:       values.name.trim(),
                tag_uids:   values.tag_uids.map( s => s.trim() )
            });

            this.setState({
                loading:    false
            });

            this.props.onFinish(hostUID);
        }
        catch ( err )
        {
            console.error(`Error adding host: ${(err as Error).message}`, err);

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
                <h1>Add a host</h1>

                <ErrorBox error={this.state.error} />

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    {({ values }) => (
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
                                Create
                            </Button>
                        </Form>
                    )}
                </Formik>

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}
