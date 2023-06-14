import React from 'react';
import { Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { TagService } from '../../../Services/TagService';


export interface AddTagStartStepProps
{
    onFinish:   ( tagUID: string ) => any;
}


interface AddTagStartStepFormData
{
    name:       string;
    color:      string;
}


interface AddTagStartStepState
{
    initialFormData:    AddTagStartStepFormData;
    loading:            boolean;
    error:              Error | null;
}


export class AddTagStartStep extends React.Component<AddTagStartStepProps, AddTagStartStepState>
{
    private readonly _tagService: TagService;


    constructor ( props: AddTagStartStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                name:       '',
                color:      ''
            },
            loading:    false,
            error:      null
        };

        this._tagService = TagService.getInstance();

        this._submit = this._submit.bind(this);
    }


    private async _submit ( values: AddTagStartStepFormData ): Promise<void>
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

            const tagUID = await this._tagService.addTag({
                name:       values.name.trim(),
                color:      values.color.trim()
            });

            this.setState({
                loading:    false
            });

            this.props.onFinish(tagUID);
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
            <div className='AddTagStartStep'>
                <h1>Add a tag</h1>

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
                           
                            <InputText
                                name='color'
                                label='Color'
                                required={true}
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
