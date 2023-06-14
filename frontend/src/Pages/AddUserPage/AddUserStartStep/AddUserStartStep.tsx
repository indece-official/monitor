import React from 'react';
import { Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { UserService, UserV1Role } from '../../../Services/UserService';


export interface AddUserStartStepProps
{
    onFinish:   ( userUID: string ) => any;
}


interface AddUserStartStepFormData
{
    username:           string;
    name:               string;
    email:              string;
    password:           string;
    password_confirm:   string;
}


interface AddUserStartStepState
{
    initialFormData:    AddUserStartStepFormData;
    loading:            boolean;
    error:              Error | null;
}


export class AddUserStartStep extends React.Component<AddUserStartStepProps, AddUserStartStepState>
{
    private readonly _userService: UserService;


    constructor ( props: AddUserStartStepProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                username:           '',
                name:               '',
                email:              '',
                password:           '',
                password_confirm:   ''
            },
            loading:    false,
            error:      null
        };

        this._userService = UserService.getInstance();

        this._submit = this._submit.bind(this);
    }


    private async _submit ( values: AddUserStartStepFormData ): Promise<void>
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

            const userUID = await this._userService.addUser({
                username:   values.username.trim(),
                name:       values.name.trim() || null,
                email:      values.email.trim() || null,
                password:   values.password.trim(),
                roles:      [
                    UserV1Role.Admin    // TODO
                ]
            });

            this.setState({
                loading:    false
            });

            this.props.onFinish(userUID);
        }
        catch ( err )
        {
            console.error(`Error adding user: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public render ( )
    {
        return (
            <div className='AddUserStartStep'>
                <h1>Add an user</h1>

                <ErrorBox error={this.state.error} />

                <Formik
                    initialValues={this.state.initialFormData}
                    onSubmit={this._submit}
                    enableReinitialize={true}>
                    {({ values }) => (
                        <Form>
                            <InputText
                                name='username'
                                label='Username'
                                required={true}
                            />

                            <InputText
                                name='name'
                                label='Name'
                            />

                            <InputText
                                type='email'
                                name='email'
                                label='Email'
                            />

                            <InputText
                                type='password'
                                name='password'
                                label='Password'
                                required={true}
                            />
                            
                            <InputText
                                type='password'
                                name='password_confirm'
                                label='Confirm password'
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
