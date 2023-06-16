import React from 'react';
import { UserService } from '../../Services/UserService';
import { RouteComponentProps, withRouter } from '../../utils/withRouter';


export interface AuthControllerProps extends RouteComponentProps
{
}


class $AuthController extends React.Component<AuthControllerProps>
{
    private readonly _userService:  UserService;


    constructor ( props: any )
    {
        super(props);

        this._userService = UserService.getInstance();
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._userService.load();

        if ( !this._userService.isLoggedIn().get() &&
             this.props.router.location.pathname !== '/login' &&
             this.props.router.location.pathname !== '/setup' )
        {
            this.props.router.navigate('/login');
        }
    }


    public render ( )
    {
        return null;
    }
}


export const AuthController = withRouter($AuthController);
