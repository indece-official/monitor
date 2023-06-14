import { StateSubject } from 'ts-subject';
import { BackendService } from './BackendService';


export enum UserV1Source
{
    Local   = 'LOCAL'
}


export enum UserV1Role
{
    Show    = 'SHOW',
    Admin   = 'ADMIN'
}


export const UserRoles: Array<UserV1Role> = [
    UserV1Role.Show,
    UserV1Role.Admin
]


export interface UserV1
{
    uid:        string;
    source:     UserV1Source;
    username:   string;
    name:       string | null;
    email:      string | null;
    roles:      Array<UserV1Role>;
}


export interface AddUserV1Request
{
    username:   string;
    name:       string | null;
    email:      string | null;
    password:   string;
    roles:      Array<string>;
}


export interface LoginV1Request
{
    username:   string;
    password:   string;
}


export interface UpdateUserV1Request
{
    name:       string | null;
    email:      string | null;
    roles:      Array<string>;
}


export interface UpdateUserPasswordV1Request
{
    password:   string;
}


export function isAdmin ( user: UserV1 | null ): boolean
{
    return !!(user && user.roles.includes(UserV1Role.Admin));
}


export class UserService
{
    private static _instance:           UserService;
    private readonly _backendService:   BackendService;
    private readonly _subjectLoggedIn:  StateSubject<UserV1 | null>;
    private readonly _subjectLoaded:    StateSubject<boolean>;
    
    
    public static getInstance ( ): UserService
    {
        if ( ! this._instance )
        {
            this._instance = new UserService();
        }
        
        return this._instance;
    }


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
        this._subjectLoggedIn = new StateSubject<UserV1 | null>(null);
        this._subjectLoaded = new StateSubject<boolean>(false);
    }


    public async getUsers ( ): Promise<Array<UserV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/user`,
            {
                method: 'GET',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        return resp.users;
    }


    public async getOwnUser ( ): Promise<UserV1>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/user/self`,
            {
                method: 'GET',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        return resp.user;
    }


    public async getUser ( userUID: string ): Promise<UserV1>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/user/${encodeURIComponent(userUID)}`,
            {
                method: 'GET',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        return resp.user;
    }


    public async addUser ( params: AddUserV1Request ): Promise<string>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/user`,
            {
                method: 'POST',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify(params)
            }
        );

        return resp.user_uid;
    }
    
    
    public async updateUser ( userUID: string, params: UpdateUserV1Request ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/user/${encodeURIComponent(userUID)}`,
            {
                method: 'PUT',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify(params)
            }
        );
    }
    
    
    public async updateUserPassword ( userUID: string, params: UpdateUserPasswordV1Request ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/user/${encodeURIComponent(userUID)}/password`,
            {
                method: 'PUT',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify(params)
            }
        );
    }
    
    
    public async deleteUser ( userUID: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/user/${encodeURIComponent(userUID)}`,
            {
                method: 'DELETE',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );
    }


    public async load ( ): Promise<void>
    {
        try
        {
            const user = await this.getOwnUser();

            this._subjectLoggedIn.next(user);
            this._subjectLoaded.next(true);
        }
        catch ( err )
        {
            console.info(`Couldn't load own user: ${(err as Error).message}`, err);

            this._subjectLoggedIn.next(null);
            this._subjectLoaded.next(true);
        }
    }


    public async login ( params: LoginV1Request ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/login`,
            {
                method: 'POST',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify(params)
            }
        );

        await this.load();
    }


    public async logout (): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/logout`,
            {
                method: 'POST',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        this._subjectLoggedIn.next(null);
    }


    public isLoggedIn ( ): StateSubject<UserV1 | null>
    {
        return this._subjectLoggedIn;
    }
   
   
    public isLoaded ( ): StateSubject<boolean>
    {
        return this._subjectLoaded;
    }
}
