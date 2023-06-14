import React from 'react';
import { Location, NavigateFunction, useLocation, useNavigate, useParams } from 'react-router-dom';

export interface RouteComponentPropsRouter<P extends { [K in keyof P]?: string } = {}> {
    params:   Readonly<P>;
    navigate: NavigateFunction;
    location: Location;
}

export interface RouteComponentProps<P extends { [K in keyof P]?: string } = {}> {
    router:   RouteComponentPropsRouter<P>;
}

export function withRouter <R extends RouteComponentProps<any>, T extends React.ComponentClass<R>> (WrappedComponent: T) {
    return (props: Omit<R, 'router'>) => {
        const params = useParams();
        const navigate = useNavigate();
        const location = useLocation();

        return <WrappedComponent {...props as any} router={{params, navigate, location}} />;
    };
}
