import React from 'react';
import { Form, Formik } from 'formik';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { NotifierV1Type, NotifierV1Types } from '../../../Services/NotifierService';
import { InputSelect } from '../../../Components/Input/InputSelect';


export interface AddNotifierNameTypeStepProps
{
    onFinish:   ( name: string, type: NotifierV1Type ) => any;
}


interface AddNotifierNameTypeStepFormData
{
    name:   string;
    type:   string;
}


interface AddNotifierNameTypeStepState
{
    initialFormData:    AddNotifierNameTypeStepFormData;
}


export class AddNotifierNameTypeStep extends React.Component<AddNotifierNameTypeStepProps, AddNotifierNameTypeStepState>
{
    constructor ( props: AddNotifierNameTypeStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                name:       '',
                type:       ''
            }
        };

        this._submit = this._submit.bind(this);
    }


    private async _submit ( values: AddNotifierNameTypeStepFormData ): Promise<void>
    {
        this.props.onFinish(
            values.name.trim(),
            values.type as NotifierV1Type
        );
    }


    public render ( )
    {
        return (
            <div className='AddNotifierNameTypeStep'>
                <h1>Add a notifier</h1>

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

                            <InputSelect
                                name='type'
                                label='Type'
                                options={NotifierV1Types.map( ( s ) => ({
                                    label:  s,
                                    value:  s
                                }))}
                            />

                            <Button type='submit'>
                                Continue
                            </Button>
                        </Form>
                    )}
                </Formik>
            </div>
        );
    }
}
