import React from 'react';
import { UserService } from '../../Services/UserService';
import { ErrorBox } from '../../Components/ErrorBox/ErrorBox';
import { Form, Formik } from 'formik';
import { InputText } from '../../Components/Input/InputText';
import { Button } from '../../Components/Button/Button';
import { Spinner } from '../../Components/Spinner/Spinner';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';
import { PageContent, PageContentSize } from '../../Components/PageContent/PageContent';
import { Box } from '../../Components/Box/Box';

import './LoginPage.css';
import { LinkUtils } from '../../utils/LinkUtils';


export interface LoginPageProps extends RouteComponentProps
{
}


interface LoginPageFormData
{
    username:   string;
    password:   string;
}


interface LoginPageState
{
    initialFormData:    LoginPageFormData;
    loading:            boolean;
    error:              Error | null;
}


class $LoginPage extends React.Component<LoginPageProps, LoginPageState>
{
    private readonly _userService: UserService;


    constructor ( props: LoginPageProps )
    {
        super(props);

        this.state = {
            initialFormData: {
                username:   '',
                password:   ''
            },
            loading:    false,
            error:      null
        };

        this._userService = UserService.getInstance();

        this._submit = this._submit.bind(this);
    }


    private async _submit ( values: LoginPageFormData ): Promise<void>
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

            await this._userService.login({
                username:       values.username.trim(),
                password:       values.password.trim()
            });

            this.setState({
                loading:    false
            });

            this.props.router.navigate(LinkUtils.make());
        }
        catch ( err )
        {
            console.error(`Error loggin in user: ${(err as Error).message}`, err);

            this.setState({
                loading:    false,
                error:      err as Error
            });
        }
    }


    public render ( )
    {
        return (
            <PageContent
                size={PageContentSize.Small}
                centered={true}
                className='LoginPage'>
                <Box className='LoginPage-login'>
                    <h1>Login</h1>

                    <ErrorBox error={this.state.error} />

                    <Formik
                        initialValues={this.state.initialFormData}
                        onSubmit={this._submit}
                        enableReinitialize={true}>
                        <Form>
                            <InputText
                                name='username'
                                label='Username'
                                required={true}
                            />

                            <InputText
                                name='password'
                                label='Password'
                                type='password'
                                required={true}
                            />

                            <Button
                                type='submit'
                                disabled={this.state.loading}>
                                Login
                            </Button>
                        </Form>
                    </Formik>
                </Box>

                <Spinner active={this.state.loading} />
            </PageContent>
        );
    }
}


export const LoginPage = withRouter($LoginPage);
