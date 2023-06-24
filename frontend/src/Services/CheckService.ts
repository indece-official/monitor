import { BackendService } from './BackendService';


export enum CheckStatusV1Status
{
    Ok          = 'OK',
    Warning     = 'WARNING',
    Critical    = 'CRITICAL',
    Unknown     = 'UNKNOWN'
}


export interface CheckStatusV1Value
{
    name:   string;
    value:  string;
}


export interface CheckStatusV1
{
    uid:                string;
    status:             CheckStatusV1Status;
    message:            string;
    data:               Record<string, any>;
    datetime_created:   string;
}


export interface CheckV1Param
{
    name:   string;
    value:  string;
}


export interface CheckV1
{
    uid:            string;
    name:           string;
    checker_uid:    string;
    custom:         boolean;
    schedule:       string | null;
    disabled:       boolean;
    params:         Array<CheckV1Param>;
    status:         CheckStatusV1 | null;
}


export interface AddCheckV1
{
    name:           string;
    checker_uid:    string;
    schedule:       string | null;
    params:         Array<CheckV1Param>;
}


export interface UpdateCheckV1
{
    name:           string;
    schedule:       string | null;
    params:         Array<CheckV1Param>;
}


export class CheckService
{
    private static _instance:           CheckService;
    private readonly _backendService:   BackendService;
    
    
    public static getInstance ( ): CheckService
    {
        if ( ! this._instance )
        {
            this._instance = new CheckService();
        }
        
        return this._instance;
    }


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
    }


    public async getChecks ( ): Promise<Array<CheckV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/check`,
            {
                method: 'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.checks;
    }
   
   
    public async getHostChecks ( hostUID: string ): Promise<Array<CheckV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/host/${encodeURIComponent(hostUID)}/check`,
            {
                method: 'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.checks;
    }


    public async getCheck ( checkUID: string ): Promise<CheckV1>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/check/${encodeURIComponent(checkUID)}`,
            {
                method: 'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.check;
    }


    public async addCheck ( params: AddCheckV1 ): Promise<string>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/check`,
            {
                method: 'POST',
                headers:    {
                    'Accept':       'application/json',
                    'Content-Type': 'application/json'
                },
                body:   JSON.stringify(params)
            }
        );

        return resp.check_uid;
    }
    
    
    public async updateCheck ( checkUID: string, params: UpdateCheckV1 ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/check/${encodeURIComponent(checkUID)}`,
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
   
   
    public async deleteCheck ( checkUID: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/check/${encodeURIComponent(checkUID)}`,
            {
                method: 'DELETE',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );
    }
   
   
    public async executeCheck ( checkUID: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/check/${encodeURIComponent(checkUID)}/execute`,
            {
                method: 'POST',
                headers:    {
                    'Accept':   'application/json'
                }
            }
        );
    }
}
