import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { HomePage } from './Pages/HomePage/HomePage';
import { SideNav } from './Components/SideNav/SideNav';
import { ConnectorsPage } from './Pages/ConnectorsPage/ConnectorsPage';
import { AddConnectorPage } from './Pages/AddConnectorPage/AddConnectorPage';
import { UsersPage } from './Pages/UsersPage/UsersPage';
import { AddUserPage } from './Pages/AddUserPage/AddUserPage';
import { AuthController } from './Controllers/AuthController/AuthController';
import { LoginPage } from './Pages/LoginPage/LoginPage';
import { SetupPage } from './Pages/SetupPage/SetupPage';
import { RouteGuard } from './Components/RouteGuard/RouteGuard';
import { Environment } from './utils/Environment';
import { Error404Page } from './Pages/Error404Page/Error404Page';
import { HostsPage } from './Pages/HostsPage/HostsPage';
import { UserV1Role } from './Services/UserService';
import { DashboardPage } from './Pages/DashboardPage/DashboardPage';
import { EditUserPage } from './Pages/EditUserPage/EditUserPage';
import { UserPage } from './Pages/UserPage/UserPage';
import { ConnectorPage } from './Pages/ConnectorPage/ConnectorPage';
import { DeleteUserPage } from './Pages/DeleteUserPage/DeleteUserPage';
import { DeleteConnectorPage } from './Pages/DeleteConnectorPage/DeleteConnectorPage';
import { DeleteHostPage } from './Pages/DeleteHostPage/DeleteHostPage';
import { AddHostPage } from './Pages/AddHostPage/AddHostPage';
import { HostPage } from './Pages/HostPage/HostPage';
import { AddCheckPage } from './Pages/AddCheckPage/AddCheckPage';
import { ConfigPage } from './Pages/ConfigPage/ConfigPage';
import { EditHostPage } from './Pages/EditHostPage/EditHostPage';
import { EditTagPage } from './Pages/EditTagPage/EditTagPage';
import { TagPage } from './Pages/TagPage/TagPage';
import { TagsPage } from './Pages/TagsPage/TagsPage';
import { AddTagPage } from './Pages/AddTagPage/AddTagPage';
import { DeleteTagPage } from './Pages/DeleteTagPage/DeleteTagPage';
import { EditCheckPage } from './Pages/EditCheckPage/EditCheckPage';
import { DeleteCheckPage } from './Pages/DeleteCheckPage/DeleteCheckPage';
import { NotifiersPage } from './Pages/NotifiersPage/NotifiersPage';
import { AddNotifierPage } from './Pages/AddNotifierPage/AddNotifierPage';
import { EditNotifierPage } from './Pages/EditNotifierPage/EditNotifierPage';

import './App.css';
import { DeleteNotifierPage } from './Pages/DeleteNotifierPage/DeleteNotifierPage';


export class App extends React.Component
{
    public render ( )
    {
        return (
            <BrowserRouter>
                <div className='App'>
                    <div className='App-content'>
                        <SideNav />

                        <div className='App-content-main'>
                            <Routes>
                                <Route path='/' element={<HomePage />} />

                                <Route path='/login' element={<LoginPage />} />
                                
                                {Environment.setup.enabled ?
                                    <Route path='/setup' element={<SetupPage />} />
                                : null}
                                
                                <Route
                                    path='/dashboard'
                                    element={<RouteGuard element={<DashboardPage />} />}
                                />

                                <Route
                                    path='/connectors'
                                    element={<RouteGuard element={<ConnectorsPage />} />}
                                />

                                <Route
                                    path='/connector/add'
                                    element={<RouteGuard
                                        element={<AddConnectorPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />

                                <Route
                                    path='/connector/:connectorUID'
                                    element={<RouteGuard
                                        element={<ConnectorPage />}
                                    />}
                                />

                                <Route
                                    path='/connector/:connectorUID/delete'
                                    element={<RouteGuard
                                        element={<DeleteConnectorPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />

                                <Route
                                    path='/users'
                                    element={<RouteGuard
                                        element={<UsersPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                                
                                <Route
                                    path='/user/add'
                                    element={<RouteGuard
                                        element={<AddUserPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                                
                                <Route
                                    path='/user/:userUID'
                                    element={<RouteGuard
                                        element={<UserPage />}
                                    />}
                                />
                                
                                <Route
                                    path='/user/:userUID/edit'
                                    element={<RouteGuard
                                        element={<EditUserPage />}
                                    />}
                                />
                                
                                <Route
                                    path='/user/:userUID/delete'
                                    element={<RouteGuard
                                        element={<DeleteUserPage />}
                                    />}
                                />
                                
                                <Route
                                    path='/hosts'
                                    element={<RouteGuard
                                        element={<HostsPage />}
                                    />}
                                />

                                <Route
                                    path='/host/add'
                                    element={<RouteGuard
                                        element={<AddHostPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                               
                                <Route
                                    path='/host/:hostUID'
                                    element={<RouteGuard
                                        element={<HostPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                               
                                <Route
                                    path='/host/:hostUID/edit'
                                    element={<RouteGuard
                                        element={<EditHostPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                                
                                <Route
                                    path='/host/:hostUID/delete'
                                    element={<RouteGuard
                                        element={<DeleteHostPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />

                                <Route
                                    path='/host/:hostUID/check/add'
                                    element={<RouteGuard
                                        element={<AddCheckPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                                
                                <Route
                                    path='/host/:hostUID/check/:checkUID/edit'
                                    element={<RouteGuard
                                        element={<EditCheckPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                                
                                <Route
                                    path='/host/:hostUID/check/:checkUID/delete'
                                    element={<RouteGuard
                                        element={<DeleteCheckPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />

                                <Route
                                    path='/tags'
                                    element={<RouteGuard
                                        element={<TagsPage />}
                                    />}
                                />

                                <Route
                                    path='/tag/add'
                                    element={<RouteGuard
                                        element={<AddTagPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                               
                                <Route
                                    path='/tag/:tagUID'
                                    element={<RouteGuard
                                        element={<TagPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                               
                                <Route
                                    path='/tag/:tagUID/edit'
                                    element={<RouteGuard
                                        element={<EditTagPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                                
                                <Route
                                    path='/tag/:tagUID/delete'
                                    element={<RouteGuard
                                        element={<DeleteTagPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />

                                <Route
                                    path='/notifiers'
                                    element={<RouteGuard
                                        element={<NotifiersPage />}
                                    />}
                                />

                                <Route
                                    path='/notifier/add'
                                    element={<RouteGuard
                                        element={<AddNotifierPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />

                                <Route
                                    path='/notifier/:notifierUID/edit'
                                    element={<RouteGuard
                                        element={<EditNotifierPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                               
                                <Route
                                    path='/notifier/:notifierUID/delete'
                                    element={<RouteGuard
                                        element={<DeleteNotifierPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />
                               
                                <Route
                                    path='/config'
                                    element={<RouteGuard
                                        element={<ConfigPage />}
                                        roles={[UserV1Role.Admin]}
                                    />}
                                />

                                <Route path='*' element={<Error404Page />} />
                            </Routes>
                        </div>
                    </div>
                </div>

                <AuthController />
            </BrowserRouter>
        );
    }
}
