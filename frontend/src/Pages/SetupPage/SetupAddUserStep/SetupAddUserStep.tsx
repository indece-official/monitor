import React from 'react';
import { Form, Formik } from 'formik';
import { ErrorBox } from '../../../Components/ErrorBox/ErrorBox';
import { InputText } from '../../../Components/Input/InputText';
import { Button } from '../../../Components/Button/Button';
import { Spinner } from '../../../Components/Spinner/Spinner';
import { UserService, UserV1Role } from '../../../Services/UserService';


export interface SetupAddUserStepProps
{
    onFinish:   ( userUID: string ) => any;
}


interface SetupAddUserStepFormData
{
    username:           string;
    name:               string;
    email:              string;
    password:           string;
    password_confirm:   string;
}


interface SetupAddUserStepState
{
    initialFormData:    SetupAddUserStepFormData;
    loading:            boolean;
    error:              Error | null;
}


export class SetupAddUserStep extends React.Component<SetupAddUserStepProps, SetupAddUserStepState>
{
    private readonly _userService: UserService;


    constructor ( props: SetupAddUserStepProps )
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


    private async _submit ( values: SetupAddUserStepFormData ): Promise<void>
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
                    UserV1Role.Admin
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
            <div className='SetupAddUserStep'>
                <h1>Setup - Create admin user</h1>

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
                                name='email'
                                label='Email'
                            />

                            <InputText
                                name='password'
                                label='Password'
                                type='password'
                                required={true}
                            />

                            <InputText
                                name='password_confirm'
                                label='Confirm password'
                                type='password'
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
