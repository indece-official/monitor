import { BackendService } from './BackendService';


export interface NotifierV1Filter
{
    tag_uids:       Array<string>;
    critical:       boolean;
    warning:        boolean;
    unknown:        boolean;
    decline:        boolean;
    min_duration:   string;
}


export interface NotifierV1ConfigParamsEmailSmtp
{
    host:       string;
    port:       number;
    user:       string;
    password:   string;
    from:       string;
    to:         Array<string>;
}


export enum NotifierV1ConfigParamsHttpMethod
{
    Get     = 'GET',
    Post    = 'POST',
    Put     = 'PUT'
}


export const NotifierV1ConfigParamsHttpMethods: Array<NotifierV1ConfigParamsHttpMethod> = [
    NotifierV1ConfigParamsHttpMethod.Get,
    NotifierV1ConfigParamsHttpMethod.Post,
    NotifierV1ConfigParamsHttpMethod.Put
];


export interface NotifierV1ConfigParamsHttpHeader
{
    name:   string;
    value:  string;
}


export interface NotifierV1ConfigParamsHttp
{
    url:        string;
    method:     NotifierV1ConfigParamsHttpMethod;
    headers:    Array<NotifierV1ConfigParamsHttpHeader>;
    body:       string | null;
}


export interface NotifierV1ConfigParamsMicrosoftTeams
{
    webhook_url:    string;
}


export interface NotifierV1ConfigParams
{
    email_smtp?:        NotifierV1ConfigParamsEmailSmtp | null;
    http?:              NotifierV1ConfigParamsHttp | null;
    microsoft_teams?:   NotifierV1ConfigParamsMicrosoftTeams | null;
}


export interface NotifierV1Config
{
    filters:    Array<NotifierV1Filter>;
    params:     NotifierV1ConfigParams;
}


export enum NotifierV1Type
{
    EmailSmtp       = 'EMAIL_SMTP',
    Http            = 'HTTP',
    MicrosoftTeams  = 'MICROSOFT_TEAMS'
}


export const NotifierV1Types: Array<NotifierV1Type> = [
    NotifierV1Type.EmailSmtp,
    NotifierV1Type.Http,
    NotifierV1Type.MicrosoftTeams
];


export interface NotifierV1
{
    uid:        string;
    name:       string;
    type:       NotifierV1Type;
    config?:    NotifierV1Config | null;
}


export interface AddNotifierV1Request
{
    name:   string;
    type:   NotifierV1Type;
    config: NotifierV1Config;
}


export interface UpdateNotifierV1Request
{
    name:   string;
    config: NotifierV1Config;
}


export class NotifierService
{
    private static _instance:           NotifierService;
    private readonly _backendService:   BackendService;
    
    
    public static getInstance ( ): NotifierService
    {
        if ( ! this._instance )
        {
            this._instance = new NotifierService();
        }
        
        return this._instance;
    }


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
    }


    public async getNotifiers ( ): Promise<Array<NotifierV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/notifier`,
            {
                method: 'GET',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        return resp.notifiers;
    }
    
    
    public async getNotifier ( notifierUID: string ): Promise<NotifierV1>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/notifier/${encodeURIComponent(notifierUID)}`,
            {
                method: 'GET',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );

        return resp.notifier;
    }

    
    public async deleteNotifier ( notifierUID: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/notifier/${encodeURIComponent(notifierUID)}`,
            {
                method:     'DELETE',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );
    }


    public async addNotifier ( params: AddNotifierV1Request ): Promise<string>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/notifier`,
            {
                method: 'POST',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify(params)
            }
        );

        return resp.notifier_uid;
    }


    public async updateNotifier ( notifierUID: string, params: UpdateNotifierV1Request ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/notifier/${encodeURIComponent(notifierUID)}`,
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
}
