import React from 'react';
import { FieldArray, Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';
import { TagService } from '../../../Services/TagService';


export interface SetupAddTagsStepProps
{
    onFinish:   ( ) => any;
}


interface SetupAddTagsStepFormDataTag
{
    name:   string;
    color:  string;
}


interface SetupAddTagsStepFormData
{
    tags: Array<SetupAddTagsStepFormDataTag>;
}


interface SetupAddTagsStepState
{
    initialFormData:    SetupAddTagsStepFormData;
    loading:            boolean;
    error:              Error | null;
}


export class SetupAddTagsStep extends React.Component<SetupAddTagsStepProps, SetupAddTagsStepState>
{
    private readonly _tagService: TagService;


    constructor ( props: SetupAddTagsStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                tags: [
                    {
                        name:   'production',
                        color:  '#f00'
                    }
                ]
            },
            loading:    false,
            error:      null
        };

        this._tagService = TagService.getInstance();

        this._submit = this._submit.bind(this);
    }


    private async _submit ( values: SetupAddTagsStepFormData ): Promise<void>
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

            for ( const tag of values.tags )
            {
                await this._tagService.addTag({
                    name:   tag.name.trim(),
                    color:  tag.color.trim()
                });
            }

            this.setState({
                loading:    false
            });

            this.props.onFinish();
        }
        catch ( err )
        {
            console.error(`Error adding tag: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public render ( )
    {
        return (
            <div className='SetupAddTagsStep'>
                <h1>Setup - Add tags</h1>

                <ErrorBox error={this.state.error} />

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    {({ values }) => (
                        <Form>
                            <FieldArray
                                name='tags'>
                                {( arrayHelpers ) => (
                                    <div className='SetupAddTagsStep-tags'>
                                        <Button
                                            type='button'
                                            onClick={() => arrayHelpers.push({
                                                labels: []
                                            })}>
                                            <FontAwesomeIcon icon={faPlus} />
                                            Add a tag
                                        </Button>

                                        <div className='SetupAddTagsStep-tags-list'>
                                            {values.tags.length === 0 ?
                                                <div className='SetupAddTagsStep-tags-empty'>
                                                    Na tags added yet.
                                                </div>
                                            : null}
                                        </div>

                                        {values.tags.map( ( tag, index ) => (
                                            <div key={index} className='AdminAddEventPage-tag'>
                                                <div className='AdminAddEventPage-tag-form'>
                                                    <InputText
                                                        name={`tags.${index}.name`}
                                                        label='Name'
                                                        required={true}
                                                    />
                                                   
                                                    <InputText
                                                        name={`tags.${index}.color`}
                                                        label='Color'
                                                        required={true}
                                                    />

                                                </div>

                                                <Button
                                                    title='Delete tag'
                                                    type='button'
                                                    onClick={() => arrayHelpers.remove(index)}>
                                                    <FontAwesomeIcon icon={faTimes} />
                                                </Button>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </FieldArray>

                            <Button
                                type='submit'
                                disabled={this.state.loading}>
                                Save
                            </Button>
                        </Form>
                    )}
                </Formik>

                <Spinner active={this.state.loading} />
            </div>
        );
    }
}
