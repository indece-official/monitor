import React from 'react';
import { Spinner } from '../Spinner/Spinner';
import { UserService, UserV1, UserV1Role } from '../../Services/UserService';


class EmptyPage extends React.Component
{
    public render ( )
    {
        return (
            <Spinner />
        );
    }
}


export interface RouteGuardProps
{
    roles?:     Array<UserV1Role>;
    element:    React.ReactNode | null;
}


interface RouteGuardState
{
    user:   UserV1 | null;
}


export class RouteGuard extends React.Component<RouteGuardProps, RouteGuardState>
{
    private readonly _userService:  UserService;


    constructor ( props: RouteGuardProps )
    {
        super(props);

        this.state = {
            user:   null
        };

        this._userService = UserService.getInstance();
    }


    public componentDidMount ( ): void
    {
        this._userService.isLoggedIn().subscribe(this, ( user ) =>
        {
            this.setState({
                user
            });
        });

        const user = this._userService.isLoggedIn().get();
        this.setState({
            user
        });
    }


    public componentWillUnmount ( ): void
    {
        this._userService.isLoggedIn().unsubscribe(this);
    }


    public render ( )
    {
        if ( !this.state.user || (this.props.roles && !this.state.user.roles.find( r => this.props.roles!.includes(r) ) ))
        {
            return (
                <EmptyPage />
            );
        }

        return (
            <>
                {this.props.element}
            </>
        );
    }
}
