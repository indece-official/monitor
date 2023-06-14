import React from 'react';
import { UserService, UserV1, UserV1Role, isAdmin } from '../../Services/UserService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Form, Formik } from 'formik';
import { InputText } from '../../Components/Input/InputText';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { SuccessBox } from '../../Components/SuccessBox/SuccessBox';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';


export interface EditUserPageRouteParams
{
    userUID:    string;
}


export interface EditUserPageProps extends RouteComponentProps<EditUserPageRouteParams>
{
}


interface EditUserPageFormData
{
    name:       string;
    email:      string;
}


interface EditUserPageState
{
    initialFormData:    EditUserPageFormData;
    user:               UserV1 | null;
    ownUser:            UserV1 | null;
    loading:            boolean;
    error:              Error | null;
    success:            string | null;
}


class $EditUserPage extends React.Component<EditUserPageProps, EditUserPageState>
{
    private readonly _userService: UserService;


    constructor ( props: EditUserPageProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                name:       '',
                email:      ''
            },
            user:           null,
            ownUser:        null,
            loading:        false,
            error:          null,
            success:        null
        };

        this._userService = UserService.getInstance();

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

            const user = await this._userService.getUser(this.props.router.params.userUID);

            this.setState({
                loading:    false,
                user,
                initialFormData: {
                    name: user.name || '',
                    email: user.email || '',
                }
            });
        }
        catch ( err )
        {
            console.error(`Error loading user: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    private async _submit ( values: EditUserPageFormData ): Promise<void>
    {
        try
        {
            if ( this.state.loading || !this.state.user )
            {
                return;
            }

            this.setState({
                loading:    true,
                error:      null
            });

            await this._userService.updateUser(
                this.state.user.uid,
                {
                    name:       values.name.trim() || null,
                    email:      values.email.trim() || null,
                    roles:      [
                        UserV1Role.Admin    // TODO
                    ]
                }
            );

            this.setState({
                loading:    false,
                success:    'The user was successfully updated.'
            });
        }
        catch ( err )
        {
            console.error(`Error updating user: ${(err as Error).message}`, err);

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
            <div className='AddUserStartStep'>
                <h1>Edit user</h1>

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
                            name='email'
                            label='Email'
                            required={true}
                        />

                        {isAdmin(this.state.ownUser) ?
                            /* TODO */ null
                        : null}

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


export const EditUserPage = withRouter($EditUserPage);
